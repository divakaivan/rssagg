#!/bin/bash
set -e

echo "Waiting for the database to start..."
sleep 10

echo "Running database migrations..."
/bin/goose -dir=sqlc/schema postgres "postgres://postgres:postgres@db:5432/rssagg?sslmode=disable" up

echo "Generating SQL code..."
/bin/sqlc generate

echo "Starting the application..."
exec ./main
