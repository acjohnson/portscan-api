portscan-api
============

**Building**

Clone this repo and build using golang 1.11+ (uses go modules)
```sh
go build
```

Install mysql/mariadb and create database:

```sql
create database `portscan-api` character set UTF8 collate utf8_bin;
grant all privileges on `portscan-api`.* to 'portscan-api'@'%' identified by "<secret>";
flush privileges;
```

Table creation is handled automatically by the database package

**Routes**

`/scans`

- Run a port scan send a PUT to `/scans` with URL value ipv4 set to the IP of the host to scan (eg. using httpie)
```sh
http PUT http://127.0.0.1:10000/scans?ipv4='127.0.0.1'
```

Should return http status code `202` and json
```json
{
    "111": "open",
    "139": "open",
    "22": "open"
}
```

- Retrieving a scan by its id
```sh
http GET http://127.0.0.1:10000/scans?id=11
```

Should return http status code `200` and json
```json
{
    "host_id": 1,
    "id": 11,
    "last_scanned": "2020-04-23 14:26:29",
    "port_number": 80,
    "port_status": "open"
}
```
