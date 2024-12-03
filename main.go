package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ini8labs/ai-dba-cli/pkg/handlers"
)

type QueryResult struct {
	Query string                   `json:"query"`
	Data  []map[string]interface{} `json:"data"`
	Error string                   `json:"error"`
}

type Output struct {
	ConnectionString string        `json:"connection_string"`
	Data             []QueryResult `json:"data"`
}

// const WebServerURL = "http://localhost:3000/v1/data"

var (
	OptimizationQuery = `
		SELECT
    c.table_schema AS schema_name,
    c.table_name,
    c.column_name,
    c.data_type,
    c.is_nullable
	FROM
    information_schema.columns c
	WHERE
    c.table_schema = 'public'
    AND c.table_name IN (
        SELECT t.table_name
        FROM information_schema.tables t
        WHERE t.table_schema = 'public'
          AND t.table_type = 'BASE TABLE'
          AND t.table_name NOT LIKE 'pg_%'
          AND t.table_name NOT LIKE 'sql_%'
    )
	ORDER BY
    c.table_name,
    c.ordinal_position;
	`
	PerformanceQuery = `
	SELECT
	ss.queryid,
	ss.query,
	ss.calls,
	ss.total_exec_time AS total_time,  -- Total execution time
	ss.mean_exec_time AS mean_time,    -- Average execution time
	ss.max_exec_time AS max_time,      -- Maximum execution time
	ss.min_exec_time AS min_time,      -- Minimum execution time
	ss.rows,
	ss.shared_blks_hit,
	ss.shared_blks_read,
	ss.shared_blks_written,
	ss.local_blks_hit,
	ss.local_blks_read,
	ss.local_blks_written,
	ss.temp_blks_read,
	ss.temp_blks_written,
	sa.pid AS backend_pid,
	sa.state AS query_state,
	sa.wait_event AS current_wait_event,
	sa.wait_event_type AS current_wait_event_type,
	sa.query_start AS query_start_time,
	sa.state_change AS last_state_change,
	sa.xact_start AS transaction_start_time
	FROM
	pg_stat_statements ss
	LEFT JOIN
	pg_stat_activity sa
	ON ss.query = sa.query
	ORDER BY
	ss.total_exec_time DESC
	LIMIT 10;
	`

	SecurityQuery = `
		WITH role_privileges AS (
		SELECT
			rolname AS role_name,
			rolsuper AS is_superuser,
			rolcreaterole AS can_create_roles,
			rolcreatedb AS can_create_db,
			rolcanlogin AS can_login,
			rolreplication AS can_replicate
		FROM pg_roles
	),
	table_privileges AS (
		SELECT
			grantee,
			table_schema,
			table_name,
			privilege_type
		FROM information_schema.role_table_grants
		WHERE table_schema = 'public' -- Focus on public schema
	),
	connection_settings AS (
		SELECT
			name AS setting_name,
			setting AS value
		FROM pg_settings
		WHERE name IN ('max_connections', 'ssl', 'log_connections', 'log_disconnections')
	),
	active_connections AS (
		SELECT
			datname AS database_name,
			usename AS username,
			client_addr AS client_address,
			backend_start,
			state
		FROM pg_stat_activity
		WHERE state = 'active' -- Focus only on active connections
	)
	SELECT DISTINCT
		rp.role_name,
		rp.is_superuser,
		rp.can_create_roles,
		rp.can_create_db,
		rp.can_login,
		rp.can_replicate,
		tp.table_schema,
		tp.table_name,
		tp.privilege_type,
		cs.setting_name,
		cs.value AS setting_value,
		ac.database_name,
		ac.username,
		ac.client_address,
		ac.backend_start,
		ac.state
	FROM
		role_privileges rp
	LEFT JOIN table_privileges tp ON rp.role_name = tp.grantee
	LEFT JOIN connection_settings cs ON true -- Include all connection settings
	LEFT JOIN active_connections ac ON rp.role_name = ac.username
	WHERE
		rp.is_superuser = true -- Focus on superuser roles
		OR tp.privilege_type IS NOT NULL -- Include roles with privileges
		OR cs.setting_name IS NOT NULL -- Include relevant connection settings
		OR ac.username IS NOT NULL -- Include active connections
	ORDER BY
		rp.role_name, tp.table_name, cs.setting_name;
`
)

