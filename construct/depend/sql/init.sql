-- 创建一个单独的用户
CREATE USER IF NOT EXISTS 'jjzzchtt'@'%' IDENTIFIED BY 'jjzzchtt';
GRANT ALL PRIVILEGES ON *.* TO 'jjzzchtt'@'%';

create database if not exists mall character set utf8mb4;
