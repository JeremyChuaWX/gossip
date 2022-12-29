package handlers

import (
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createPostInput struct {
	UserID string `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

type updatePostInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostHandler struct {
	DB *gorm.DB
}

func (h PostHandler) CreatePost(c *fiber.Ctx) error {
	var err error
	var input createPostInput

	// input validation
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	post := models.Post{
		UserID: input.UserID,
		Title:  input.Title,
		Body:   input.Body,
	}

	// create post
	if err = h.DB.Create(&post).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": post})
}

func (h PostHandler) GetAllPosts(c *fiber.Ctx) error {
	var err error
	var posts []models.Post

	// get all posts
	if err = h.DB.Preload(clause.Associations).Find(&posts).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Posts not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": posts})
}

func (h PostHandler) GetPostById(c *fiber.Ctx) error {
	var err error
	var post models.Post
	id := c.Params("id")

	// get post by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&post).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Post not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": post})
}

func (h PostHandler) UpdatePost(c *fiber.Ctx) error {
	var err error
	var input updatePostInput
	var post models.Post
	id := c.Params("id")

	// get post by id
	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Post not found")
	}

	// input validation
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	updatePost := models.Post{
		Title: input.Title,
		Body:  input.Body,
	}

	// update post
	if err = h.DB.Model(&post).Updates(updatePost).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": post})
}

func (h PostHandler) DeletePost(c *fiber.Ctx) error {
	var err error
	var post models.Post
	id := c.Params("id")

	// get post by id
	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Post not found")
	}

	// delete post
	if err = h.DB.Delete(&post).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": true})
}
