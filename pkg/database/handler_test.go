package database

import (
	"testing"

	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/stretchr/testify/require"
)

func Test_getNormalizedData(t *testing.T) {
	tests := []struct {
		name              string
		scanInput         *scanning.Scan
		desiredScanV2Data string
		wantErr           bool
	}{
		{
			name: "Valid Scan",
			scanInput: &scanning.Scan{
				Ip:          "192.0.2.1",
				Port:        80,
				Service:     "http",
				DataVersion: scanning.V1,
				Data:        &scanning.V1Data{ResponseBytesUtf8: []byte("Hello World")},
			},
			desiredScanV2Data: "Hello World",
			wantErr:           false,
		},
		{
			name: "Valid V2 Scan",
			scanInput: &scanning.Scan{
				Ip:          "192.0.2.1",
				Port:        80,
				Service:     "http",
				DataVersion: scanning.V2,
				Data:        &scanning.V2Data{ResponseStr: "Hello World"},
			},
			desiredScanV2Data: "Hello World",
			wantErr:           false,
		},
		{
			name: "Invalid Scan Version",
			scanInput: &scanning.Scan{
				Ip:          "192.0.2.1",
				Port:        80,
				Service:     "http",
				DataVersion: 3,
				Data:        scanning.V1Data{ResponseBytesUtf8: []byte("Hello World")},
			},
			desiredScanV2Data: "Hello World",
			wantErr:           true,
		},
		{
			name: "Mismatched Data Version V1",
			scanInput: &scanning.Scan{
				Ip:          "192.0.2.1",
				Port:        80,
				Service:     "http",
				DataVersion: scanning.V1,
				Data:        scanning.V2Data{ResponseStr: "Hello World"},
			},
			desiredScanV2Data: "",
			wantErr:           true,
		},
		{
			name: "Mismatched Data Version V2",
			scanInput: &scanning.Scan{
				Ip:          "192.0.2.1",
				Port:        80,
				Service:     "http",
				DataVersion: scanning.V2,
				Data:        scanning.V1Data{ResponseBytesUtf8: []byte("Hello World")},
			},
			desiredScanV2Data: "",
			wantErr:           true,
		},
		{
			name:              "Nil Scan",
			scanInput:         nil,
			desiredScanV2Data: "",
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNormalizedData(tt.scanInput)
			if tt.wantErr {
				require.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err, "Expected no error but got one")
			if got != tt.desiredScanV2Data {
				t.Errorf("getNormalizedData() = %v, want %v", got, tt.desiredScanV2Data)
			}
		})
	}
}
