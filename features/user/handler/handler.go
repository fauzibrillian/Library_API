package handler

import (
	"context"
	"errors"
	"library_api/features/user"
	"library_api/helper/cld"
	"library_api/helper/jwt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	gojwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type UserController struct {
	srv    user.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s user.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) user.Handler {
	return &UserController{
		srv:    s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

// Login implements user.Handler.
func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := uc.srv.Login(input.Phone, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]any{
					"message": "phone tidak ditemukan",
				})
			}
			if strings.Contains(err.Error(), "password salah") {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "password salah",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "phone tidak di temukan",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID, result.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var response = new(LoginResponse)
		response.ID = result.ID
		response.Name = result.Name
		response.Phone = result.Phone
		response.Role = result.Role
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Success Login Data",
			"data":    response,
		})
	}
}

// Register implements user.Handler.
func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input tidak sesuai",
			})
		}
		var inputProses = new(user.User)
		inputProses.Phone = input.Phone
		inputProses.Name = input.Name
		inputProses.Password = input.Password

		result, err := uc.srv.Register(*inputProses)
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
		var response = new(RegisterResponse)
		response.ID = result.ID
		response.Phone = result.Phone
		response.Name = result.Name

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Register Data",
			"data":    response,
		})
	}
}

// ResetPassword implements user.Handler.
func (uc *UserController) ResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(ResetPasswordRequest)
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Harap Login dulu",
				"data":    nil,
			})
		}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
				"data":    nil,
			})
		}

		var inputProcess = new(user.User)
		inputProcess.ID = uint(userID)
		inputProcess.Password = input.Password
		inputProcess.NewPassword = input.NewPassword

		result, err := uc.srv.ResetPassword(c.Get("user").(*gojwt.Token), *inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "id tidak cocok") {
				statusCode = http.StatusUnauthorized
				message = "Tidak Mempunyai Akses"
			}
			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			}
			if strings.Contains(err.Error(), "yang lama") {
				statusCode = http.StatusBadRequest
				message = "harap masukkan password yang lama jika ingin mengganti password"
			}

			return c.JSON(statusCode, map[string]any{
				"message": message,
			})
		}

		var response = new(ResetPasswordResponse)
		response.ID = result.ID
		response.Name = result.Name

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Reset Password Data User",
			"data":    response,
		})
	}
}

// UpdateUser implements user.Handler.
func (uc *UserController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(PutUserRequest)
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Harap Login dulu",
				"data":    nil,
			})
		}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
				"data":    nil,
			})
		}

		formHeader, err := c.FormFile("avatar")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": "ID user tidak valid",
					})
				}
				var inputProcess = new(user.User)
				inputProcess.Avatar = ""
				inputProcess.ID = uint(userID)
				inputProcess.Phone = input.Phone
				inputProcess.Name = input.Name

				result, err := uc.srv.UpdateUser(c.Get("user").(*gojwt.Token), *inputProcess)

				if err != nil {
					c.Logger().Error("ERROR Register, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "kesalahan pada database") {
						statusCode = http.StatusNotFound
						message = "User Tidak Terdaftar"
					}
					if strings.Contains(err.Error(), "id tidak cocok") {
						statusCode = http.StatusUnauthorized
						message = "Tidak Mempunyai Akses"
					}
					if strings.Contains(err.Error(), "terdaftar") {
						statusCode = http.StatusConflict
						message = "data yang diinputkan sudah terdaftar ada sistem"
					}
					if strings.Contains(err.Error(), "yang lama") {
						statusCode = http.StatusBadRequest
						message = "harap masukkan password yang lama jika ingin mengganti password"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(PutUserResponse)
				response.ID = result.ID
				response.Name = result.Name
				response.Phone = result.Phone
				response.Avatar = result.Avatar

				return c.JSON(http.StatusCreated, map[string]any{
					"message": "Success Updated Data User",
					"data":    response,
				})
			}
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formheader error",
					"data":    nil,
				})

		}

		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formfile error",
					"data":    nil,
				})
		}
		defer formFile.Close()

		link, err := cld.UploadImage(uc.cl, uc.ct, formFile, uc.folder)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "harap pilih gambar",
					"data":    nil,
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"message": "kesalahan pada server",
					"data":    nil,
				})
			}
		}

		var inputProcess = new(user.User)
		inputProcess.Avatar = link
		inputProcess.ID = uint(userID)
		inputProcess.Phone = input.Phone
		inputProcess.Name = input.Name

		result, err := uc.srv.UpdateUser(c.Get("user").(*gojwt.Token), *inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "id tidak cocok") {
				statusCode = http.StatusUnauthorized
				message = "Tidak Mempunyai Akses"
			}
			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			}
			if strings.Contains(err.Error(), "yang lama") {
				statusCode = http.StatusBadRequest
				message = "harap masukkan password yang lama jika ingin mengganti password"
			}

			return c.JSON(statusCode, map[string]any{
				"message": message,
			})
		}

		var response = new(PutUserResponse)
		response.ID = result.ID
		response.Name = result.Name
		response.Phone = result.Phone
		response.Avatar = result.Avatar

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Updated Data User",
			"data":    response,
		})
	}
}

// Delete implements user.Handler.
func (uc *UserController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
			})
		}

		err = uc.srv.DeleteUser(c.Get("user").(*gojwt.Token), uint(userID))
		if err != nil {
			c.Logger().Error("ERROR Delete User, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika menghapus user"

			if strings.Contains(err.Error(), "tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "user tidak ditemukan"
			} else if strings.Contains(err.Error(), "tidak memiliki izin") {
				statusCode = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk menghapus user ini"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success Deleted Data User",
		})
	}
}
