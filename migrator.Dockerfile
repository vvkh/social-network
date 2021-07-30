FROM migrate/migrate

COPY migrations /migrations
COPY scripts/migrate.sh /migrate.sh
ENTRYPOINT /migrate.sh
