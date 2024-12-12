package model

import "github.com/Fufuhu/go-clamav/internal/queue/clients"

type ScanResult struct {
	clients.QueueMessage

	ScanResult string `json:"scan_result"`
	ScannedAt  string `json:"scanned_at"`
	Message    string `json:"message"`
}

const ScanResultInfected = "infected"
const ScanResultClean = "clean"
