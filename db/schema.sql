CREATE TABLE IF NOT EXISTS scans (
    ip TEXT NOT NULL,
    port INTEGER NOT NULL,
    service TEXT NOT NULL,
    last_scanned INTEGER NOT NULL,
    response TEXT,
    PRIMARY KEY (ip, port, service)
);