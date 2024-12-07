package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ini8labs/ai-dba-cli/pkg/config"
)

var (
	analyseCmd = &cobra.Command{
		Use:   "analyse",
		Short: "Analyse your database",
		RunE:  analyse,
	}

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

	SecurityQuery = `WITH role_privileges AS (
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
        jsonb_agg(jsonb_build_object(
            'table_schema', table_schema,
            'table_name', table_name,
            'privilege_type', privilege_type
        )) AS privileges
    FROM information_schema.role_table_grants
    WHERE table_schema = 'public' -- Focus on public schema
    GROUP BY grantee
),
connection_settings AS (
    SELECT
        jsonb_agg(jsonb_build_object(
            'setting_name', name,
            'value', setting
        )) AS settings
    FROM pg_settings
    WHERE name IN ('max_connections', 'ssl', 'log_connections', 'log_disconnections')
),
active_connections AS (
    SELECT
        usename AS username,
        jsonb_agg(jsonb_build_object(
            'database_name', datname,
            'client_address', client_addr,
            'backend_start', backend_start,
            'state', state
        )) AS connections
    FROM pg_stat_activity
    WHERE state = 'active' -- Focus only on active connections
    GROUP BY usename
)
SELECT
    rp.role_name,
    rp.is_superuser,
    rp.can_create_roles,
    rp.can_create_db,
    rp.can_login,
    rp.can_replicate,
    COALESCE(tp.privileges, '[]'::jsonb) AS table_privileges,
    COALESCE(cs.settings, '[]'::jsonb) AS connection_settings,
    COALESCE(ac.connections, '[]'::jsonb) AS active_connections
FROM
    role_privileges rp
LEFT JOIN table_privileges tp ON rp.role_name = tp.grantee
LEFT JOIN connection_settings cs ON true -- Include all connection settings
LEFT JOIN active_connections ac ON rp.role_name = ac.username
WHERE
    rp.is_superuser = true -- Focus on superuser roles
    OR tp.privileges IS NOT NULL -- Include roles with privileges
    OR cs.settings IS NOT NULL -- Include relevant connection settings
    OR ac.connections IS NOT NULL -- Include active connections
ORDER BY
    rp.role_name;`
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

func init() {

	analyseCmd.Flags().StringP("connection-string", "c", "", "Database connection string")

	// Set a custom help function
	analyseCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Printf("Example:\n\n\t%s analyse -c postgresql://user:pass@localhost:5432/dbname\n\t%s analyse --connection-string postgresql://user:pass@127.0.0.1:5432/dbname\n\nNote: `localhost` and `120.0.0.1` can be used interchangeably.", Binary, Binary)
	})

	rootCmd.AddCommand(analyseCmd)
}

func analyse(cmd *cobra.Command, args []string) error {
	connectionString, err := cmd.Flags().GetString("connection-string")
	if err != nil {
		return err
	}

	// Load the token
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	if config.Token == "" {
		return fmt.Errorf("Please login using the `login` command.")
	}

	// Build the connection string if not provided
	dsn := connectionString
	if dsn == "" {

		return fmt.Errorf("Connection-string cannot be empty; please provide it using the --connection-string flag")
	}

	dsn, err = validateAndNormalizeConnectionString(dsn)
	if err != nil {
		return err
	}

	// Connect to the PostgreSQL database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}

	db = db.Session(&gorm.Session{
		Logger: logger.Discard,
	})

	// Close the DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}
	defer sqlDB.Close()

	// Check if the database is connected
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("Failed to ping the database: %w. Make sure the database is accessible.", err)
	}

	// Define your queries
	queries := map[string]string{
		"optimization": OptimizationQuery,
		"performance":  PerformanceQuery,
		"security":     SecurityQuery,
	}

	results := []QueryResult{}

	for key, query := range queries {
		var result []map[string]interface{}
		var queryResult QueryResult

		// Execute the query
		if err := db.Raw(query).Scan(&result).Error; err != nil {
			if strings.Contains(err.Error(), "pg_stat_statements must be loaded via") {
				logrus.Warnln("pg_stat_statements is not enabled on the database. Please ensure it is properly configured.")
			}
			// TODO: deprecate
			// }else {
			// 	logrus.Warnf("Failed to execute query (%s): %v\n", key, err)
			// }
			continue
		}

		queryResult.Query = key
		queryResult.Data = result

		results = append(results, queryResult)
	}

	// Prepare the output
	outputData := Output{
		ConnectionString: dsn,
		Data:             results,
	}

	jsonData, err := json.Marshal(outputData)
	if err != nil {
		return fmt.Errorf("Could not process the request. Please try again later.")
	}

	// Send the results to a webhook
	client := &http.Client{}
	req, err := http.NewRequest("POST", WebhookURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("Could not process the request. Please try again later.")
	}

	// Add the token to the Authorization header
	req.Header.Add("Authorization", "Bearer "+config.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could not process the request. Please try again later.")
	}
	defer resp.Body.Close()

	// Process the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Could not process the request. Please try again later.")
	}

	var jsonRespData map[string]interface{}
	if err := json.Unmarshal(body, &jsonRespData); err != nil {
		return fmt.Errorf("Could not process the request. Please try again later.")
	}

	// Handle the response based on the status code
	switch resp.StatusCode {
	case http.StatusOK:
		logrus.Infof("Success: Please log into web app to view the results. %s/dashboard", UIURL)
	case http.StatusUnauthorized:
		return fmt.Errorf("Error: Unauthorized. Please try logging in again.")
	default:
		return fmt.Errorf("Error: Could not process the request. Please try again later.")
	}

	return nil
}

// validateAndNormalizeConnectionString checks if the connection string has all required fields
// and replaces "localhost" with "127.0.0.1".
func validateAndNormalizeConnectionString(dsn string) (string, error) {

	// Ensure the connection string starts with "postgres://" or "postgresql://"
	if !strings.HasPrefix(dsn, "postgresql://") {
		if strings.HasPrefix(dsn, "postgresql:/") {
			return "", errors.New("invalid connection string: missing '://' after the scheme")
		}
		return "", errors.New("invalid connection string: must start with 'postgres://' or 'postgresql://'")
	}

	// Parse the DSN
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("failed to parse connection string: %w", err)
	}

	if parsedDSN.Scheme != "postgresql" {
		return "", errors.New("invalid connection string: must start with 'postgresql://'")
	}

	// Extract query parameters and host information
	host := parsedDSN.Hostname()
	port := parsedDSN.Port()
	user := parsedDSN.User.Username()
	password, _ := parsedDSN.User.Password()
	dbname := strings.Trim(parsedDSN.Path, "/") // Remove leading slash from the path

	// Validate the components
	if host == "" {
		return "", errors.New("connection string must include a host")
	}
	if port == "" {
		return "", errors.New("connection string must include a port")
	}
	if user == "" {
		return "", errors.New("connection string must include a username")
	}
	if password == "" {
		return "", errors.New("connection string must include a password")
	}
	if dbname == "" {
		return "", errors.New("connection string must include a database name")
	}

	// Replace "localhost" with "127.0.0.1"
	if host == "localhost" {
		host = "127.0.0.1"
	}

	// Rebuild the connection string with the normalized host
	parsedDSN.Host = fmt.Sprintf("%s:%s", host, port)

	return parsedDSN.String(), nil
}
