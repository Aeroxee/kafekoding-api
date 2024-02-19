package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aeroxee/kafekoding-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ClassHandlerV1 struct{}

func NewClassHandlerV1() ClassHandlerV1 {
	return ClassHandlerV1{}
}

// CreateHandler is function to handler creating class.
func (c ClassHandlerV1) CreateHandler(ctx *gin.Context) {
	user, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication required",
		})
		return
	}

	if user.Type != models.ADMIN {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access danied",
		})
		return
	}

	payloads := struct {
		Title string `form:"title" validate:"required"`
		// Slug        string                `form:"slug" validate:"required"`
		Description string                `form:"description" validate:"required"`
		Logo        *multipart.FileHeader `form:"logo" validate:"required"`
		IsActive    bool                  `form:"is_active"`
	}{}
	err = ctx.ShouldBindWith(&payloads, binding.FormMultipart)
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
			errorMessages = append(errorMessages, fmt.Sprintf("Error on field: %s with %s.", err.Field(), err.ActualTag()))
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":   "error",
			"message":  "Validation error",
			"messages": errorMessages,
		})
		return
	}

	newSlug := slug.MakeLang(payloads.Title, "id")

	class := models.Class{
		Title:       payloads.Title,
		Slug:        newSlug,
		Description: payloads.Description,
		IsActive:    payloads.IsActive,
	}

	// check extension file
	if !isAllowedExtension(filepath.Ext(payloads.Logo.Filename)) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Please upload image only.",
		})
		return
	}
	filename := uuid.NewString() + filepath.Ext(payloads.Logo.Filename)
	destination := fmt.Sprintf("media/classes/%s/%s", newSlug, filename)

	// upload image to server
	err = ctx.SaveUploadedFile(payloads.Logo, destination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	class.Logo = &destination

	// save to db
	err = models.CreateNewClass(&class)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, class)
}

// Get is handler with request GET.
func (c ClassHandlerV1) Get(ctx *gin.Context) {
	is_active := getQueryBool(ctx.Request, "is_active", true)

	classes := models.GetAllClass(is_active)
	ctx.JSON(http.StatusOK, classes)
}

// Detail is handler to get detail of class.
func (c ClassHandlerV1) Detail(ctx *gin.Context) {
	slugClass := ctx.Param("slug")
	add_mentor := ctx.QueryArray("add_mentor")
	add_member := ctx.QueryArray("add_member")
	delete_mentor := ctx.QueryArray("delete_mentor")
	delete_member := ctx.QueryArray("delete_member")

	class, err := models.GetClassBySlug(slugClass)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Class not found.",
		})
		return
	}

	// add mentor
	if len(add_mentor) > 0 {
		isUser, err := getUserFromContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		// check if this user is mentor for this class.
		for _, m := range class.Mentors {
			if m.Username != isUser.Username {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "This method is not permitted",
				})
				return
			}
		}

		var users []*models.User
		for _, m := range add_mentor {
			user, err := models.GetUserByUsername(m)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "User with username " + m + " is not found.",
				})
				return
			}

			var isJoined bool = false
			for _, mentor := range class.Mentors {
				if mentor.Username == user.Username {
					isJoined = true
				}
			}

			if !isJoined {
				users = append(users, &user)
			}
		}

		models.DB().Model(&class).Association("Mentors").Append(users)
		models.DB().Save(&class)
	}

	// add member
	if len(add_member) > 0 {
		isUser, err := getUserFromContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		// check if this user is mentor for this class.
		for _, m := range class.Mentors {
			if m.Username != isUser.Username {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "This method is not permitted",
				})
				return
			}
		}

		var users []*models.User
		for _, m := range add_member {
			user, err := models.GetUserByUsername(m)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "User with username " + m + " is not found.",
				})
				return
			}

			var isJoined bool = false
			for _, member := range class.Members {
				if member.Username == user.Username {
					isJoined = true
				}
			}

			if !isJoined {
				users = append(users, &user)
			}
		}

		models.DB().Model(&class).Association("Members").Append(users)
	}

	// remove mentor
	if len(delete_mentor) > 0 {
		isUser, err := getUserFromContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		// check if this user is mentor for this class.
		for _, m := range class.Mentors {
			if m.Username != isUser.Username {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "This method is not permitted",
				})
				return
			}
		}

		var users []*models.User
		for _, m := range delete_mentor {
			user, err := models.GetUserByUsername(m)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "User with username " + m + " is not found.",
				})
				return
			}

			var isJoined bool = false
			for _, mentor := range class.Mentors {
				if mentor.Username == user.Username {
					isJoined = true
				}
			}

			if isJoined {
				users = append(users, &user)
			}
		}

		models.DB().Model(&class).Association("Mentors").Delete(users)
	}

	// remove member
	if len(delete_member) > 0 {
		isUser, err := getUserFromContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		// check if this user is mentor for this class.
		for _, m := range class.Mentors {
			if m.Username != isUser.Username {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "This method is not permitted",
				})
				return
			}
		}

		var users []*models.User
		for _, m := range delete_member {
			user, err := models.GetUserByUsername(m)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "User with username " + m + " is not found.",
				})
				return
			}

			var isJoined bool = false
			for _, member := range class.Members {
				if member.Username == user.Username {
					isJoined = true
				}
			}

			if isJoined {
				users = append(users, &user)
			}
		}

		models.DB().Model(&class).Association("Members").Delete(users)
	}

	ctx.JSON(http.StatusOK, class)
}

