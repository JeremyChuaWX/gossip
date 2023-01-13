package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostHandler struct {
	DB *gorm.DB
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	type createPostInput struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	var err error
	var input createPostInput

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ServerResponse{
			Error: true,
			Msg:   "Invalid input",
			Data:  errors,
		})
	}

	post := models.Post{
		UserID: currId,
		Title:  input.Title,
		Body:   input.Body,
	}

	// create post
	if err = h.DB.Create(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post created",
		Data:  post,
	})
}

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	var err error
	var posts []models.Post

	// get all posts
	if err = h.DB.Preload(clause.Associations).Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Posts not found")
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Posts found",
		Data:  posts,
	})
}

func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	var err error
	var post models.Post
	id := c.Params("id")

	// get post by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post found",
		Data:  post,
	})
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	type updatePostInput struct {
		Title string `json:"title,omitempty"`
		Body  string `json:"body,omitempty"`
	}

	var err error
	var input updatePostInput
	var post models.Post
	id := c.Params("id")

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get post by id
	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	// check authorised
	if currId != post.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ServerResponse{
			Error: true,
			Msg:   "Invalid input",
			Data:  errors,
		})
	}

	updatePost := models.Post{
		Title: input.Title,
		Body:  input.Body,
	}

	// update post
	if err = h.DB.Model(&post).Updates(updatePost).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post updated",
		Data:  post,
	})
}

func (h *PostHandler) UpdatePostScore(c *fiber.Ctx) error {
	type updatePostScoreInput struct {
		PostScore int `json:"post_score"`
	}

	var err error
	var input updatePostScoreInput
	var post models.Post
	id := c.Params("id")

	// get post by id
	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ServerResponse{
			Error: true,
			Msg:   "Invalid input",
			Data:  errors,
		})
	}

	updatePost := models.Post{PostScore: input.PostScore}

	// update post score
	if err = h.DB.Model(&post).Select("post_score").Updates(updatePost).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post score updated",
		Data:  post,
	})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	var err error
	var post models.Post
	id := c.Params("id")

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get post by id
	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	// check authorised
	if currId != post.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
	}

	// delete post
	if err = h.DB.Delete(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post deleted",
		Data:  post,
	})
}
