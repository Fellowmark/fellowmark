package staff

import (
	"context"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (ur StaffRoute) PasswordHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(utils.DecodeBodyContextKey).(*models.Staff)
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusInternalServerError)
		} else {
			user.Password = hash
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}
