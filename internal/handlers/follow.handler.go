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

// POST /follow/:id
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

// DELETE /follow/:id
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

// GET /users/:id/followers
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

// GET /users/:id/following
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
