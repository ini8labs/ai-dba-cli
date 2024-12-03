# AI-DBA-CLI

A powerful command-line tool to interact with your PostgreSQL databases and record insights. This tool allows you to perform various actions such as logging in and running commands directly from your terminal.

## Features

- **Login with email and password**: Use the `login` command to authenticate with your credentials by providing an email and password using the `-e`/`--email` and `-p`/`--password` flags respectively.
- **Run Postgres Analysis**: `analyse` command enables you to run PostgreSQL commands using the `-c` or `--connection-string` flag with your PostgreSQL connection string.

## Requirements

- **Go Runtime**: If you wish to build the application yourself, you will need Go installed. Download Go from [here](https://golang.org/dl/).
- **PostgreSQL Database**: Ensure you have a PostgreSQL database available for testing and use.

## Installation

### Download Pre-built Binary (Recommended)

1. Go to the [Releases](https://github.com/ini8labs/ai-dba-cli/releases) page.
2. Download the appropriate binary for your system:
   - `dba.exe` for Windows amd64
   - `dba-linux-amd64` for Linux amd64
   - `db-darwin-arm64` for MacOS arm64
3. Clone the repo on your system.
4. Place the downloaded binary in the root folder of the cloned repo.
5. Make sure the `.env` file is in the same directory as the binary.

## Running the Binary (Windows AMD64)

### 1. Login to the Database

Use the `login` command to authenticate to the PostgreSQL database with your email and password.

#### Command:

```bash
./dba.exe login -e <email> -p <password>
```

### 2. Analyze PostgreSQL Database

Use the `analyse` command to connect to a PostgreSQL database using a connection string and run analysis.

```bash
./dba.exe analyse -c <postgres_connection_string>
```

## Running the Binary (MacOS ARM64)

### 1. Login to the Database

Use the `login` command to authenticate to the PostgreSQL database with your email and password.

#### Command:

```bash
./dba-darwin-arm64 login -e <email> -p <password>
```

### 2. Analyze PostgreSQL Database

Use the `analyse` command to connect to a PostgreSQL database using a connection string and run analysis.

```bash
./dba-darwin-arm64 analyse -c <postgres_connection_string>
```

## Running the Binary (Linux AMD64)

### 1. Login to the Database

Use the `login` command to authenticate to the PostgreSQL database with your email and password.

#### Command:

```bash
./dba-linux-amd64 login -e <email> -p <password>
```

### 2. Analyze PostgreSQL Database

Use the `analyse` command to connect to a PostgreSQL database using a connection string and run analysis.

```bash
./dba-linux-amd64 analyse -c <postgres_connection_string>
```