#!/bin/sh

while ! nc -z $PROXY_HOST $PROXY_PORT; do
  sleep 0.1
done

echo "INSERT INTO mysql_servers(hostgroup_id,hostname,port) VALUES (1, '${PROXY_MYSQL_HOST}', ${PROXY_MYSQL_PORT});" > init.sql

echo "INSERT INTO mysql_users(username,password,default_hostgroup) VALUES ('${MYSQL_USER}','${MYSQL_PASSWORD}',1);" >> init.sql

echo "UPDATE global_variables SET variable_value='${MYSQL_USER}' WHERE variable_name='mysql-monitor_username';" >> init.sql
echo "UPDATE global_variables SET variable_value='${MYSQL_PASSWORD}' WHERE variable_name='mysql-monitor_password';" >> init.sql

echo "LOAD MYSQL USERS TO RUNTIME;" >> init.sql
echo "SAVE MYSQL USERS TO DISK;" >> init.sql

echo "LOAD MYSQL SERVERS TO RUNTIME;" >> init.sql
echo "SAVE MYSQL SERVERS TO DISK;" >> init.sql

echo "LOAD MYSQL VARIABLES TO RUNTIME;" >> init.sql
echo "SAVE MYSQL VARIABLES TO DISK;" >> init.sql

cat init.sql | mysql -u $PROXY_ADMIN_USER -p$PROXY_ADMIN_PASS -h $PROXY_HOST -P$PROXY_PORT