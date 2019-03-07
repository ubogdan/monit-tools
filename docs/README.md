
# Monit tools

### Mysql Replication

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

### MxToolbox

##### Monit configuration
Create file /etc/monit.d/blacklist.conf with content
```
check program mysql-replication with path "/usr/local/bin/mxtoolbox-blacklist-monitor -host=mx.domain.com -apikey=123e4567-e89b-12d3-a456-426655440000" every 2 cycles
     if status != 0 then alert
```

