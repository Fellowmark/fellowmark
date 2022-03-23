package module

import (
	"net/http"
	"strings"
	"context"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
	"errors"
)

const EmailsNotFoundKey = "emailNotFoundIndexes"
const EnrollmentExistIndexesKey = "enrollmentExistIndexes"
const EnrollErrorsKey = "enrollErrors"
const SupervisionExistIndexesKey = "supervisionExistIndexes"
const SuperviseErrorsKey = "superviseErrors"

type EnrollmentResult struct {
	SuccessCount int `json:"success"`
	EnrollErrors []string `json:"enrollErrors"`
}

type SupervisionResult struct {
	SuccessCount int `json:"success"`
	SuperviseErrors []string `json:"superviseErrors"`
}

func (controller ModuleController) ModuleCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		module := r.Context().Value(utils.DecodeBodyContextKey).(*models.Module)
		user := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
		
		txError := controller.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(module).Create(module).Error; err != nil {
				return err
			}
			supervision := &models.Supervision{ModuleID: module.ID, StaffID: user.ID}
			if err := tx.Model(supervision).Create(supervision).Error; err != nil {
				return err
			}
			return nil
		})
		if txError != nil {
			var errMessage string
			if strings.Contains(txError.Error(), "duplicate key value violates unique constraint") {
				errMessage = "Creation failed: Module already exists."
			} else {
				errMessage = "Creation failed: " + txError.Error()
			}
			utils.HandleResponse(w, errMessage, http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, module, http.StatusOK)
	}
}

func (controller ModuleController) EnrollmentCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enrollments := r.Context().Value(utils.DecodeBodyContextKey).(*[]models.Enrollment)
		enrollErrors := r.Context().Value(EnrollErrorsKey).(*[]string)
		result := controller.DB.Model(enrollments).Create(enrollments)
		if result.Error != nil {
			if result.Error.Error() == "empty slice found" {
				response := EnrollmentResult{0, *enrollErrors}
				utils.HandleResponseWithObject(w, &response, http.StatusOK)
			} else {
				utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			}
		} else {
			response := EnrollmentResult{len(*enrollments), *enrollErrors}
			utils.HandleResponseWithObject(w, &response, http.StatusOK)
		}
	}
}

func (controller ModuleController) EnrollmentDeleteHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enrollment := r.Context().Value(utils.DecodeBodyContextKey).(*models.Enrollment)
		txError := controller.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(enrollment).Where(enrollment).Take(enrollment).Error; err != nil {
				return err
			}
			if err := tx.Model(enrollment).Delete(enrollment).Error; err != nil {
				return err
			}
			return nil
		})
		if txError != nil {
			var errMessage string
			if errors.Is(txError, gorm.ErrRecordNotFound) {
				errMessage = "Deletion failed: Student not found."
			} else {
				errMessage = "Deletion failed: " + txError.Error()
			}
			utils.HandleResponse(w, errMessage, http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, enrollment, http.StatusOK)
	}
}

func (controller ModuleController) SupervisionCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supervisions := r.Context().Value(utils.DecodeBodyContextKey).(*[]models.Supervision)
		superviseErrors := r.Context().Value(SuperviseErrorsKey).(*[]string)
		result := controller.DB.Model(supervisions).Create(supervisions)
		if result.Error != nil {
			if result.Error.Error() == "empty slice found" {
				response := SupervisionResult{0, *superviseErrors}
				utils.HandleResponseWithObject(w, &response, http.StatusOK)
			} else {
				utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			}
		} else {
			response := SupervisionResult{len(*supervisions), *superviseErrors}
			utils.HandleResponseWithObject(w, &response, http.StatusOK)
		}
	}
}

