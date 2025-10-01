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

type CommentHandler struct {
	repo     *repos.CommentRepo
	postRepo *repos.PostRepo
}

func NewCommentHandler(r *repos.CommentRepo, p *repos.PostRepo) *CommentHandler {
	return &CommentHandler{repo: r, postRepo: p}
}

// CreateComment godoc
// @Summary Post comment
// @Description Create a new comment on a post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body dtos.CommentRequest true "Comment body"
// @Security BearerAuth
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Router /posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post id",
		})
		return
	}

	var body dtos.CommentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid body",
		})
		return
	}

	comment := models.Comment{
		UserID:  userId,
		PostID:  postId,
		Content: body.Content,
	}

	if err := h.repo.CreateComment(c.Request.Context(), &comment); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Comment added",
	})
}

// GetComments godoc
// @Summary Get comments by post ID
// @Description Get all comments for a post
// @Tags Comments
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=[]models.Comment}
// @Failure 400 {object} dtos.Response
// @Router /posts/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post id",
		})
		return
	}

	comments, err := h.repo.GetCommentsByPost(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to get comments",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get comments successfully",
		Data:    comments,
	})
}

// UpdateComment godoc
// @Summary Update comment
// @Description Update a comment by ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body dtos.CommentUpdateRequest true "Comment update body"
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Router /posts/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid comment id",
		})
		return
	}

	var body dtos.CommentUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid body",
		})
		return
	}

	if err := h.repo.UpdateComment(c.Request.Context(), commentId, body.Content); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update comment",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Comment updated",
	})
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Delete a comment by ID
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Router /posts/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid comment id",
		})
		return
	}

	if err := h.repo.DeleteComment(c.Request.Context(), commentId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to delete comment",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Comment deleted",
	})
}
