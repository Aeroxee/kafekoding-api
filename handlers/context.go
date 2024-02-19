package handlers

import (
	"net/http"

	"github.com/Aeroxee/kafekoding-api/auth"
	"github.com/Aeroxee/kafekoding-api/models"
)

// get user info from request context.
func getUserFromContext(r *http.Request) (models.User, error) {
	claims := r.Context().Value(&auth.UserAuth{}).(auth.Claims)
	user, err := models.GetUserByID(claims.Credential.UserID)
	return user, err
}
