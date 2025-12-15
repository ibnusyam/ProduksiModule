package main

import (
	"Steril-App/handler"
	handlersensor "Steril-App/handler_sensor"
	"Steril-App/internal/repository"
	"Steril-App/internal/service"
	"Steril-App/ws"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	// Hapus import yang ini: "github.com/labstack/echo/middleware"

	// Gunakan path yang benar untuk Echo v4
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware" // <-- Sudah dikoreksi
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// Menggunakan %v untuk error
		fmt.Printf("file .env tidak ada: %v\n", err)
	}

	db, err := repository.ConnectDB()
	if err != nil {
		fmt.Println("gagal conect db", err)
		return // Tambahkan return jika gagal koneksi ke DB
	}

	addFingerRepository := repository.NewAddFingerRepository(db)
	fingerRepository := repository.NewFingerRepository(db)
	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepository, fingerRepository)
	userHandler := handler.NewUserHandler(userService)

	logFingerRepository := repository.NewFingerLogRepostory(db)
	fingerLogHandler := handler.NewLogFingerHanlere(logFingerRepository)

	wsHandler := ws.NewWebSocketHandler(addFingerRepository, fingerRepository, logFingerRepository)
	// addFingerLog := handlersensor.NewFingerLog(logFingerRepository, fingerRepository)

	// Inisialisasi Echo
	e := echo.New()

	// Pastikan CORS Middleware diaktifkan paling awal
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPut, http.MethodPost,
			http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType,
			echo.HeaderAccept, echo.HeaderAuthorization,
		},
		AllowCredentials: false,
	}))

	// Routes
	e.POST("/create", userHandler.CreateUser)
	e.DELETE("/delete/:id", userHandler.DeleteUser)

	e.GET("/users", userHandler.GetAllUser)

	e.POST("/get", fingerLogHandler.GetFingerLog)
	e.POST("insert", fingerLogHandler.AddManualFingerLog)
	e.POST("/remove", fingerLogHandler.DeleteFingerLog)
	e.POST("/notes", fingerLogHandler.SaveNote)
	e.GET("/notes", fingerLogHandler.GetNotes)

	//Sensor
	e.GET("/ws", wsHandler.HandleWebSocket)
	e.GET("/scan", handlersensor.ScanRegisteredFinger)
	e.POST("/add", handlersensor.AddFingerByID)
	e.POST("/del", handlersensor.DeleteFingerByID)

	// Jalankan server
	e.Logger.Fatal(e.Start(":8082"))
}
