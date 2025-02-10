package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"article-golang/config"
	"article-golang/models"

	"github.com/gin-contrib/cors"
)

var validate = validator.New()

type ArticleInput struct {
	Title    string `json:"title" binding:"required,min=20"`
	Content  string `json:"content" binding:"required,min=200"`
	Category string `json:"category" binding:"required,min=3"`
	Status   string `json:"status" binding:"required,oneof=publish draft trash"`
}

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.Post{})

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/article", createArticle)
	r.GET("/articles/:limit/:offset", getArticles)
	r.GET("/article/:id", getArticleByID)
	r.PUT("/article/:id", updateArticle)
	r.DELETE("/article/:id", deleteArticle)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}

func validateArticle(article models.Post) string {
	if len(article.Title) < 20 {
		return "Title harus memiliki minimal 20 karakter"
	}
	if len(article.Content) < 200 {
		return "Content harus memiliki minimal 200 karakter"
	}
	if len(article.Category) < 3 {
		return "Category harus memiliki minimal 3 karakter"
	}
	if article.Status != "publish" && article.Status != "draft" && article.Status != "trash" {
		return "Status harus bernilai 'publish', 'draft', atau 'trash'"
	}
	return ""
}


func createArticle(c *gin.Context) {
	var article models.Post
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errMsg := validateArticle(article); errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	config.DB.Create(&article)
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

func getArticles(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	var posts []models.Post
	config.DB.Limit(limit).Offset(offset).Find(&posts)
	c.JSON(http.StatusOK, posts)
}

func getArticleByID(c *gin.Context) {
	id := c.Param("id")
	var article models.Post
	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func updateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Post

	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errMsg := validateArticle(input); errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	config.DB.Model(&article).Updates(models.Post{
		Title:    input.Title,
		Content:  input.Content,
		Category: input.Category,
		Status:   input.Status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Post

	if err := config.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	config.DB.Delete(&article)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
