package assignment

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (controller AssignmentController) CreateAssignmentPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			assignment := r.Context().Value(utils.DecodeBodyContextKey).(*models.Assignment)
			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) CreateQuestionPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			question := r.Context().Value(utils.DecodeBodyContextKey).(*models.Question)
			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", question.AssignmentID).Find(&assignment)
			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) CreateRubricPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			rubric := r.Context().Value(utils.DecodeBodyContextKey).(*models.Rubric)

			var question models.Question
			controller.DB.Model(&models.Assignment{}).Where("id = ?", rubric.QuestionID).Find(&question)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", question.AssignmentID).Find(&assignment)

			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) resolveAssignmentIdFromMuxVar(r *http.Request) string {
	return mux.Vars(r)["assignmentId"]
}

func (controller AssignmentController) resolveAssignmentIdFromContext(r *http.Request) string {
	assignment := r.Context().Value("data").(*models.Assignment)
	return strconv.FormatUint(uint64(assignment.ID), 10)
}

func (controller AssignmentController) AssignPairings(w http.ResponseWriter, r *http.Request) {
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

func (controller AssignmentController) InitializePairings(w http.ResponseWriter, r *http.Request) {
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

func (controller AssignmentController) CreatePairingsPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			input := r.Context().Value(utils.DecodeBodyContextKey).(*models.Assignment)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", input.ID).Find(&assignment)

			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) GetPairingsPermissionsCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			input := r.Context().Value(utils.DecodeParamsContextKey).(*models.Pairing)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", input.AssignmentID).Find(&assignment)

			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) GetAllPairings(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPagination(r)
	pagination.Rows = &[]models.Pairing{}
	data := r.Context().Value(utils.DecodeParamsContextKey).(*models.Pairing)
	scope := utils.Paginate(controller.DB, func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.Pairing{}).Where(models.Pairing{AssignmentID: data.AssignmentID, Active: true})
	}, r, &pagination)
	result := controller.DB.Scopes(scope).Preload(clause.Associations).Find(pagination.Rows)
	if result.Error != nil {
		utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
	} else {
		utils.HandleResponseWithObject(w, pagination, http.StatusOK)
	}
}

func (controller AssignmentController) GetPairingsForRevieweeHandleFunc(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
	pagination := utils.GetPagination(r)
	pagination.Rows = &[]models.Pairing{}
	data := r.Context().Value(utils.DecodeParamsContextKey).(*models.Pairing)

	scope := utils.Paginate(controller.DB, func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.Pairing{}).Where(models.Pairing{AssignmentID: data.AssignmentID, StudentID: claims.ID, Active: true})
	}, r, &pagination)

	result := controller.DB.Scopes(scope).Preload(clause.Associations).Find(pagination.Rows)

	if result.Error != nil {
		utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
	} else {
		utils.HandleResponseWithObject(w, pagination, http.StatusOK)
	}
}

func (controller AssignmentController) GetPairingsForMarkerHandleFunc(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
	pagination := utils.GetPagination(r)
	pagination.Rows = &[]models.Pairing{}
	data := r.Context().Value(utils.DecodeParamsContextKey).(*models.Pairing)

	scope := utils.Paginate(controller.DB, func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.Pairing{}).Where(models.Pairing{AssignmentID: data.AssignmentID, MarkerID: claims.ID, Active: true})
	}, r, &pagination)

	result := controller.DB.Scopes(scope).Preload(clause.Associations).Find(pagination.Rows)

	if result.Error != nil {
		utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
	} else {
		utils.HandleResponseWithObject(w, pagination, http.StatusOK)
	}
}
