#!/bin/bash
set -e

pushd /home/postgres/migrations/
/migrate -database "postgres://${POSTGRES_USER}@/${POSTGRES_DB}?host=/var/run/postgresql" -path "/migrations" up
pg_isready -d order-service


