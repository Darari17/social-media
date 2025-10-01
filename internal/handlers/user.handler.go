package handlers

import (
	"log"
	"net/http"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repos.UserRepo
}

func NewUserHandler(ur *repos.UserRepo) *UserHandler {
	return &UserHandler{
		userRepo: ur,
	}
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.userRepo.GetAllUsers(c.Request.Context())
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

func (uh *UserHandler) GetUserByID(c *gin.Context) {
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

	user, err := uh.userRepo.GetUserByID(c.Request.Context(), userId)
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

func (uh *UserHandler) UpdateUser(c *gin.Context) {
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

	if err := uh.userRepo.UpdateUser(
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
