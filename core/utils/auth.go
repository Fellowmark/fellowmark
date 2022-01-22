package utils

import (
	"context"
	"errors"
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
	Data models.User `json:"data"`
	jwt.StandardClaims
}

const JWTClaimContextKey = "claims"

func GenerateJWT(user models.User) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))

	claims := ClaimsData{
		user,
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
			data := r.Context().Value(JWTClaimContextKey).(*models.User)
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

func ValidateJWT(r *http.Request) (*ClaimsData, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		return nil, errors.New("Unauthorized")
	}
	tokenString := strings.Split(authHeader, "Bearer ")[1]
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, errors.New("Unauthenticated")
	}
	return claims, nil
}

func AuthenticationMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if claims, err := ValidateJWT(r); err != nil {
				HandleResponse(w, err.Error(), http.StatusUnauthorized)
			} else {
				ctxWithUser := context.WithValue(r.Context(), JWTClaimContextKey, &claims.Data)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}

func IsAdmin(user models.User, db *gorm.DB) bool {
	result := db.Take(&models.Admin{}, "id = ?", user.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func IsSupervisor(user models.User, moduleId uint, db *gorm.DB) bool {
	if bypass := IsAdmin(user, db); bypass {
		return true
	}
	resultStaff := db.Take(&models.Staff{}, "id = ?", user.ID)
	if errors.Is(resultStaff.Error, gorm.ErrRecordNotFound) {
		return false
	}
	resultSupervision := db.Take(&models.Supervision{}, "staff_id = ? AND module_id = ?", user.ID, moduleId)
	if errors.Is(resultSupervision.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func IsEnrolled(user models.User, moduleId uint, db *gorm.DB) bool {
	if bypass := IsAdmin(user, db); bypass {
		return true
	}
	resultStaff := db.Take(&models.Student{}, "id = ?", user.ID)
	if errors.Is(resultStaff.Error, gorm.ErrRecordNotFound) {
		return false
	}
	resultSupervision := db.Take(&models.Enrollment{}, "student_id = ? AND module_id = ? ", user.ID, moduleId)
	if errors.Is(resultSupervision.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// Both Supervisor and Marker are authorized for marking
func IsMarker(user models.User, assignmentId uint, studentId uint, db *gorm.DB) bool {
	if bypass := IsAdmin(user, db); bypass {
		return true
	}

	var assignment models.Assignment
	db.Model(&models.Assignment{}).Where("id = ?", assignmentId).Find(&assignment)
	if bypass := IsSupervisor(user, assignment.ModuleID, db); bypass {
		return true
	}

	result := db.Model(&models.Pairing{}).Take(&models.Pairing{},
		models.Pairing{MarkerID: user.ID, StudentID: studentId, AssignmentID: assignmentId})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func IsReviewee(user models.User, assignmentId uint, markerId uint, db *gorm.DB) bool {
	if bypass := IsAdmin(user, db); bypass {
		return true
	}
	result := db.Model(&models.Pairing{}).Take(&models.Pairing{},
		models.Pairing{StudentID: user.ID, MarkerID: markerId, AssignmentID: assignmentId})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func IsPair(user models.User, assignmentId uint, otherStudentId uint, db *gorm.DB) bool {
	return IsMarker(user, assignmentId, otherStudentId, db) || IsReviewee(user, assignmentId, otherStudentId, db)
}

func IsMemberOf(claims models.User, moduleId uint, db *gorm.DB) bool {
	return IsEnrolled(claims, moduleId, db) || IsSupervisor(claims, moduleId, db)
}

func IsAdminMiddleware(db *gorm.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(JWTClaimContextKey).(*models.User)
			if IsAdmin(*data, db) {
				next.ServeHTTP(w, r)
			} else {
				HandleResponse(w, "Insufficient Permissions", http.StatusUnauthorized)
			}
		})
	}
}

func LoginHandleFunc(db *gorm.DB, scope func(db *gorm.DB) *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.User
		if err := DecodeParams(r, &input); err != nil {
			HandleResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.Scopes(scope).Take(&user, "email = ?", input.Email)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
			return
		}

		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(user)
		if err != nil {
			HandleResponse(w, "Internal Error", http.StatusInternalServerError)
		} else {
			HandleResponse(w, token, http.StatusOK)
		}
	}
}
