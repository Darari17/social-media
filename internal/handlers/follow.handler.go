package handlers

import (
	"net/http"
	"strconv"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/internal/utils"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followRepo *repos.FollowRepo
}

func NewFollowHandler(repo *repos.FollowRepo) *FollowHandler {
	return &FollowHandler{followRepo: repo}
}

// FollowUser godoc
// @Summary Follow user
// @Description Follow another user by ID
// @Tags Follow
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID to follow"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Router /follow/{id} [post]
func (fh *FollowHandler) FollowUser(c *gin.Context) {
	followerId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	followingIdStr := c.Param("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid user id",
		})
		return
	}

	if followerId == followingId {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "You cannot follow yourself",
		})
		return
	}

	res, err := fh.followRepo.FollowUser(c.Request.Context(), followerId, followingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to follow user",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Followed user successfully",
		Data:    res,
	})
}

// UnfollowUser godoc
// @Summary Unfollow user
// @Description Unfollow another user by ID
// @Tags Follow
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID to unfollow"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Router /follow/{id} [delete]
func (fh *FollowHandler) UnfollowUser(c *gin.Context) {
	followerId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	followingIdStr := c.Param("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid user id",
		})
		return
	}

	rows, err := fh.followRepo.UnfollowUser(c.Request.Context(), followerId, followingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to unfollow user",
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Follow relation not found",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Unfollowed user successfully",
	})
}

// GetFollowers godoc
// @Summary Get followers
// @Description Get list of followers for a user
// @Tags Follow
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.Response{data=[]models.User}
// @Failure 400 {object} dtos.Response
// @Router /users/{id}/followers [get]
func (fh *FollowHandler) GetFollowers(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid user id",
		})
		return
	}

	users, err := fh.followRepo.GetFollowers(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch followers",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get followers successfully",
		Data:    users,
	})
}

// GetFollowing godoc
// @Summary Get following
// @Description Get list of users that a user is following
// @Tags Follow
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.Response{data=[]models.User}
// @Failure 400 {object} dtos.Response
// @Router /users/{id}/following [get]
func (fh *FollowHandler) GetFollowing(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid user id",
		})
		return
	}

	users, err := fh.followRepo.GetFollowing(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch following",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get following successfully",
		Data:    users,
	})
}
