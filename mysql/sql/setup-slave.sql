# vim: filetype=mysql

CHANGE MASTER TO
  MASTER_HOST     = '172.17.42.1',
  MASTER_PORT     = 3307,
  MASTER_USER     = 'slave_user',
  MASTER_PASSWORD = 'slavepassword',

  MASTER_CONNECT_RETRY = 1;

START SLAVE;
