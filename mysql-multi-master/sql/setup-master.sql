# vim: filetype=mysql

CREATE USER 'replicator'@'%' IDENTIFIED BY 'replicatorpassword';
GRANT ALL ON *.* TO 'replicator'@'%';
FLUSH PRIVILEGES;
SHOW MASTER STATUS;

