
# Introduction
The payroll The service provides the ability to download a report on working hours and generate payment reports for employees based on the downloaded data.

# Getting Started
TODO: Guide users through getting your code up and running on their own system. In this section, you can talk about:
1.	Installation process:
   - For build app needs GO version 1.21.1 and higher.
   - For download all dependencies need for build project run commands:
    `go mod tidy`
    `go mod download`
    `go mod vendor`
   - To run the application, you need to connect to the Postgres Database.
      Below is the Docker command to download and run the Postgres database image.
     - `docker pull Postgres` - download docker image of Postgres DB
     - `docker run –name pgsql-dev -e POSTGRES_PASSWORD=Welcome4$ -p 5432:5432 Postgres` - run Docker with Postgres DB on default postgres port 5432
   - Set variables for DB connection in file run.sh
   - Run file from console: `sh run.sh` - it build project and started

# Build and Test
1. For build project run command: `go build -o payroll-service ./cmd/api`
1. For run unit tests run command: go test -v -cover `go list ./...`

# Documentation
Swagger available in http://localhost:8080/doc/api/

