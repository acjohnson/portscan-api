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
	var sql_str string
	var err error

	logger, err := logger.Load()
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

	sql_str = fmt.Sprintf("select id, " +
		"host_id, " +
		"last_scanned, " +
		"port_number, " +
		"port_status from scans")

	if values.Get("id") != "" {
		sql_str = fmt.Sprintf("select id, "+
			"host_id, "+
			"last_scanned, "+
			"port_number, "+
			"port_status from scans where id = %s", values.Get("id"))
	}
	if values.Get("ipv4") != "" {
		sql_str = fmt.Sprintf("select s.id, "+
			"s.host_id, "+
			"s.last_scanned, "+
			"s.port_number, "+
			"s.port_status from hosts h "+
			"inner join scans s on h.id = s.host_id "+
			"where h.id = (select id from hosts where ipv4 = INET_ATON('%s')"+
			")", values.Get("ipv4"))
	}

	if sql_str != "" {
		rows, err := db.Query(sql_str)

		if err != nil {
			logger.Println(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&id, &host_id, &last_scanned, &port_number, &port_status); err != nil {
				return m, err
			}

			m[strconv.Itoa(port_number)] = map[string]string{
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
	}
	return m, err
}

func UpdateScans(db *sql.DB,
	values url.Values,
	port_status map[string]map[string]string,
	previous_scan map[string]map[string]string) (map[string]map[string]string, error) {

	var sql_str string
	var err error

	logger, err := logger.Load()
	if err != nil {
		log.Fatal(err)
	}

	if values.Get("ipv4") != "" {
		// Insert host if not already in DB
		sql_str = fmt.Sprintf("insert ignore into hosts(ipv4) values (INET_ATON('%s'))", values.Get("ipv4"))
		rows, err := db.Query(sql_str)
		if err != nil {
			logger.Println(err)
		}
		defer rows.Close()

		// Delete previous scan rows
		sql_str = fmt.Sprintf("delete from scans "+
			"where host_id = (select id from hosts where ipv4 = INET_ATON('%s')"+
			")", values.Get("ipv4"))
		rows, err = db.Query(sql_str)
		if err != nil {
			logger.Println(err)
		}
		defer rows.Close()

		// Iterate port_status map from nmap return
		for port, status := range port_status[values.Get("ipv4")] {
			sql_str = fmt.Sprintf("insert into scans(host_id, "+
				"port_number, "+
				"port_status) "+
				"values ((select id from hosts "+
				"where ipv4 = INET_ATON('%s')), %s, '%s'"+
				")", values.Get("ipv4"), port, status)
			rows, err = db.Query(sql_str)
			if err != nil {
				logger.Println(err)
			}
			defer rows.Close()

			err = rows.Err()
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
