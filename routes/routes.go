package routes

import (
	"library_api/config"
	"library_api/features/book"
	"library_api/features/transaction"
	"library_api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uc user.Handler, bh book.Handler, th transaction.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	RouteUser(e, uc)
	RouteBook(e, bh)
	RouteTransaction(e, th)
}

func RouteUser(e *echo.Echo, uc user.Handler) {
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.PATCH("/resetpassword/:id", uc.ResetPassword(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.PATCH("/users/:id", uc.UpdateUser(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.DELETE("/users/:id", uc.Delete(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.GET("/users", uc.SearchUser(), echojwt.JWT([]byte(config.InitConfig().JWT)))

}

func RouteBook(e *echo.Echo, bh book.Handler) {
	e.POST("/books", bh.AddBook(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.PATCH("/books/:id", bh.UpdateBook(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.DELETE("/books/:id", bh.DeleteBook(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.GET("/books", bh.SearchBook())
	e.GET("/books/:id", bh.GetBook())
}

func RouteTransaction(e *echo.Echo, th transaction.Handler) {
	e.POST("/transactions", th.Borrow(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.GET("/transactions", th.AllTransaction(), echojwt.JWT([]byte(config.InitConfig().JWT)))
	e.PATCH("/transactions/:id", th.UpdateReturn(), echojwt.JWT([]byte(config.InitConfig().JWT)))

}
