#!/bin/bash

. ../shared.sh

MYSQL_ALPHA="mysql -h $DOCKER_IP -P 3306 -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE"
MYSQL_BETA="mysql -h $DOCKER_IP -P 3307 -u root -p$MYSQL_ROOT_PASSWORD -D $MYSQL_DATABASE"

echo $MYSQL_ALPHA
echo $MYSQL_BETA

println SETUP REPLICATION CYCLE
cat sql/setup-master.sql | $MYSQL_ALPHA
sleep 1;
cat sql/setup-replicator.sql | $MYSQL_BETA
sleep 1;
cat sql/setup-reverse-replicator.sql | $MYSQL_ALPHA

println SEED ALPHA
cat sql/seed.sql | $MYSQL_ALPHA

println QUERY ALPHA
echo "SELECT * FROM todo_items;" | $MYSQL_ALPHA

println CHECK BETA
sleep 1
echo "SELECT * FROM todo_items;" | $MYSQL_BETA

