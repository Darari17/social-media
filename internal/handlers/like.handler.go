package handlers

import (
	"net/http"
	"strconv"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/internal/utils"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeRepo *repos.LikeRepo
	postRepo *repos.PostRepo
}

func NewLikeHandler(r *repos.LikeRepo, p *repos.PostRepo) *LikeHandler {
	return &LikeHandler{
		likeRepo: r,
		postRepo: p,
	}
}

// POST /posts/:id/like
func (h *LikeHandler) LikePost(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	like := models.Like{
		UserID: userId,
		PostID: postId,
	}

	if err := h.likeRepo.CreateLike(c.Request.Context(), &like); err != nil {
		c.JSON(http.StatusConflict, dtos.Response{
			Code:    http.StatusConflict,
			Success: false,
			Message: "You already liked this post",
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Post liked",
	})
}

// DELETE /posts/:id/like
func (h *LikeHandler) UnlikePost(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.likeRepo.DeleteLike(c.Request.Context(), userId, postId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to unlike post",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Post unliked",
	})
}

// GET /posts/:id/likes
func (h *LikeHandler) GetLikes(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	users, err := h.likeRepo.GetLikesByPost(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to get likes",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get likes successfully",
		Data:    users,
	})
}
