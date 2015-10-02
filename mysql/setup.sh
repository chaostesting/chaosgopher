#!/bin/bash

. ../shared.sh

MYSQL_MASTER="mysql -h $DOCKER_IP -P 3307 -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE"
MYSQL_SLAVE="mysql -h $DOCKER_IP -P 3308 -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE"

echo $MYSQL_MASTER
echo $MYSQL_SLAVE

println PREP MASTER-SLAVE
cat sql/setup-master.sql | $MYSQL_MASTER
cat sql/setup-slave.sql | $MYSQL_SLAVE

println SEED MASTER
cat sql/seed.sql | $MYSQL_MASTER

println QUERY MASTER
sleep 1
echo "SELECT * FROM todo_items;" | $MYSQL_MASTER

println CHECK SLAVE
sleep 1
echo "SELECT * FROM todo_items;" | $MYSQL_SLAVE

