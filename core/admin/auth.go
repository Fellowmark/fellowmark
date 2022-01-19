package admin

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

func (ur AdminRoute) AdminLoginHandleFunc(w http.ResponseWriter, r *http.Request) {
	var input models.Admin
	if err := utils.DecodeParams(r, &input); err != nil {
		utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggers.InfoLogger.Println(input)

	var user models.Admin
	result := ur.DB.Take(&user, "email = ?", input.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		utils.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		return
	}

	isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if !isEqual {
		utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(utils.ADMIN, user)
	if err != nil {
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, token, http.StatusOK)
	}
}
