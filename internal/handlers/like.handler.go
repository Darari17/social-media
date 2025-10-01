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

// LikePost godoc
// @Summary Like post
// @Description Like a post by ID
// @Tags Likes
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Router /posts/{id}/like [post]
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

// UnlikePost godoc
// @Summary Unlike post
// @Description Unlike a post by ID
// @Tags Likes
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Router /posts/{id}/unlike [post]
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

// GetLikes godoc
// @Summary Get likes
// @Description Get list of users who liked a post
// @Tags Likes
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=[]models.User}
// @Failure 400 {object} dtos.Response
// @Router /posts/{id}/likes [get]
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
