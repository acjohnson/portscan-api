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

_Note: Output includes previous scan results as well as current port status under the IP address_
```json
{
    "127.0.0.1": {
        "111": "open",
        "139": "open",
        "22": "open"
    },
    "port_number_111_host_id_1": {
        "host_id": "1",
        "id": "20",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "111",
        "port_status": "open"
    },
    "port_number_139_host_id_1": {
        "host_id": "1",
        "id": "21",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "139",
        "port_status": "open"
    },
    "port_number_22_host_id_1": {
        "host_id": "1",
        "id": "19",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "22",
        "port_status": "open"
    }
}
```

- List all scans
```
http GET http://127.0.0.1:10000/scans
```

Should return http status code `200` and json
```json
{
    "port_number_111_host_id_1": {
        "host_id": "1",
        "id": "20",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "111",
        "port_status": "open"
    },
    "port_number_111_host_id_13": {
        "host_id": "13",
        "id": "17",
        "last_scanned": "2020-04-27 15:31:03",
        "port_number": "111",
        "port_status": "open"
    },
    "port_number_111_host_id_36": {
        "host_id": "36",
        "id": "22",
        "last_scanned": "2020-04-27 16:25:29",
        "port_number": "111",
        "port_status": "open"
    }
}
```

- Retrieving a scan by scan id
```sh
http GET http://127.0.0.1:10000/scans?id=21
```

Should return http status code `200` and json
```json
{
    "port_number_139_host_id_1": {
        "host_id": "1",
        "id": "21",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "139",
        "port_status": "open"
    }
}
```

- Retrieving a scan by host ip
```sh
http GET http://127.0.0.1:10000/scans?ipv4=127.0.0.1
```

Should return http status code `200` and json
```json
{
    "port_number_111_host_id_1": {
        "host_id": "1",
        "id": "20",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "111",
        "port_status": "open"
    },
    "port_number_139_host_id_1": {
        "host_id": "1",
        "id": "21",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "139",
        "port_status": "open"
    },
    "port_number_22_host_id_1": {
        "host_id": "1",
        "id": "19",
        "last_scanned": "2020-04-27 15:43:52",
        "port_number": "22",
        "port_status": "open"
    }
}
```
