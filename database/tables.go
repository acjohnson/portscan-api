package database

import (
	"database/sql"
	"log"
)

func Tables(db *sql.DB) error {
	var err error
	// Create hosts table
	const hosts_sql = `CREATE TABLE IF NOT EXISTS hosts (` +
		`id int NOT NULL PRIMARY KEY AUTO_INCREMENT, ` +
		`ipv4 int unsigned UNIQUE)`
	_, err = db.Exec(hosts_sql)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	// Create scans table
	const scans_sql = `CREATE TABLE IF NOT EXISTS scans (` +
		`id int NOT NULL PRIMARY KEY AUTO_INCREMENT, ` +
		`host_id int NOT NULL, ` +
		`CONSTRAINT host_id FOREIGN KEY (host_id) REFERENCES hosts(id), ` +
		`last_scanned timestamp DEFAULT CURRENT_TIMESTAMP, ` +
		`port_number int NOT NULL, ` +
		`port_status varchar(32))`
	_, err = db.Exec(scans_sql)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return err
}
