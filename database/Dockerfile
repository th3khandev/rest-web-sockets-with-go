FROM mysql:8.0.29-oracle

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["mysqld", "--user=root", "--skip-grant-tables"]