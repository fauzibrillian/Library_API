package main

import (
	"library_api/config"
	"library_api/database"
	uh "library_api/features/user/handler"
	ur "library_api/features/user/repository"
	us "library_api/features/user/service"
	"library_api/routes"

	ek "library_api/helper/enkrip"

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

	db.AutoMigrate(&ur.UserModel{})
	ekrip := ek.New()
	userRepo := ur.New(db)
	userService := us.New(userRepo, ekrip)
	userHandler := uh.New(userService)
	routes.InitRoute(e, userHandler)

	e.Logger.Fatal(e.Start(":8000"))

}
