
GOPATH := ${PWD}
export GOPATH

TARGETS = bin/mysql-replication-monitor

GOFLAGS = --ldflags '-s -w'

all: $(TARGETS)

bin/mysql-replication-monitor:
	go get github.com/go-sql-driver/mysql
	go build ${GOFLAGS} -o $@ cmd/mysql-replication-monitor/main.go

clean:
	rm bin/*