package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"encoding/json"
)

func QueryScans(db *sql.DB, values url.Values) (interface{}, error) {
	var sql_str string
	var b []byte
	var r interface{}
	var err error

	var (
		id           int
		host_id      int
		last_scanned string
		port_number  int
		port_status  string
	)

	if values.Get("id") != "" {
		sql_str = fmt.Sprintf("select id, " +
		                      "host_id, " +
				      "last_scanned, " +
				      "port_number, " +
				      "port_status from scans where id = %s", values.Get("id"))
	}
	if values.Get("ipv4") != "" {
		sql_str = fmt.Sprintf("select s.id, " +
		                      "s.host_id, " +
				      "s.last_scanned, " +
				      "s.port_number, " +
				      "s.port_status from hosts h " +
				      "inner join scans s on h.id = s.host_id " +
				      "where h.id = (select id from hosts where ipv4 = INET_ATON('%s')" +
			              ")", values.Get("ipv4"))
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
		json_str := fmt.Sprintf("{\"id\": %d, " +
		                        "\"host_id\": %d, " +
					"\"last_scanned\": \"%s\", " +
					"\"port_number\": %d, " +
					"\"port_status\": \"%s\"}", id, host_id, last_scanned, port_number, port_status)
		b = []byte(json_str)
		err = json.Unmarshal(b, &r)
	}
	return r, err
}

func UpdateScans(db *sql.DB, values url.Values, port_status map[string]string) (map[string]string, error) {
        var sql_str string
	var err error
        if values.Get("ipv4") != "" {
		// Insert host if not already in DB
		sql_str = fmt.Sprintf("insert ignore into hosts(ipv4) values (INET_ATON('%s'))", values.Get("ipv4"))
                rows, err := db.Query(sql_str)
                if err != nil {
                        log.Fatal(err)
                }
                defer rows.Close()

		// Delete previous scan rows
		sql_str = fmt.Sprintf("delete from scans " +
		                      "where host_id = (select id from hosts where ipv4 = INET_ATON('%s')" +
			              ")", values.Get("ipv4"))
                rows, err = db.Query(sql_str)
                if err != nil {
                        log.Fatal(err)
                }
                defer rows.Close()

		// Iterate port_status map from nmap return
		for port, status := range port_status {
			sql_str = fmt.Sprintf("insert into scans(host_id, " +
			                      "port_number, " +
					      "port_status) " +
					      "values ((select id from hosts " +
					      "where ipv4 = INET_ATON('%s')), %s, '%s'" +
				              ")", values.Get("ipv4"), port, status)
			rows, err = db.Query(sql_str)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	return port_status, err
}
