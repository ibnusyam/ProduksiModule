package model

type ScanCommand struct {
	Command string `json:"Command"`
	ID      string `json:"id"`
}

type SensorResponse struct {
	Action  string `json:"action"`
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Trigger string `json:"trigger"`
}

type AddFingerRequest struct {
	Nik string `json:"nik"`
}

type FingerID struct {
	FingerID string
}