func main() {

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		logrus.Fatal("WEBHOOK_URL environment variable is not set")
	}

	// CLI flags
	connStr := flag.String("conn", "", "PostgreSQL connection string")
	host := flag.String("host", "localhost", "PostgreSQL host")
	port := flag.String("port", "5432", "PostgreSQL port")
	user := flag.String("user", "postgres", "PostgreSQL username")
	password := flag.String("password", "", "PostgreSQL password")
	dbname := flag.String("dbname", "", "PostgreSQL database name")

	// Login flags
	login := flag.Bool("login", false, "Login to the CLI")
	loginUsername := flag.String("u", "", "Username for login (use with -login)")
	loginPassword := flag.String("p", "", "Password for login (use with -login)")

	// Query flags
	runQueries := flag.Bool("run", false, "Run queries and send data to webhook")

	flag.Parse()

	// tokenStore := handlers.NewFileTokenStore("./dblyser/v1/internal/tokens.json")
	tokenStore := handlers.NewFileTokenStore("./dblyser/v1/internal/")
	if tokenStore == nil {
		logrus.Errorf("Failed to create token store")
		return
	}

	bseURL := os.Getenv("BASE_URL")
	if bseURL == "" {
		logrus.Fatal("BASE_URL environment variable is not set")
		return
	}

	newCLI := handlers.NewCLI(bseURL, *tokenStore)
	if newCLI == nil {
		logrus.Errorf("Failed to create CLI")
		return
	}

	// Check if login mode is enabled
	if *login {
		if *loginUsername == "" || *loginPassword == "" {
			logrus.Println("Error: -u (username) and -p (password) are required when using -login")
			return
		}

		if err := newCLI.Login(*loginUsername, *loginPassword); err != nil {
			logrus.Infof("Failed to login: %v", err)
			return
		}
	}

	if *runQueries {
		// Build connection string if not provided
		var dsn string
		if *connStr == "" {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				*host, *port, *user, *password, *dbname)
		} else {
			dsn = *connStr
		}

		// Connect to the PostgreSQL database using GORM
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		//close db connection
		sqlDB, err := db.DB()
		if err != nil {
			logrus.WithError(err).Error("failed to get underlying *sql.DB")
			fmt.Println("failed to get underlying *sql.DB")
			return
		}
		defer sqlDB.Close()

		// Check if the database is connected
		if err := sqlDB.Ping(); err != nil {
			logrus.Fatalf("Failed to ping the database: %v", err)
		}

		// Queries to execute
		// queries := []string{OptimizationQuery, PerformanceQuery, SecurityQuery}
		queries := map[string]string{
			"optimization": OptimizationQuery,
			"performance":  PerformanceQuery,
			"security":     SecurityQuery,
		}

		results := []QueryResult{}

		for key, query := range queries {
			var result []map[string]interface{}

			var queryResult QueryResult

			db = db.Session(&gorm.Session{
				Logger: logger.Discard,
			})

			if err := db.Raw(query).Scan(&result).Error; err != nil {
				if strings.Contains(err.Error(), "pg_stat_statements must be loaded via") {
					// Log a concise warning without the full query
					logrus.Warnf("pg_stat_statements is not enabled on the database. Please ensure it is properly configured.")
				} else {
					// Log other errors normally
					logrus.Warnf("Failed to execute query: %v", err)
				}

				continue
			}

			queryResult.Query = key
			queryResult.Data = result

			results = append(results, queryResult)
		}

		var outputData Output

		outputData.ConnectionString = dsn
		outputData.Data = results

		// Serialize results to JSON
		// jsonData, err := json.Marshal(results)
		jsonData, err := json.Marshal(outputData)
		if err != nil {
			logrus.Errorf("Failed to marshal query results to JSON: %v", err)
			return
		}

		// read token from file
		token, err := tokenStore.GetToken("auth_token")
		if err != nil {
			logrus.Errorf("Failed to retrieve token from file: %v", err)
			return
		}

		// Prepare the HTTP request with the Bearer token in the Authorization header
		client := &http.Client{}
		req, err := http.NewRequest("POST", webhookURL, strings.NewReader(string(jsonData)))
		if err != nil {
			logrus.Errorf("Failed to create HTTP request: %v", err)
			return
		}

		// Add Bearer token to the Authorization header
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Failed to send data to the web server: %v", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body: %v", err)
			return
		}

		var jsonRespData map[string]interface{}
		if err := json.Unmarshal(body, &jsonRespData); err != nil {
			logrus.Infof("Response: %s.Failed to unmarshal response body: %v", resp.Status, err)
			return
		}

		msg, ok := jsonRespData["message"]
		if !ok {
			logrus.Infof("Data sent to the web server. Response: %s. %s", resp.Status, jsonRespData)
			return
		}

		logrus.Infof("Data sent to the web server. Response: %s.", handleError(resp.Status, msg.(string)))

		// logrus.Infof("Data sent to the web server. Response: %s. %s", resp.Status, msg)
	}

}

func handleError(status string, errStr string) string {

	if status == "401 Unauthorized" {
		return "Invalid token, try to login in again"
	}

	if errStr == "" {
		return "Data sent to the web server. Response: " + status
	}

	return errStr
}
