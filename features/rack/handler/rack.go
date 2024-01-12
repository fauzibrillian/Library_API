package handler

import (
	"library_api/features/rack"
	"net/http"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type RackHandler struct {
	r rack.Service
}

func New(r rack.Service) rack.Handler {
	return &RackHandler{
		r: r,
	}
}

// AddRack implements rack.Handler.
func (rh *RackHandler) AddRack() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RackRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input tidak sesuai",
			})
		}

		var inputProses = new(rack.Rack)
		inputProses.Name = input.Name

		result, err := rh.r.AddRack(c.Get("user").(*golangjwt.Token), *inputProses)
		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "email telah terdaftar",
			})
		}

		var response = new(RackResponse)
		response.ID = result.ID
		response.Name = result.Name

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Created Data",
			"data":    response,
		})
	}
}
