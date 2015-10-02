# vim: filetype=mysql

GRANT ALL ON *.*
  TO 'slave_user'@'%' IDENTIFIED BY 'slavepassword';
FLUSH PRIVILEGES;
SHOW MASTER STATUS;

