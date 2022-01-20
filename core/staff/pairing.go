package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (controller StaffController) AssignPairings(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(utils.DecodeBodyContextKey)
	assignment := &models.Assignment{}
	result := controller.DB.Model(&models.Assignment{}).Where(data).First(assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		result = utils.SetNewPairings(controller.DB, *assignment)
		if result.Error != nil {
			loggers.ErrorLogger.Println(result.Error.Error())
			utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			utils.HandleResponse(w, "Success", http.StatusCreated)
		}
	}
}

func (controller StaffController) InitializePairings(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(utils.DecodeBodyContextKey)
	assignment := &models.Assignment{}
	result := controller.DB.Model(&models.Assignment{}).Where(data).Find(assignment)

	if result.Error != nil {
		loggers.ErrorLogger.Println(result.Error.Error())
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		result = utils.InitializePairings(controller.DB, (*assignment))
		if result.Error != nil {
			loggers.ErrorLogger.Println(result.Error.Error())
			utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			utils.HandleResponse(w, "Success", http.StatusCreated)
		}
	}
}

func (controller StaffController) CreatePairingsPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			data := r.Context().Value(utils.DecodeBodyContextKey)
			assignment := &models.Assignment{}
			controller.DB.Model(&models.Assignment{}).Where(data).Find(assignment)
			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}
