FROM migrate/migrate

COPY migrations /migrations
COPY ./migrate.sh ./migrate.sh
ENTRYPOINT ./migrate.sh
