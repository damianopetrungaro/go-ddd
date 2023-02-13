FROM docker.io/migrate/migrate:v4.15.2 AS migrate

FROM postgres:14.6

WORKDIR /home/postgres/migrations

COPY --from=migrate /migrate /migrate
COPY config/database/migrations/ /migrations/
COPY config/database/init-db.sh /docker-entrypoint-initdb.d/
