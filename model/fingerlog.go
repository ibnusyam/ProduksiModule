package model

import "time"

type FingerLogResult struct {
	NIK        string      `json:"nik"`
	FullName   string      `json:"full_name"`
	Timestamps []time.Time `json:"timestamps"` // Slice yang akan diisi
}

// Structure sementara untuk memindai setiap baris dari database
type RawFingerLog struct {
	NIK       string
	FullName  string
	Timestamp time.Time
}

type FingerLogRequest struct {
	Date string `json:"date"`
}

type NoteRequest struct {
    Date string `json:"date"`
    NIK  string `json:"nik"`
    Note string `json:"note"`
}

type NoteResponse struct {
    NIK  string `json:"nik"`
    Note string `json:"note"`
}

type AddManualFingerLogRequest struct {
    NIK       string `json:"nik" form:"nik"`
    Timestamp string `json:"timestamp" form:"timestamp"` // Format: "YYYY-MM-DD HH:mm:ss"
}

type SuccessResponse struct {
    Message string `json:"message"`
    Status  int    `json:"status"`
}

type DeleteFingerLogRequest struct {
    NIK       string `json:"nik"`
    Timestamp string `json:"timestamp"` // Format: "YYYY-MM-DD HH:mm:ss"
}