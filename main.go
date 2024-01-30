package main

import (
	"library_api/config"
	"library_api/database"
	uh "library_api/features/user/handler"
	ur "library_api/features/user/repository"
	us "library_api/features/user/service"

	bh "library_api/features/book/handler"
	br "library_api/features/book/repository"
	bs "library_api/features/book/service"

	th "library_api/features/transaction/handler"
	tr "library_api/features/transaction/repository"
	ts "library_api/features/transaction/service"
	"library_api/routes"

	"library_api/helper/cld"
	ek "library_api/helper/enkrip"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	if cfg == nil {
		e.Logger.Fatal("tidak bisa start server kesalahan database")
	}
	cld, ctx, param := cld.InitCloudnr(*cfg)

	db, err := database.InitMySql(*cfg)
	if err != nil {
		e.Logger.Fatal("tidak bisa start bro", err.Error())
	}

	db.AutoMigrate(&ur.UserModel{}, &br.BookModel{}, &tr.TransactionModel{})
	ekrip := ek.New()
	userRepo := ur.New(db)
	userService := us.New(userRepo, ekrip)
	userHandler := uh.New(userService, cld, ctx, param)

	bookRepo := br.New(db)
	bookService := bs.New(bookRepo)
	bookHandler := bh.New(bookService, cld, ctx, param)

	transactionRepo := tr.New(db)
	transactionService := ts.New(transactionRepo)
	transactionHandler := th.New(transactionService)

	routes.InitRoute(e, userHandler, bookHandler, transactionHandler)

	e.Logger.Fatal(e.Start(":8000"))

}
