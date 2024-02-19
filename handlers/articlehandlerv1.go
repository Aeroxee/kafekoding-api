package handlers

import (
	"fmt"
	"net/http"

	"github.com/Aeroxee/kafekoding-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
)

type ArticleHandlerV1 struct{}

func NewArticleHandlerV1() ArticleHandlerV1 {
	return ArticleHandlerV1{}
}

func (ArticleHandlerV1) CreateHandler(ctx *gin.Context) {
	payloads := struct {
		Title   string               `json:"title" validate:"required"`
		Content string               `json:"content" validate:"required"`
		Status  models.ArticleStatus `json:"status" validate:"required"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&payloads)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err.Error())
			return
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Error on field: %s, with %s", err.Field(), err.ActualTag()))
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":   "error",
			"message":  "Validation error",
			"messages": errorMessages,
		})
		return
	}

	// get this user info
	thisUser, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication required",
		})
		return
	}

	article := models.Article{
		Title:   payloads.Title,
		UserID:  thisUser.ID,
		Slug:    slug.MakeLang(payloads.Title, "id"),
		Content: payloads.Content,
		Status:  payloads.Status,
	}

	articleModel := models.NewArticleModel(models.DB())
	err = articleModel.CreateNewArticle(&article)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"article": article,
		"message": "Successfullt create new article.",
	})
}

func (ArticleHandlerV1) Get(ctx *gin.Context) {
	page := getQueryInt(ctx.Request, "page", 1)
	size := getQueryInt(ctx.Request, "size", 10)
	status := getQueryString(ctx.Request, "status", "PUBLISHED")

	// calculate offset based on page and size.
	offset := (page - 1) * size

	articleModel := models.NewArticleModel(models.DB())
	articles := articleModel.GetAllArticle(models.ArticleStatus(status), size, offset)

	var count int64
	err := models.DB().Model(&models.Article{}).Count(&count).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"page":     page,
		"size":     size,
		"total":    count,
	})
}

func (ArticleHandlerV1) Update(ctx *gin.Context) {
	thisUser, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication is required.",
		})
		return
	}

	slugArticle := ctx.Param("slug")
	article, err := models.NewArticleModel(models.DB()).GetArticleBySlug(slugArticle)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Article with slug: %s is not found error.", slugArticle),
		})
		return
	}

	if thisUser.ID != article.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "You don't have permission to update this article.",
		})
		return
	}

	payloads := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}{}
	err = ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	if payloads.Title != "" {
		article.Title = payloads.Title
	}
	if payloads.Content != "" {
		article.Content = payloads.Content
	}

	article.Status = models.ArticleStatus(payloads.Status)

	// save
	models.DB().Save(&article)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Update article successfully",
		"article": article,
	})
}

func (ArticleHandlerV1) Detail(ctx *gin.Context) {
	slugArticle := ctx.Param("slug")
	article, err := models.NewArticleModel(models.DB()).GetArticleBySlug(slugArticle)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Article with slug: %s is not found error.", slugArticle),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "",
		"article": article,
	})
}

func (ArticleHandlerV1) Delete(ctx *gin.Context) {
	thisUser, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication is required.",
		})
		return
	}

	slugArticle := ctx.Param("slug")
	article, err := models.NewArticleModel(models.DB()).GetArticleBySlug(slugArticle)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Article with slug: %s is not found error.", slugArticle),
		})
		return
	}

	if thisUser.ID != article.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "You don't have permission to update this article.",
		})
		return
	}

	models.DB().Delete(&article)
	ctx.JSON(http.StatusNoContent, nil)
}
