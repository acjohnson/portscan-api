package database

import (
	"database/sql"
	"fmt"
	"github.com/acjohnson/portscan-api/logger"
	"log"
	"net/url"
	"strconv"
)

func QueryScans(db *sql.DB, values url.Values) (map[string]map[string]string, error) {
	var err error
	var rows *sql.Rows
	var sql_str string
	var query_parm string

	logger, err := logger.Load("DEBUG")
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]map[string]string)

	var (
		id           int
		host_id      int
		last_scanned string
		port_number  int
		port_status  string
	)

	sql_str = "select id, " +
		"host_id, " +
		"last_scanned, " +
		"port_number, " +
		"port_status from scans"

	if values.Get("id") != "" {
		query_parm = values.Get("id")
		sql_str = "select id, " +
			"host_id, " +
			"last_scanned, " +
			"port_number, " +
			"port_status from scans where id = ?"
	}
	if values.Get("ipv4") != "" {
		query_parm = values.Get("ipv4")
		sql_str = "select s.id, " +
			"s.host_id, " +
			"s.last_scanned, " +
			"s.port_number, " +
			"s.port_status from hosts h " +
			"inner join scans s on h.id = s.host_id " +
			"where h.id = (select id from hosts where ipv4 = INET_ATON(?))"

	}

	stmt, err := db.Prepare(sql_str)
	defer stmt.Close()

	if query_parm != "" {
		rows, err = stmt.Query(query_parm)
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		logger.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &host_id, &last_scanned, &port_number, &port_status); err != nil {
			return m, err
		}

		m[fmt.Sprintf("port_number_%s_host_id_%d", strconv.Itoa(port_number), host_id)] = map[string]string{
			"host_id":      strconv.Itoa(host_id),
			"id":           strconv.Itoa(id),
			"last_scanned": last_scanned,
			"port_number":  strconv.Itoa(port_number),
			"port_status":  port_status,
		}

		if err != nil {
			logger.Println(err)
		}
		logger.Println(id, host_id, last_scanned, port_number, port_status)
	}
	err = rows.Err()
	if err != nil {
		logger.Println(err)
	}
	return m, err
}

func UpdateScans(db *sql.DB,
	values url.Values,
	port_status map[string]map[string]string,
	previous_scan map[string]map[string]string) (map[string]map[string]string, error) {

	var err error

	logger, err := logger.Load("DEBUG")
	if err != nil {
		log.Fatal(err)
	}

	if values.Get("ipv4") != "" {
		// Insert host if not already in DB
		stmt, err := db.Prepare("insert ignore into hosts(ipv4) values (INET_ATON(?))")
		if err != nil {
			logger.Println(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(values.Get("ipv4"))
		if err != nil {
			logger.Println(err)
		}

		// Delete previous scan rows
		stmt, err = db.Prepare("delete from scans " +
			"where host_id = (select id from hosts where ipv4 = INET_ATON(?))")
		_, err = stmt.Exec(values.Get("ipv4"))
		defer stmt.Close()
		if err != nil {
			logger.Println(err)
		}

		// Iterate port_status map from nmap return
		for port, status := range port_status[values.Get("ipv4")] {
			stmt, err = db.Prepare("insert into scans(host_id, " +
				"port_number, " +
				"port_status) " +
				"values ((select id from hosts " +
				"where ipv4 = INET_ATON(?)), ?, ?)")
			defer stmt.Close()
			_, err = stmt.Exec(values.Get("ipv4"), port, status)
			if err != nil {
				logger.Println(err)
			}
		}
	}
	// merge port_status with previous_scan map for easy viewing of previous scan results
	for k, v := range port_status {
		previous_scan[k] = v
	}
	return previous_scan, err
}
