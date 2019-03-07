

TARGETS = bin/mysql-replication-monitor bin/mxtoolbox-blacklist-monitor

GOFLAGS = --ldflags '-s -w'

all: $(TARGETS)

bin/mysql-replication-monitor:
	go get github.com/go-sql-driver/mysql
	go build ${GOFLAGS} -o $@ github.com/ubogdan/monit-tools/cmd/mysql-replication-monitor

bin/mxtoolbox-blacklist-monitor:
	go build ${GOFLAGS} -o $@ github.com/ubogdan/monit-tools/cmd/mysql-replication-monitor

clean:
	rm bin/*