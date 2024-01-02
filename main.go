package main

import (
	"library_api/config"
	"library_api/database"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	if cfg == nil {
		e.Logger.Fatal("tidak bisa start server kesalahan database")
	}

	db, err := database.InitMySql(*cfg)
	if err != nil {
		e.Logger.Fatal("tidak bisa start bro", err.Error())
	}

	db.AutoMigrate()
	e.Logger.Fatal(e.Start(":8000"))

}
