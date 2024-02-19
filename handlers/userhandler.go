package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Aeroxee/kafekoding-api/auth"
	"github.com/Aeroxee/kafekoding-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type UserHandlerV1 struct{}

func NewUserHandlerV1() UserHandlerV1 {
	return UserHandlerV1{}
}

// validate
var validate *validator.Validate

// ActivationData struct for save activation data
type ActivationData struct {
	Email          string
	ExpirationTime time.Time
}

// activation data
var activationData = make(map[string]ActivationData)

// generate activation code
func (u *UserHandlerV1) generateActivationCode() string {
	return uuid.New().String()
}

// send activation code to email target.
func (u *UserHandlerV1) sendActivationEmail(email, activationCode string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	to := []string{email}
	subject := "Activate Your Account"
	body := fmt.Sprintf("Click the following link to activate your account: http://localhost:8000/v1/activate/%s", activationCode)
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", strings.Join(to, ","), subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUsername, to, []byte(message))
	if err == nil {
		return true
	} else {
		return false
	}
}

// RegisterHandler is handler to regitration user.
func (u *UserHandlerV1) RegisterHandler(ctx *gin.Context) {
	payloads := struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	// validate field
	validate = validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&payloads)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err.Error())
			return
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Error on %s, with %s", err.Field(), err.ActualTag()))
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":   "error",
			"message":  "Validation error",
			"messages": errorMessages,
		})
		return
	}

	user := models.User{
		FirstName: payloads.FirstName,
		LastName:  payloads.LastName,
		Username:  payloads.Username,
		Email:     payloads.Email,
		Password:  payloads.Password,
	}

	// save user to db
	err = models.CreateNewUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// generate activation code
	activationCode := u.generateActivationCode()

	// Save activation code and expiration time in the map
	activationData[activationCode] = ActivationData{
		Email:          payloads.Email,
		ExpirationTime: time.Now().Add(24 * time.Hour),
	}

	// Send email activation
	if !u.sendActivationEmail(activationData[activationCode].Email, activationCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to send activation code to your email.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registration successful. Check your email for activation instructions.",
	})
}

// ActivationHandler is handler for activation account/user
func (u *UserHandlerV1) ActivationHandler(ctx *gin.Context) {
	activationCode := ctx.Param("activationCode")

	data, ok := activationData[activationCode]
	if !ok || time.Now().After(data.ExpirationTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid or expired activation code.",
		})
		return
	}

	user, err := models.GetUserByEmail(data.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("User with email: %s is not found.", data.Email),
		})
		return
	}

	user.IsActive = true
	models.DB().Save(&user)
	delete(activationData, activationCode)

	ctx.Writer.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(ctx.Writer, "<h1>Your account is activated successfully.</h1>")
}

// GetTokenHandler is handler to generate new token.
func (u *UserHandlerV1) GetTokenHandler(ctx *gin.Context) {
	payloads := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	var user models.User
	if strings.Contains(payloads.Username, "@") {
		user, err = models.GetUserByEmail(payloads.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect.",
			})
			return
		}
	} else {
		user, err = models.GetUserByUsername(payloads.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username or password is incorrect.",
			})
			return
		}
	}

	if !auth.DecryptionPassword(user.Password, payloads.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username or password is incorrect.",
		})
		return
	}

	// check if user is not active
	// generate token not permitted.
	if !user.IsActive {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Your account is not yet active, please check your email for the activation code link.",
		})
		return
	}

	credential := auth.Credential{
		UserID:   user.ID,
		Username: user.Username,
	}

	token, err := auth.GetToken(credential)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Generate new token is successfully.",
		"token":   token,
	})
}

// CheckAuthHandler is handler to check authentication user.
func (u *UserHandlerV1) CheckAuthHandler(ctx *gin.Context) {
	user, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Your token is not valid.",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ChangePasswordHandler handler to change password user account.
func (u *UserHandlerV1) ChangePasswordHandler(ctx *gin.Context) {
	payloads := struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}{}

	err := ctx.ShouldBindJSON(&payloads)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload error",
		})
		return
	}

	user, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Your token is not valid.",
		})
		return
	}

	if !auth.DecryptionPassword(user.Password, payloads.OldPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Your password is incorrect.",
		})
		return
	}

	if payloads.NewPassword != payloads.NewPasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Password confirmation not same.",
		})
		return
	}

	// update and save password
	newPassword := auth.EncryptionPassword(payloads.NewPasswordConfirm)
	user.Password = newPassword
	models.DB().Save(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Update password is successfully.",
	})
}

// UpdateInfoUserHandler handler to update infor user.
func (u *UserHandlerV1) UpdateInfoUserHandler(ctx *gin.Context) {
	payloads := struct {
		FirstName string                `form:"first_name"`
		LastName  string                `form:"last_name"`
		Username  string                `form:"username"`
		Avatar    *multipart.FileHeader `form:"avatar"`
	}{}
	err := ctx.ShouldBindWith(&payloads, binding.FormMultipart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload type is not multipart/form-data",
		})
		return
	}

	// get user context
	user, err := getUserFromContext(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Your token is not valid.",
		})
		return
	}

	if payloads.FirstName != "" {
		user.FirstName = payloads.FirstName
	}
	if payloads.LastName != "" {
		user.LastName = payloads.LastName
	}
	if payloads.Username != "" {
		user.Username = payloads.Username
	}
	if payloads.Avatar != nil {
		if user.Avatar != nil {
			os.RemoveAll(*user.Avatar)
		}

		if !isAllowedExtension(filepath.Ext(payloads.Avatar.Filename)) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Please upload an image only.",
			})
			return
		}

		ext := filepath.Ext(payloads.Avatar.Filename)
		filename := uuid.NewString() + ext
		destination := fmt.Sprintf("media/user/profile/%s/%s", user.Username, filename)

		// save
		err = ctx.SaveUploadedFile(payloads.Avatar, destination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		user.Avatar = &destination
	}

	// save
	models.DB().Save(&user)
	ctx.JSON(http.StatusOK, user)
}
