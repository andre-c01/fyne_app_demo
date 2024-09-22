# fyne_app_demo
A first hands demo of a fyny (golang) app &amp; mysql

# Enviremont Setup

## 2 containers: dev , mariadb server

### mariadb server (debian bookworm)

```shell
apt update && apt install mariadb-server netplan.io -y

printf "
network:
  ethernets:
    enp0s3:
      dhcp4: true
    enp0s8:
      dhcp4: false
      addresses: [192.168.56.101/24]
      nameservers:
        addresses: [127.0.0.1,8.8.8.8]
      dhcp6: false
  version: 2
" > /etc/netplan/51-user-set.yaml

netplan apply

sed -i '/bind-address/c\bind-address\t\t= 192.168.56.101' /etc/mysql/mariadb.conf.d/50-server.cnf

systemctl enable mariadb --now

systemctl restart mariadb

mariadb -u root -p

mariadb  -s -u root -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('System32');"
mariadb  -s -u root -e "DELETE FROM mysql.user WHERE User='';"
mariadb  -s -u root -e "DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');"
mariadb  -s -u root -e "DROP DATABASE IF EXISTS test;"
mariadb  -s -u root -e "DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';"
mariadb  -s -u root -e "FLUSH PRIVILEGES;"

mariadb  -s -u root -e "CREATE USER 'teacher'@'192.168.56.1' IDENTIFIED BY 'System32';"

mariadb  -s -u root -e "CREATE DATABASE school;"

mariadb  -s -u root -D "school" -e "CREATE TABLE students(
	id INT auto_increment PRIMARY KEY,
	name VARCHAR(255) not null UNIQUE,
	passwd VARCHAR(255) not null);"

mariadb  -s -u root -e "GRANT CREATE, ALTER, DROP, INSERT, UPDATE, DELETE, SELECT, REFERENCES on school.* TO 'teacher'@'192.168.56.1' WITH GRANT OPTION;"

mariadb  -u root -D "school" -e "desc students;"

```
