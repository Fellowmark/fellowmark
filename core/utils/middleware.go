package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

func DecodeBodyMiddleware(refType interface{}, contextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsedData := refType
			if err := DecodeBody(r.Body, &parsedData); err != nil {
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

func DBCreateHandleFunc(db *gorm.DB, tableName string, contextInKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(contextInKey)
		result := db.Table(tableName).Omit("ID").Create(data)
		if result.Error != nil {
			HandleResponse(w, "Already Exists", http.StatusBadRequest)
		} else {
			HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func DBCreateMiddleware(db *gorm.DB, tableName string, contextInKey string, contextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(contextInKey)
			result := db.Table(tableName).Create(data)
			if result.Error != nil {
				HandleResponse(w, "Already Exists", http.StatusBadRequest)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func DBGetFromData(db *gorm.DB, tableName string, contextInKey string, arrayRefType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(contextInKey)
		result := db.Table(tableName).Where(data).Find(arrayRefType)
		if result.Error != nil {
			HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
		} else {
			HandleResponseWithObject(w, arrayRefType, http.StatusOK)
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
				ctxWithUser := context.WithValue(r.Context(), contextOutKey, claims)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}