// Update handler
func (c ClassHandlerV1) Update(ctx *gin.Context) {
	thisUser, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication is required.",
		})
		return
	}

	slugClass := ctx.Param("slug")
	class, err := models.GetClassBySlug(slugClass)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Class not found.",
		})
		return
	}

	// check this user what is mentor of this class.
	var isMentor bool = false
	for _, mentor := range class.Mentors {
		if mentor.Username == thisUser.Username {
			isMentor = true
		}
	}

	if !isMentor || thisUser.Type == models.MEMBER {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access danied.",
		})
		return
	}

	// payloads
	payloads := struct {
		Title       string                `form:"title"`
		Description string                `form:"description"`
		Logo        *multipart.FileHeader `form:"logo"`
		IsActive    bool                  `form:"is_active"`
	}{}
	err = ctx.ShouldBindWith(&payloads, binding.FormMultipart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	if payloads.Title != "" {
		class.Title = payloads.Title
		class.Slug = slug.MakeLang(payloads.Title, "id")
		// move file
		s := strings.Split(*class.Logo, "/")
		oldFilename := s[3]

		oldDestination := class.Logo
		newDestination := fmt.Sprintf("media/classes/%s/%s", class.Slug, oldFilename)

		// read old file
		sourceFile, err := os.Open(*oldDestination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		defer sourceFile.Close()

		os.MkdirAll(fmt.Sprintf("media/classes/%s", class.Slug), 0700)
		// make new file
		newFile, err := os.Create(newDestination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		defer newFile.Close()

		// copy file to source destination
		_, err = io.Copy(newFile, sourceFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		class.Logo = &newDestination
		models.DB().Save(&class)
	}

	if payloads.Description != "" {
		class.Description = payloads.Description
		models.DB().Save(&class)
	}

	if payloads.Logo != nil {
		// remove old file
		os.RemoveAll(*class.Logo)

		if !isAllowedExtension(filepath.Ext(payloads.Logo.Filename)) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Please upload image only.",
			})
			return
		}
		newFilename := uuid.NewString() + filepath.Ext(payloads.Logo.Filename)
		newDestination := fmt.Sprintf("media/classes/%s/%s", class.Slug, newFilename)
		err := ctx.SaveUploadedFile(payloads.Logo, newDestination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		class.Logo = &newDestination
		models.DB().Save(&class)
	}

	class.IsActive = payloads.IsActive
	models.DB().Save(&class)

	ctx.JSON(http.StatusOK, class)
}

// Delete handler
func (c ClassHandlerV1) Delete(ctx *gin.Context) {
	thisUser, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Authentication is required.",
		})
		return
	}

	slugClass := ctx.Param("slug")
	class, err := models.GetClassBySlug(slugClass)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Class not found.",
		})
		return
	}

	// only admin can delete
	if thisUser.Type != models.ADMIN {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Permission danied",
		})
		return
	}

	models.DB().Delete(&class)
	ctx.JSON(http.StatusNoContent, nil)
}