func (controller ModuleController) SupervisionDeleteHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supervision := r.Context().Value(utils.DecodeBodyContextKey).(*models.Supervision)
		txError := controller.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(supervision).Where(supervision).Take(supervision).Error; err != nil {
				return err
			}
			if err := tx.Model(supervision).Delete(supervision).Error; err != nil {
				return err
			}
			return nil
		})
		if txError != nil {
			var errMessage string
			if errors.Is(txError, gorm.ErrRecordNotFound) {
				errMessage = "Deletion failed: Supervisor not found."
			} else {
				errMessage = "Deletion failed: " + txError.Error()
			}
			utils.HandleResponse(w, errMessage, http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, supervision, http.StatusOK)
	}
}

func (controller ModuleController) GetStudentEnrollmentsHandleFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
		controller.DB.Joins("inner join enrollments e on modules.id = e.module_id").Where("e.student_id = ?", user.ID).Find(&modules)
		utils.HandleResponseWithObject(w, modules, http.StatusOK)
	})
}

func (controller ModuleController) GetStaffSupervisionsHandleFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
		controller.DB.Joins("inner join supervisions e on modules.id = e.module_id").Where("e.staff_id = ?", user.ID).Find(&modules)
		utils.HandleResponseWithObject(w, modules, http.StatusOK)
	})
}

