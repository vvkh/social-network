FROM migrate/migrate

COPY migrations /migrations
ENTRYPOINT /bin/sh -c "migrate -path /migrations -database  'mysql://${DB_URL}' up"
