# AI-DBA-CLI

**AI-DBA-CLI** is a powerful command-line tool designed to interact with PostgreSQL databases and record insights. It allows users to log in and analyze database performance directly from the terminal.

---

## Features

- **Secure Login**: Authenticate using your email and password with the `login` command.  
  Example:
  ```bash
  ./dba.exe login -e <email> -p <password>
  ```
- **Database Analysis**: Analyze your PostgreSQL database using the `analyse` command with a connection string.  
  Example:
  ```bash
  ./dba.exe analyse -c <postgres_connection_string>
  ```

---

## Requirements

- **Go Runtime**: Required to build the application from source. [Download Go](https://golang.org/dl/).
- **PostgreSQL Database**: Ensure you have access to a PostgreSQL database for testing and analysis.

---

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Visit the [Releases](https://github.com/ini8labs/ai-dba-cli/releases) page.
2. Download the appropriate binary for your operating system:
   - **Windows**: `dba.exe`
   - **Linux (AMD64)**: `dba-linux-amd64`
   - **MacOS (ARM64)**: `dba-darwin-arm64`
3. Clone the repository to your local system:
   ```bash
   git clone https://github.com/ini8labs/ai-dba-cli.git
   ```
4. Place the downloaded binary in the root folder of the cloned repository.

### Option 2: Build from Source

1. Ensure **Go** is installed on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/ini8labs/ai-dba-cli.git
   cd ai-dba-cli
   ```
3. Build the binary:
   ```bash
   go build -o dba
   ```

---

## Usage

### Running the Binary on **Windows (AMD64)**

1. **Login to the Database**  
   Authenticate using your email and password:
   ```bash
   ./dba.exe login -e your-email@example.com -p your-password
   ```

2. **Analyze PostgreSQL Database**  
   Connect to your database using a connection string and perform analysis:
   ```bash
   ./dba.exe analyse -c "postgres://username:password@localhost:5432/database_name"
   ```

### Running the Binary on **MacOS (ARM64)**

1. **Login to the Database**  
   Authenticate using your email and password:
   ```bash
   ./dba-darwin-arm64 login -e your-email@example.com -p your-password
   ```

2. **Analyze PostgreSQL Database**  
   Connect to your database using a connection string and perform analysis:
   ```bash
   ./dba-darwin-arm64 analyse -c "postgres://username:password@localhost:5432/database_name"
   ```

### Running the Binary on **Linux (AMD64)**

1. **Login to the Database**  
   Authenticate using your email and password:
   ```bash
   ./dba-linux-amd64 login -e your-email@example.com -p your-password
   ```

2. **Analyze PostgreSQL Database**  
   Connect to your database using a connection string and perform analysis:
   ```bash
   ./dba-linux-amd64 analyse -c "postgres://username:password@localhost:5432/database_name"
   ```

---

## Contribution and Feedback

Feel free to contribute to the development of AI-DBA-CLI by creating issues or submitting pull requests on the [GitHub repository](https://github.com/ini8labs/ai-dba-cli). Feedback and suggestions are always welcome!

---

**AI-DBA-CLI** — Simplify your PostgreSQL management and analysis.
