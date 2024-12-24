# Overview

Welcome to Dblyser, your go-to solution for PostgreSQL database analysis. Our platform empowers you to analyze both local databases on your machine and hosted databases on servers, offering tailored insights and actionable recommendations. With web application and command-line interface (CLI) tools, you can optimize, secure, and understand your databases like never before.

## Key Features  

- **Database Analysis**:
  Effortlessly Connect and Analyze PostgreSQL Databases
  - **Local Database Analysis**  
  Perfect for developers working in local environments or testing setups via Dblyser CLI.  

  - **Hosted Database Analysis**  
  Analyze databases hosted on remote servers, including private data centers and cloud platforms, with secure connections and precise insights.  

- **Performance Bottleneck Detection**  
  Identify and resolve slow queries, indexing issues, and other performance hindrances.  

- **Security Vulnerability Assessment**  
  Pinpoint potential vulnerabilities in your database structure, roles, and permissions.  

- **Optimization Opportunities**  
  Receive actionable recommendations to improve database efficiency and security.  

- **Saved Connections**  
  Easily reconnect to previously analyzed databases for streamlined workflows.  

- **Seamless Usability**  
  Use the **Web Application** for a fully interactive experience, or the **Dblyser CLI** specifically for analyzing PostgreSQL databases running on your local machine.

## Prerequisites

To ensure Dblyser can analyze your PostgreSQL database effectively, you need to have the `pg_stat_statements` extension installed. This extension is used to track and report statistics on SQL queries executed by PostgreSQL, which is crucial for performance analysis.

## [Web Application: Interactive Database Analysis](../webapp/README.md)

The **Dblyser Web Application** offers an intuitive, user-friendly interface to manage and analyze your databases visually.  

#### Key Highlights  

- **Visual Dashboards**:  
  - Gain actionable insights into performance, security, and optimization.  
- **Hosted Database Support**:  
  - Connect to PostgreSQL databases on hosted on remote servers.  
- **Automated Recommendations**:  
  - Get tailored suggestions for improving performance and tightening security.  

#### Who Is It For?  

- Developers and admins managing cloud-hosted databases.  
- Teams needing a visual, collaborative analysis platform.  

## [Dblyser CLI: Local Database Analysis via Command Line](../../README.md)

The **Dblyser CLI** enables users to access and analyze their local PostgreSQL databases directly from the terminal. The tool allows you to collect detailed information from your local database and securely send it to our server for analysis.


### Key Highlights  

- **Local and Hosted Database Support**:  
  - Connect to and send analysis data from both local and remote databases.   
- **Seamless Integration**:  
  - Once the analysis data is sent, the results are recorded and can be viewed on our website at [https://dblyser.com](https://dblyser.com).


### Who Is It For?

- Developers working on local databases.
- Engineers managing remote or cloud-hosted databases.


## Data Privacy and Security

At Dblyser, we prioritize your data privacy and security. We ensure that our platform operates in compliance with the highest standards to maintain your trust. Here’s what you need to know about the data we access and store:

### What We Access

- **Connection Details**:
  - To analyze your database, we securely use the connection information you provide (host, port, username, etc.).

- **Database Metadata**:
  - During analysis, we only run SQL queries to gather:
    - **Table Information**: Structure and schema details.
    - **Performance Statistics**: Data from pg_stat_statements and related extensions.
    - **Permissions**: User roles and access levels.

### What We Don’t Access

- **Your Data**: We never look at or store the actual data contained within your tables. Your sensitive information remains untouched.

### What We Store

- **Connection Details**:
  - We store only the connection details you provide for analysis, and this is encrypted and used solely to facilitate your requested analyses.

### Security Assurances

- **Compliant Practices**: We adhere to industry best practices for data privacy and security.
- **Encrypted Connections**: All communication with your database is encrypted, ensuring data remains safe in transit.
- **Transparency**: Our analysis methods involve only necessary queries to provide actionable insights without compromising your data integrity.