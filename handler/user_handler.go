package handler

import (
	"Steril-App/internal/service"
	"Steril-App/model"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	req := new(model.CreateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"Message": "Request body tidak sesuai",
		})
	}
	err = h.Service.CreateUser(req)
	if err != nil {

		if errors.Is(err, service.ErrUserAlreadyExists) {
			log.Printf("Handler: Konflik Bisnis: %v", err)
			return c.JSON(http.StatusConflict, echo.Map{
				"message": "Gagal membuat user: User sudah terdaftar",
			})
		}
		log.Printf("Handler: Internal Error dari Service: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal memproses pendaftaran user",
			"detail":  "Terjadi kesalahan internal saat mengakses database.",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User berhasil dibuat",
		"nik":     req.NIK,
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {

	id := c.Param("id")
	fmt.Println(id)
	err := h.Service.DeleteUser(id)

	if err != nil {
		if err.Error() == fmt.Sprintf("user dengan ID %s tidak ditemukan", id) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal menghapus pengguna",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Pengguna berhasil dihapus",
	})
}

func (h *UserHandler) GetAllUser(c echo.Context) error {

	rows, err := h.Service.GetAllUser()
	if err != nil {
		return c.JSON(http.StatusBadGateway, "error")
	}
	return c.JSON(http.StatusOK, rows)
}
