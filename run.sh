#!/bin/bash

# Local DB connection
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=payroll
export DB_USER=payroll_db_user
export DB_PASSWORD=jw8s0F4
export DB_PASSWORD=jw8s0F4
export DB_SCHEMA=payroll
export SSL_MODE=false


go build -o payroll-service ./cmd/api

if [ $? -ne 0 ]; then
        echo "Build failed"
        exit $?
fi

./payroll-service