portscan-api
============

**Build**

Clone this repo and build using golang 1.11+ (uses go modules)
```sh
go build
```

Install `nmap` and `mysql` or `mariadb` using your favorite package manager
and create the database

```sql
create database `portscan-api` character set UTF8 collate utf8_bin;
create user 'portscan-api'@'%' identified by '<secret>';
grant all privileges on `portscan-api`.* to 'portscan-api'@'%';
flush privileges;
```

Table creation is handled automatically by the database package

**Configuration**

Create a configuration file called `conf.json` in the root of this project.

conf.json.example has been provided to help you get started.

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

- Retrieving a scan by scan id
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

- Retrieving a scan by host ip
```sh
http GET http://127.0.0.1:10000/scans?ipv4=127.0.0.1
```

Should return http status code `200` and json
```json
{
    "111": {
        "host_id": 4,
        "id": 28,
        "last_scanned": "2020-04-23 15:40:44",
        "port_number": 111,
        "port_status": "open"
    },
    "139": {
        "host_id": 4,
        "id": 29,
        "last_scanned": "2020-04-23 15:40:44",
        "port_number": 139,
        "port_status": "open"
    },
    "22": {
        "host_id": 4,
        "id": 27,
        "last_scanned": "2020-04-23 15:40:44",
        "port_number": 22,
        "port_status": "open"
    }
}
```
