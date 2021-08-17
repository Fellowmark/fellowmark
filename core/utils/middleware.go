package utils

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DecodeBodyMiddleware(refType interface{}, contextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsedData := reflect.New(reflect.TypeOf(refType)).Interface()
			if err := DecodeBody(r.Body, parsedData); err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				ctxWithUser := context.WithValue(r.Context(), contextOutKey, parsedData)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}

func SanitizeDataMiddleware(contextInKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := validator.Validate(r.Context().Value(contextInKey))
			if err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func DBCreateHandleFunc(db *gorm.DB, model interface{}, contextInKey string, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(contextInKey)
		result := db.Model(model).Omit("ID").Create(data)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRegistered) {
				if update {
					DBUpdateHandleFunc(db, model, contextInKey).ServeHTTP(w, r)
				}
			} else {
				HandleResponse(w, "Already Exists", http.StatusBadRequest)
			}
		} else {
			HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func DBUpdateHandleFunc(db *gorm.DB, model interface{}, contextInKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(contextInKey)
		result := db.Model(model).Updates(data)
		if result.Error != nil {
			HandleResponse(w, "Already Exists", http.StatusBadRequest)
		} else {
			HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func DBCreateMiddleware(db *gorm.DB, model interface{}, contextInKey string, update bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(DBCreateHandleFunc(db, model, contextInKey, update))
	}
}

func SuccessMiddleware(db *gorm.DB, contextInKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(contextInKey)
		HandleResponseWithObject(w, data, http.StatusOK)
	}
}

func DBGetFromData(db *gorm.DB, model interface{}, contextInKey string, arrayRefType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := GetPagination(r)
		pagination.Rows = reflect.New(reflect.TypeOf(arrayRefType)).Interface()
		data := r.Context().Value(contextInKey)
		scope := Paginate(db, func(tx *gorm.DB) *gorm.DB {
			return tx.Model(model).Where(data)
		}, r, &pagination)
		result := db.Scopes(scope).Preload(clause.Associations).Where(data).Find(pagination.Rows)
		if result.Error != nil {
			HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
		} else {
			HandleResponseWithObject(w, pagination, http.StatusOK)
		}
	}
}

func LoginHandleFunc(db *gorm.DB, role string, contextInKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextInKey)
		token, err := GenerateJWT(role, user)
		if err != nil {
			HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			HandleResponse(w, token, http.StatusOK)
		}
	}
}

func ValidateJWTMiddleware(role string, contextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.Contains(authHeader, "Bearer") {
				HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.Split(authHeader, "Bearer ")[1]
			claims, err := ParseJWT(tokenString)
			if err != nil || claims.Role != role {
				HandleResponse(w, err.Error(), http.StatusUnauthorized)
			} else {
				ctxWithUser := context.WithValue(r.Context(), contextOutKey, &claims.Data)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}

func EnrollmentCheckMiddleware(db *gorm.DB, contextInKey string, muxVarKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(contextInKey).(*models.Student)
			moduleId := mux.Vars(r)[muxVarKey]
			var count int64
			db.Model(&models.Enrollment{}).Where("student_id = ? and module_id = ?", data.ID, moduleId).Count(&count)
			if count == 0 {
				HandleResponse(w, "Not enrolled in module", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func SupervisionCheckMiddleware(db *gorm.DB, contextInKey string, muxVarKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(contextInKey).(*models.Staff)
			moduleId := mux.Vars(r)[muxVarKey]
			var count int64
			db.Model(&models.Enrollment{}).Where("staff_id = ? and module_id = ?", data.ID, moduleId).Count(&count)
			if count == 0 {
				HandleResponse(w, "Not enrolled in module", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func MarkerCheckMiddleware(db *gorm.DB, contextInKey string, studentContextKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var count int64
			data := r.Context().Value(contextInKey).(*models.Grade)
			claims := r.Context().Value(studentContextKey).(*models.Student)
			db.Model(&models.Pairing{}).Where("id = ? AND marker_id = ?", data.PairingID, claims.ID).Count(&count)
			if count == 0 {
				HandleResponse(w, "Please don't cheat", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func MarkeeCheckMiddleware(db *gorm.DB, contextInKey string, claimsContextKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var count int64
			data := r.Context().Value(contextInKey).(*models.Grade)
			student := r.Context().Value(claimsContextKey).(*models.Student)
			db.Model(&models.Pairing{}).Where("id = ? AND student_id = ?", data.PairingID, student.ID).Count(&count)
			if count == 0 {
				HandleResponse(w, "Please don't cheat", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func ValidateAssignmentIdMiddlware(db *gorm.DB, muxVarKey string, moduleIdMuxVarKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var assignment models.Assignment
			urlVars := mux.Vars(r)
			assignmentId := urlVars[muxVarKey]
			result := db.Model(&models.Assignment{}).Where("id = ?", assignmentId).First(&assignment)
			urlVars[moduleIdMuxVarKey] = strconv.Itoa(int(assignment.ModuleID))
			if result.Error != nil {
				HandleResponse(w, "Assignment not found", http.StatusNotFound)
			} else {
				next.ServeHTTP(w, mux.SetURLVars(r, urlVars))
			}
		})
	}
}

func GetAssignedPairingsHandlerFunc(db *gorm.DB, claimsContextKey string, muxVarKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pairings []models.Pairing
		assignmentId := mux.Vars(r)[muxVarKey]
		student := r.Context().Value(claimsContextKey).(*models.Student)
		result := db.Model(&models.Pairing{}).Where("assignment_id = ? AND (student_id = ? OR marker_id = ?)", assignmentId, student.ID, student.ID).Find(&pairings)
		if result.Error != nil {
			HandleResponse(w, "Something went wrong", http.StatusInternalServerError)
		} else if len(pairings) == 0 {
			HandleResponse(w, "No pairing found", http.StatusNotFound)
		} else {
			HandleResponseWithObject(w, pairings, http.StatusOK)
		}
	}
}