func (controller ModuleController) CheckStaffSupervision() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			data := r.Context().Value(utils.DecodeBodyContextKey)
			var moduleID uint
			switch data.(type) {
				case *BatchEnrollment:
					moduleID = data.(*BatchEnrollment).ModuleID
				case *BatchSupervision:
					moduleID = data.(*BatchSupervision).ModuleID
				case *models.Supervision:
					moduleID = data.(*models.Supervision).ModuleID
				case *models.Enrollment:
					moduleID = data.(*models.Enrollment).ModuleID
			}
			if pass := utils.IsSupervisor(*claims, moduleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller ModuleController) EnrollmentDataPrepare() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const duplicateErrorMessage = "Enrollment exists"
			const studentNotFoundErrorMessage = "Student not found"
			data := r.Context().Value(utils.DecodeBodyContextKey).(*BatchEnrollment)
			if (data.StudentID > 0) {
				enrollments := make([]models.Enrollment, 0, 1)
				enrollErrors := make([]string, 1, 1)
				existEnrollment := models.Enrollment{}
				if controller.DB.Take(&models.Student{}, data.StudentID).Error != nil {
					enrollErrors[0] = studentNotFoundErrorMessage
				} else if controller.DB.Model(&existEnrollment).Where("student_id = ? and module_id = ?", data.StudentID, data.ModuleID).Take(&existEnrollment).Error == nil {
					enrollErrors[0] = duplicateErrorMessage
				} else {
					enrollments = append(enrollments, models.Enrollment{ModuleID: data.ModuleID, StudentID: data.StudentID})
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &enrollments)
				ctx = context.WithValue(ctx, EnrollErrorsKey, &enrollErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else if (len(data.StudentIDs) > 0) {
				enrollments := make([]models.Enrollment, 0, len(data.StudentIDs))
				enrollErrors := make([]string, len(data.StudentIDs), len(data.StudentIDs))
				for i, studentID := range data.StudentIDs {
					existEnrollment := models.Enrollment{}
					if controller.DB.Take(&models.Student{}, studentID).Error != nil {
						enrollErrors[i] = studentNotFoundErrorMessage
					} else if controller.DB.Model(&existEnrollment).Where("student_id = ? and module_id = ?", studentID, data.ModuleID).Take(&existEnrollment).Error == nil {
						enrollErrors[i] = duplicateErrorMessage
					} else {
						enrollments = append(enrollments, models.Enrollment{ModuleID: data.ModuleID, StudentID: studentID})
					}
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &enrollments)
				ctx = context.WithValue(ctx, EnrollErrorsKey, &enrollErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else if (len(data.StudentEmails) > 0) {
				enrollments := make([]models.Enrollment, 0, len(data.StudentEmails))
				enrollErrors := make([]string, len(data.StudentEmails), len(data.StudentEmails))
				for i, email := range data.StudentEmails {
					student := models.Student{}
					studentResult := controller.DB.Model(&student).Where("email = ?", email).Take(&student)
					existEnrollment := models.Enrollment{}
					if studentResult.Error != nil {
						enrollErrors[i] = studentNotFoundErrorMessage
					} else if controller.DB.Model(&existEnrollment).Where("student_id = ? and module_id = ?", student.ID, data.ModuleID).Take(&existEnrollment).Error == nil {
						enrollErrors[i] = duplicateErrorMessage
					} else {
						enrollments = append(enrollments, models.Enrollment{ModuleID: data.ModuleID, StudentID: student.ID})
					}
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &enrollments)
				ctx = context.WithValue(ctx, EnrollErrorsKey, &enrollErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.HandleResponse(w, "Empty Enrollments Data", http.StatusBadRequest)
			}
		})
	}
}

func (controller ModuleController) SupervisionDataPrepare() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const duplicateErrorMessage = "Supervision exists"
			const staffNotFoundErrorMessage = "Staff not found"
			data := r.Context().Value(utils.DecodeBodyContextKey).(*BatchSupervision)
			if (data.StaffID > 0) {
				supervisions := make([]models.Supervision, 0, 1)
				superviseErrors := make([]string, 1, 1)
				existSupervision := models.Supervision{}
				if controller.DB.Take(&models.Staff{}, data.StaffID).Error != nil {
					superviseErrors[0] = staffNotFoundErrorMessage
				} else if controller.DB.Model(&existSupervision).Where("staff_id = ? and module_id = ?", data.StaffID, data.ModuleID).Take(&existSupervision).Error == nil {
					superviseErrors[0] = duplicateErrorMessage
				} else {
					supervisions = append(supervisions, models.Supervision{ModuleID: data.ModuleID, StaffID: data.StaffID})
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &supervisions)
				ctx = context.WithValue(ctx, SuperviseErrorsKey, &superviseErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else if (len(data.StaffIDs) > 0) {
				supervisions := make([]models.Supervision, 0, len(data.StaffIDs))
				superviseErrors := make([]string, len(data.StaffIDs), len(data.StaffIDs))
				for i, staffID := range data.StaffIDs {
					existSupervision := models.Supervision{}
					if controller.DB.Take(&models.Staff{}, staffID).Error != nil {
						superviseErrors[i] = staffNotFoundErrorMessage
					} else if controller.DB.Model(&existSupervision).Where("staff_id = ? and module_id = ?", staffID, data.ModuleID).Take(&existSupervision).Error == nil {
						superviseErrors[i] = duplicateErrorMessage
					} else {
						supervisions = append(supervisions, models.Supervision{ModuleID: data.ModuleID, StaffID: staffID})
					}
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &supervisions)
				ctx = context.WithValue(ctx, SuperviseErrorsKey, &superviseErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else if (len(data.StaffEmails) > 0) {
				supervisions := make([]models.Supervision, 0, len(data.StaffEmails))
				superviseErrors := make([]string, len(data.StaffEmails), len(data.StaffEmails))
				for i, email := range data.StaffEmails {
					staff := models.Staff{}
					staffResult := controller.DB.Model(&staff).Where("email = ?", email).Take(&staff)
					existSupervision := models.Supervision{}
					if staffResult.Error != nil {
						superviseErrors[i] = staffNotFoundErrorMessage
					} else if controller.DB.Model(&existSupervision).Where("staff_id = ? and module_id = ?", staff.ID, data.ModuleID).Take(&existSupervision).Error == nil {
						superviseErrors[i] = duplicateErrorMessage
					} else {
						supervisions = append(supervisions, models.Supervision{ModuleID: data.ModuleID, StaffID: staff.ID})
					}
				}
				ctx := context.WithValue(r.Context(), utils.DecodeBodyContextKey, &supervisions)
				ctx = context.WithValue(ctx, SuperviseErrorsKey, &superviseErrors)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.HandleResponse(w, "Empty Supervisions Data", http.StatusBadRequest)
			}
		})
	}
}