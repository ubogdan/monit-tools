package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var socket, credentials, slave string
var delay int64

func init() {
	var username, password string
	flag.StringVar(&socket, "socket", "", "Mysql unix socket (if empty test will run over tcp 127.0.0.1:3306)")
	flag.StringVar(&username, "u", "root", "Mysql username")
	flag.StringVar(&password, "p", "", "Mysql database")
	flag.StringVar(&slave, "slave", "", "Mariadb Slave name)")
	flag.Int64Var(&delay, "delay", 60, "Max seconds behind master that will trigger alert (default 60s)")
	flag.Parse()
	credentials = username
	if password != "" {
		credentials += ":" + password
	}
}

func main() {

	address := "unix(" + socket + ")"
	if socket == "" {
		address = "tcp(127.0.0.1:3306)"
	}

	db, err := sql.Open("mysql", credentials+"@"+address+"/")
	if err != nil {
		ExitState(1000, "Invalid dsn %s", err)
	}
	//db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(1)

	if err = db.Ping(); err != nil {
		ExitState(1016, "Connect: %s", err)
	}

	query_slave := "show slave status"
	if slave != "" {
		query_slave = fmt.Sprintf("show slave %q status", slave)
	}

	rows, err := db.Query(query_slave)
	if err != nil {
		ExitState(1065, "Query: %s", err)
	}

	if !rows.Next() {
		ExitState(0, "Server is not running as slave", err)
	}

	cols, _ := rows.ColumnTypes()

	status := &struct {
		SlaveIORunning      string
		SlaveSQLRunning     string
		LastIOError         string
		LastError           string
		LastErrorNo         int
		SecondsBehindMaster sql.NullInt64
	}{}

	result := make([]interface{}, 0)
	for _, col := range cols {

		switch col.Name() {
		case "Slave_IO_Running":
			result = append(result, &status.SlaveIORunning)

		case "Slave_SQL_Running":
			result = append(result, &status.SlaveSQLRunning)

		case "Last_Error":
			result = append(result, &status.LastError)

		case "Last_Errno":
			result = append(result, &status.LastErrorNo)

		case "Last_IO_Error":
			result = append(result, &status.LastIOError)

		case "Seconds_Behind_Master":
			result = append(result, &status.SecondsBehindMaster)

		default:
			receiver := sql.NullString{}
			result = append(result, &receiver)
		}
	}

	err = rows.Scan(result...)
	if err != nil {
		ExitState(1066, "Query error %s", err)
	}
	defer rows.Close()
	db.Close()

	if status.SlaveIORunning == "Yes" && status.SlaveSQLRunning == "Yes" {
		if status.SecondsBehindMaster.Valid && status.SecondsBehindMaster.Int64 > delay {
			ExitState(5, "Replication is lagging (%d seconds behind master)", status.SecondsBehindMaster.Int64)
		}
		ExitState(0, "Resuming normal operation")
	}

	if status.LastErrorNo > 0 {
		ExitState(10, "%s", status.LastError)
	}

	if status.SlaveIORunning == "No" {
		ExitState(20, "Lost remote: %s", status.LastIOError)
	}

	ExitState(100, "Unhandled software exception . Please investigate !")

}

func ExitState(code int, format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", params...)
	os.Exit(code)
}
