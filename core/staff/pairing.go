package staff

import (
	"net/http"

	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (ur StaffController) AssignPairings(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("assignment")
	assignment := &models.Assignment{}
	result := ur.DB.Model(&models.Assignment{}).Where(data).First(assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		result = utils.SetNewPairings(ur.DB, *assignment)
		if result.Error != nil {
			loggers.ErrorLogger.Println(result.Error.Error())
			utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			utils.HandleResponse(w, "Success", http.StatusCreated)
		}
	}
}

func (ur StaffController) InitializePairings(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("assignment")
	assignment := &models.Assignment{}
	result := ur.DB.Model(&models.Assignment{}).Where(data).Find(assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		result = utils.InitializePairings(ur.DB, (*assignment))
		if result.Error != nil {
			loggers.ErrorLogger.Println(result.Error.Error())
			utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			utils.HandleResponse(w, "Success", http.StatusCreated)
		}
	}
}
