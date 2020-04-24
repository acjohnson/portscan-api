package database

import (
	"database/sql"
	"github.com/acjohnson/portscan-api/logger"
	"log"
)

func Tables(db *sql.DB) error {
	var err error

	logger, err := logger.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create hosts table
	const sql_hosts = `CREATE TABLE IF NOT EXISTS hosts (` +
		`id int NOT NULL PRIMARY KEY AUTO_INCREMENT, ` +
		`ipv4 int unsigned UNIQUE)`
	_, err = db.Exec(sql_hosts)
	if err != nil {
		logger.Println(err)
		panic(err.Error())
	}

	// Create scans table
	const sql_scans = `CREATE TABLE IF NOT EXISTS scans (` +
		`id int NOT NULL PRIMARY KEY AUTO_INCREMENT, ` +
		`host_id int NOT NULL, ` +
		`CONSTRAINT host_id FOREIGN KEY (host_id) REFERENCES hosts(id), ` +
		`last_scanned timestamp DEFAULT CURRENT_TIMESTAMP, ` +
		`port_number int NOT NULL, ` +
		`port_status varchar(32))`
	_, err = db.Exec(sql_scans)
	if err != nil {
		logger.Println(err)
		panic(err.Error())
	}
	return err
}
