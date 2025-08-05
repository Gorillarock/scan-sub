# scan-sub
scan-sub is a series of services which coordinate the scanning of ports, publishing of scan data, subscribing to and processing scanning messages, and storing unique entries into a db.

<br><br>

## Services

### scanner
A toy scanner which publishes fake network scanning data via pubsub client.

### db
A sqlite db for storing network scanning information.

### processor
Subscribes to pubsub topic and receives network scanning data messages.
Processes messages:
    - Normalizes V1 messages to V2 data type.
    - upserts "ip", "port", "service" (composite key), with "last_scanned" (timestamp int64) and service "response", into db.
    - throws out "new" scan data messages which are older than the matching unique entry in the db

### test_db_reader
Can be run alongside the scanner, db, and processor services to log out some basic database info.
Logs:
    - Unique Scan Entry count.
    - Last 10 scan entries.

<br><br>

## How to Build and run the program

### Build
```
make build
```

NOTE: Building is a pre-requisite to running the program.

### Run the program
```
make run
```

### Run the program with a test db reader
```
make run-with-test-db-reader
```

NOTE: This starts an extra service which just queries the db and logs info about the number of unique scans and the last 10 scan entries.


## To Run Unit Tests
```
make test
```
