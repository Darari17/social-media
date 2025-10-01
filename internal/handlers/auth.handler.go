package handlers

import (
	"log"
	"net/http"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authRepo *repos.AuthRepo
}

func NewAuthHandler(authRepo *repos.AuthRepo) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
	}
}

// Register godoc
// @Summary Register user
// @Description Register with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.UserRequest true "Register request"
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 409 {object} dtos.Response
// @Router /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var body dtos.UserRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid body request",
		})
		return
	}

	hashedPwd, err := pkg.HashPassword(body.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: hashedPwd,
	}

	if err := ah.authRepo.CreateAccount(c.Request.Context(), &user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusConflict, dtos.Response{
			Code:    http.StatusConflict,
			Success: false,
			Message: "Email already exists",
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Register successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.UserRequest true "Login request"
// @Success 200 {object} dtos.Response{data=dtos.UserTokenResponse}
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var body dtos.UserRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid body request",
		})
		return
	}

	user, err := ah.authRepo.GetEmail(c.Request.Context(), body.Email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Something went wrong",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Email or Password",
		})
		return
	}

	if ok := pkg.VerifyPassword(user.Password, body.Password); !ok {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Email or Password",
			Data:    nil,
		})
		return
	}

	claim := pkg.NewJWTClaims(user.ID)
	token, err := claim.GenerateToken()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to Generate Token",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login Succesfully",
		Data: dtos.UserTokenResponse{
			Token: token,
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout by invalidating JWT
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /auth/logout [get]
func (ah *AuthHandler) Logout(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	if err := ah.authRepo.Logout(c.Request.Context(), bearerToken); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Logout Succesfully",
	})
}
