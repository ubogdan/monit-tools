### Mysql replication watcher

#### Install intructions

##### build application
```
make 
make install
```

##### Grant mysql access to app
```
CREATE USER 'monit'@'localhost' IDENTIFIED BY '';
GRANT SUPER,REPLICATION CLIENT ON *.* TO 'monit'@'localhost';
```

##### Monit configuration
Create file /etc/monit.d/replication.conf with content
```
check program mysql-replication with path "/usr/local/bin/mysql-replication-monitor -u=monit"
     if status != 0 then alert
```