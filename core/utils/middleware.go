package utils

import (
	"context"
	"errors"
	"net/http"
	"reflect"
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
		pagination.Rows = &[]models.Question{}
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
				ctxWithUser := context.WithValue(r.Context(), contextOutKey, &claims)
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

func MarkerCheckMiddleware(db *gorm.DB, contextInKey string, claimsContextKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var count int64
			data := r.Context().Value(contextInKey).(*models.Grade)
			claims := r.Context().Value(claimsContextKey).(*ClaimsData)
			db.Model(&models.Pairing{}).Where("id = ? AND marker_id = ?", data.PairingID, claims.Data.(*models.Student).ID).Count(&count)
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
			claims := r.Context().Value(claimsContextKey).(*ClaimsData)
			db.Model(&models.Pairing{}).Where("id = ? AND student_id = ?", data.PairingID, claims.Data.(*models.Student).ID).Count(&count)
			if count == 0 {
				HandleResponse(w, "Please don't cheat", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
