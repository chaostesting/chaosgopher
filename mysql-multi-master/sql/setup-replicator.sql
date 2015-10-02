# vim: filetype=mysql

CHANGE MASTER TO
  MASTER_HOST     = '172.17.42.1',
  MASTER_PORT     = 3306,
  MASTER_USER     = 'replicator',
  MASTER_PASSWORD = 'replicatorpassword',

  MASTER_CONNECT_RETRY = 1;

START SLAVE;
