package handlersensor

import (
	"Steril-App/model"
	"Steril-App/ws"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteFingerByID(c echo.Context) error {
	payload := model.ScanCommand{}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Format data salah")
	}

	// 2. Set Perintah
	payload.Command = "DELETE"
	fmt.Println(payload)

	// 3. Kirim
	if err := ws.SendCommand(payload); err != nil {
		return err
	}

	return nil
}
