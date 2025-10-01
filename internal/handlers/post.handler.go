package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/Darari17/social-media/internal/repos"
	"github.com/Darari17/social-media/internal/utils"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postRepo *repos.PostRepo
}

func NewPostHandler(postRepo *repos.PostRepo) *PostHandler {
	return &PostHandler{postRepo: postRepo}
}

// CreatePost godoc
// @Summary Create post
// @Description Create a new post
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param content formData string false "Post content"
// @Param image formData file false "Post image"
// @Security BearerAuth
// @Success 201 {object} dtos.Response{data=dtos.PostResponse}
// @Failure 400 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Router /posts [post]
func (ph *PostHandler) CreatePost(c *gin.Context) {
	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var body dtos.PostRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid form data",
		})
		return
	}

	var imagePath *string
	if body.Image != nil {
		if filename, err := utils.FileUpload(c, body.Image, "posts"); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, dtos.Response{
				Code:    http.StatusBadRequest,
				Success: false,
				Message: err.Error(),
			})
			return
		} else {
			imagePath = &filename
		}
	}

	post := models.Post{
		UserID:  userId,
		Content: &body.Content,
		Image:   imagePath,
	}

	if err := ph.postRepo.CreatePost(c.Request.Context(), &post); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to create post",
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Post created successfully",
		Data: dtos.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			DeletedAt: post.DeletedAt,
		},
	})
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Get list of all posts
// @Tags Posts
// @Produce json
// @Success 200 {object} dtos.Response{data=[]dtos.PostResponse}
// @Router /posts [get]
func (ph *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := ph.postRepo.GetAllPosts(c.Request.Context())
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch posts",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get posts successfully",
		Data:    posts,
	})
}

// GetPostByID godoc
// @Summary Get post by ID
// @Description Get a single post by its ID
// @Tags Posts
// @Produce json
// @Param postId path int true "Post ID"
// @Success 200 {object} dtos.Response{data=dtos.PostResponse}
// @Failure 404 {object} dtos.Response
// @Router /posts/{postId} [get]
func (ph *PostHandler) GetPostByID(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post id",
		})
		return
	}

	post, err := ph.postRepo.GetPostByID(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get post successfully",
		Data: dtos.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	})
}

// UpdatePost godoc
// @Summary Update post
// @Description Update a post by ID
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param postId path int true "Post ID"
// @Param content formData string false "Post content"
// @Param image formData file false "Post image"
// @Security BearerAuth
// @Success 200 {object} dtos.Response{data=dtos.PostResponse}
// @Failure 400 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Router /posts/{postId} [patch]
func (ph *PostHandler) UpdatePost(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post id",
		})
		return
	}

	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	existingPost, err := ph.postRepo.GetPostByID(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Post not found",
		})
		return
	}

	if existingPost.UserID != userId {
		c.JSON(http.StatusForbidden, dtos.Response{
			Code:    http.StatusForbidden,
			Success: false,
			Message: "You are not allowed to update this post",
		})
		return
	}

	var body dtos.PostUpdateRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid form data",
		})
		return
	}

	var imagePath *string
	if body.Image != nil {
		filename, err := utils.FileUpload(c, body.Image, "posts")
		if err != nil {
			c.JSON(http.StatusBadRequest, dtos.Response{
				Code:    http.StatusBadRequest,
				Success: false,
				Message: err.Error(),
			})
			return
		}
		imagePath = &filename
	}

	updated := models.Post{ID: postId}

	if body.Content != nil && *body.Content != "" {
		updated.Content = body.Content
	}
	if imagePath != nil {
		updated.Image = imagePath
	}

	if updated.Content == nil && updated.Image == nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "No fields to update",
		})
		return
	}

	if err := ph.postRepo.UpdatePost(c.Request.Context(), &updated); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update post",
		})
		return
	}

	postAfterUpdate, _ := ph.postRepo.GetPostByID(c.Request.Context(), postId)

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Post updated successfully",
		Data:    postAfterUpdate,
	})
}

// DeletePost godoc
// @Summary Delete post
// @Description Delete a post by ID
// @Tags Posts
// @Produce json
// @Param postId path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Router /posts/{postId} [delete]
func (ph *PostHandler) DeletePost(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid post id",
		})
		return
	}

	userId, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	existingPost, err := ph.postRepo.GetPostByID(c.Request.Context(), postId)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Post not found",
		})
		return
	}

	if existingPost.UserID != userId {
		c.JSON(http.StatusForbidden, dtos.Response{
			Code:    http.StatusForbidden,
			Success: false,
			Message: "You are not allowed to delete this post",
		})
		return
	}

	if err := ph.postRepo.DeletePost(c.Request.Context(), postId); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Post deleted successfully",
	})
}
