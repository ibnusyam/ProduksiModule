package handlersensor

import (
	"Steril-App/model"
	"Steril-App/ws"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ScanRegisteredFinger(c echo.Context) error {
	// 1. Buat Payload
	payload := model.ScanCommand{
		Command: "SCAN",
	}

	// 2. Kirim
	if err := ws.SendCommand(payload); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Perintah Scan Terkirim"})
}
