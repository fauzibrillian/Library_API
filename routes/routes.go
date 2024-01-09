package routes

import (
	"library_api/features/book"
	"library_api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uc user.Handler, bh book.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	RouteUser(e, uc)
	RouteBook(e, bh)
}

func RouteUser(e *echo.Echo, uc user.Handler) {
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.PATCH("/resetpassword/:id", uc.ResetPassword(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PATCH("/user/:id", uc.UpdateUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/user/:id", uc.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func RouteBook(e *echo.Echo, bh book.Handler) {
	e.POST("/addbook", bh.AddBook(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
