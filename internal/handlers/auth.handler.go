package handlers

import (
	"log"
	"net/http"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/internal/utils"
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

func (ah *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := ah.authRepo.GetAllUsers(c.Request.Context())
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch all users",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get all users successfully",
		Data:    users,
	})
}

func (ah *AuthHandler) GetUserByID(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	user, err := ah.authRepo.GetUserByID(c.Request.Context(), userId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "User not found",
		})
		return
	}

	response := dtos.UserResponse{
		ID:        userId,
		Name:      user.Name,
		Email:     &user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get profile successfully",
		Data:    response,
	})
}

func (ah *AuthHandler) UpdateUser(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	var body dtos.UserUpdateRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid body request",
		})
		return
	}

	var avatarPath *string
	file, err := c.FormFile("avatar")
	if err == nil {
		if filename, err := utils.FileUpload(c, file, "avatar"); err != nil {
			c.JSON(http.StatusBadRequest, dtos.Response{
				Code:    http.StatusBadRequest,
				Success: false,
				Message: err.Error(),
			})
			return
		} else {
			avatarPath = &filename
		}
	}

	updatedUser := models.User{
		ID:     userId,
		Name:   body.Name,
		Avatar: avatarPath,
		Bio:    body.Bio,
	}

	if err := ah.authRepo.UpdateUser(
		c.Request.Context(),
		&updatedUser,
	); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Profile updated successfully",
	})
}
