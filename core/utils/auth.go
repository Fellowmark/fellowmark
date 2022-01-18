package utils

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"gorm.io/gorm"
)

type ClaimsData struct {
	Role string      `json:"role"`
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

const JWTClaimContextKey = "claims"

func GenerateJWT(role string, object interface{}) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))

	claims := ClaimsData{
		role,
		object,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "npr-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		loggers.ErrorLogger.Println("Something Went Wrong: %s" + err.Error())
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*ClaimsData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*ClaimsData); ok && token.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}

func ParseJWTWithClaims(tokenString string, claims *ClaimsData) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if !token.Valid {
		return err
	}
	return nil
}

// returns argon2 hash of strings such as email recovery tokens
func HashString(token string) string {
	hash, err := argon2id.CreateHash(token, argon2id.DefaultParams)
	if err != nil {
		loggers.ErrorLogger.Println("Something Went Wrong: %s" + err.Error())
	}

	return hash
}

func SupervisionCheckMiddleware(db *gorm.DB, moduleIdResolver func(r *http.Request) string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(JWTClaimContextKey).(*models.Staff)
			moduleId := moduleIdResolver(r)
			var count int64
			db.Model(&models.Supervision{}).Where("staff_id = ? and module_id = ?", data.ID, moduleId).Count(&count)
			if count == 0 {
				HandleResponse(w, "Not enrolled in module", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func MarkerCheckMiddleware(db *gorm.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var count int64
			data := r.Context().Value(DecodeBodyContextKey).(*models.Grade)
			claims := r.Context().Value(JWTClaimContextKey).(*models.Student)
			db.Model(&models.Pairing{}).Where("id = ? AND marker_id = ?", data.PairingID, claims.ID).Count(&count)
			if count == 0 {
				HandleResponse(w, "Please don't cheat", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func MarkeeCheckMiddleware(db *gorm.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var count int64
			data := r.Context().Value(DecodeBodyContextKey).(*models.Grade)
			student := r.Context().Value(JWTClaimContextKey).(*models.Student)
			db.Model(&models.Pairing{}).Where("id = ? AND student_id = ?", data.PairingID, student.ID).Count(&count)
			if count == 0 {
				HandleResponse(w, "Please don't cheat", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func EnrollmentCheckMiddleware(db *gorm.DB, moduleIdResolver func(r *http.Request) string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(JWTClaimContextKey).(*models.Student)
			moduleId := moduleIdResolver(r)
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

func LoginHandleFunc(db *gorm.DB, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(JWTClaimContextKey)
		token, err := GenerateJWT(role, user)
		if err != nil {
			HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			HandleResponse(w, token, http.StatusOK)
		}
	}
}

func ValidateJWTMiddleware(role string, refType interface{}) mux.MiddlewareFunc {
	return ValidateJWTMiddlewareMultipleRoles([]string{role})
}

func ValidateJWTMiddlewareMultipleRoles(roles []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.Contains(authHeader, "Bearer") {
				HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := strings.Split(authHeader, "Bearer ")[1]
			claims, err := ParseJWT(tokenString)
			if err != nil || Contains(roles, claims.Role) {
				HandleResponse(w, err.Error(), http.StatusUnauthorized)
			} else {
				ctxWithUser := context.WithValue(r.Context(), JWTClaimContextKey, &claims.Data)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}
