FROM ubuntu:20.04

RUN apt-get update && apt-get install -y mysql-client netcat

COPY scripts/proxysql_init.sh /init.sh

CMD cat init.sh && /init.sh
