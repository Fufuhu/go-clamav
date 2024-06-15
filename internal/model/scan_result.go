package model

import "github.com/Fufuhu/go-clamav/internal/queue/clients"

type ScanResult struct {
	clients.S3Object

	ScanResult string `json:"scan_result"`
	ScannedAt  string `json:"scanned_at"`
}

const ScanResultInfected = "infected"
const ScanResultClean = "clean"
