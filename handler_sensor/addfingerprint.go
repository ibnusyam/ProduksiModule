package handlersensor

import (
	"Steril-App/model"
	"Steril-App/ws"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddFingerByID(c echo.Context) error {
	// 1. Ambil data dari Body Request
	payload := model.ScanCommand{}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Format data salah")
	}

	// 2. Set Perintah Khusus
	payload.Command = "DAFTAR_BARU" // Sesuaikan dengan logika NodeMCU kamu

	// 3. Kirim via WebSocket Helper
	if err := ws.SendCommand(payload); err != nil {
		return err // Error sudah dihandle di ws package
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Perintah Daftar Terkirim"})
}
