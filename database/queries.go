package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
)

func QueryScans(db *sql.DB, values url.Values) (string, error) {
	var sql_str string
	var r string
	var err error

	var (
		id           int
		host_id      int
		last_scanned string
		port_number  int
		port_status  string
	)

	if values.Get("id") != "" {
		sql_str = fmt.Sprintf("select id, host_id, last_scanned, port_number, port_status from scans where id = %s", values.Get("id"))
	}
	if values.Get("ipv4") != "" {
		sql_str = fmt.Sprintf("select s.id, s.host_id, s.last_scanned, s.port_number, s.port_status from hosts h inner join scans s on h.id = s.host_id where h.id = (select id from hosts where ipv4 = INET_ATON('%s'))", values.Get("ipv4"))
	}

	if sql_str != "" {
		rows, err := db.Query(sql_str)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &host_id, &last_scanned, &port_number, &port_status)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(id, host_id, last_scanned, port_number, port_status)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	if port_status != "" {
		r = fmt.Sprintf("%s %s %s %s %s", id, host_id, last_scanned, port_number, port_status)
	}
	return r, err
}
