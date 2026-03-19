-- mysql-init/init.sql

CREATE USER 'wjhaccount'@'%' IDENTIFIED BY 'Wujiahui789';
GRANT ALL PRIVILEGES ON myvedio_lab4.* TO 'wjhaccount'@'%';
FLUSH PRIVILEGES;
