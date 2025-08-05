package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/labstack/gommon/log"
)

type Handler interface {
	ProcessMessage(ctx context.Context, scan scanning.Scan) error
}

type handler struct {
	dbClient *sql.DB
}

func NewHandler(dbClient *sql.DB) *handler {
	return &handler{
		dbClient: dbClient,
	}
}
func (h *handler) Close() error {
	if h.dbClient != nil {
		return h.dbClient.Close()
	}
	return nil
}

func (h *handler) ProcessMessage(ctx context.Context, scan scanning.Scan) error {

	// Check if previous DB entry exists and is older than the current scan.
	lastScanned, err := h.getLastScanned(scan)
	if err != nil {
		return fmt.Errorf("failed to get last scanned time: %w", err)
	}
	if scan.Timestamp < lastScanned {
		log.Debug("Skipping scan for", scan.Ip, "port", scan.Port, "service", scan.Service, "as it is older than the last recorded scan.")
		return nil
	}

	dataBytes, err := getNormalizedData(&scan)
	if err != nil {
		return fmt.Errorf("failed to get normalized data: %w", err)
	}

	// Save to db
	_, err = h.writeScan(scan, string(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to insert scan into database: %w", err)
	}
	return nil
}

func (h *handler) writeScan(scan scanning.Scan, dataBytes string) (sql.Result, error) {
	stmt := `INSERT OR REPLACE INTO scans (ip, port, service, last_scanned, response) VALUES (?, ?, ?, ?, ?)`
	return h.dbClient.Exec(stmt, scan.Ip, scan.Port, scan.Service, scan.Timestamp, dataBytes)
}

func (h *handler) getLastScanned(scan scanning.Scan) (int64, error) {
	var lastScanned int64
	stmt := `SELECT last_scanned FROM scans WHERE ip = ? AND port = ? AND service = ?`
	err := h.dbClient.QueryRow(stmt, scan.Ip, scan.Port, scan.Service).Scan(&lastScanned)
	if err != nil {
		if err == sql.ErrNoRows {
			return lastScanned, nil
		}
		return lastScanned, fmt.Errorf("failed to get last scanned time: %w", err)
	}
	return lastScanned, nil
}

func getNormalizedData(scan *scanning.Scan) (string, error) {
	if scan == nil {
		return "", fmt.Errorf("recieved empty scan")
	}

	// Use reflection to handle different data versions
	switch scan.DataVersion {
	case scanning.V1:
		v2Data, err := convertV1ToV2(scan.Data)
		if err != nil {
			return "", fmt.Errorf("failed to convert V1 data to V2: %w", err)
		}
		if v2Data.ResponseStr == "" {
			return "", fmt.Errorf("V1 data conversion resulted in empty response")
		}
		return v2Data.ResponseStr, nil
	case scanning.V2:
		var v2Data scanning.V2Data
		dataBytes, _ := json.Marshal(scan.Data)
		if err := json.Unmarshal(dataBytes, &v2Data); err != nil {
			return "", fmt.Errorf("failed to unmarshal V2 data: %w", err)
		}
		if v2Data.ResponseStr == "" {
			return "", fmt.Errorf("V2 data resulted in empty response")
		}
		return v2Data.ResponseStr, nil
	}

	return "", fmt.Errorf("unknown DataVersion")
}

func convertV1ToV2(data any) (scanning.V2Data, error) {
	var v1Data scanning.V1Data
	if data == nil {
		return scanning.V2Data{}, fmt.Errorf("data is nil")
	}

	dataBytes, _ := json.Marshal(data)
	if err := json.Unmarshal(dataBytes, &v1Data); err != nil {
		return scanning.V2Data{}, err
	}

	return scanning.V2Data{
		ResponseStr: string(v1Data.ResponseBytesUtf8),
	}, nil
}
