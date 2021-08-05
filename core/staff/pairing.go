package staff

import (
	"net/http"

	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (ur StaffRoute) AssignPairings(w http.ResponseWriter, r *http.Request) {
	assignment := r.Context().Value("assignment").(models.Assignment)
	result := utils.SetNewPairings(ur.DB, assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, "Success", http.StatusCreated)
	}
}

func (ur StaffRoute) InitializePairings(w http.ResponseWriter, r *http.Request) {
	assignment := r.Context().Value("assignment").(models.Assignment)
	result := utils.InitializePairings(ur.DB, assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, "Success", http.StatusOK)
	}
}
