package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/nus-utils/nus-peer-review/models"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var SchemaDecoder = schema.NewDecoder()

const DecodeBodyContextKey = "data"
const DecodeParamsContextKey = "params"

func SetupCors() mux.MiddlewareFunc {
	return handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodOptions}))
}

func SetHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request (not recommended in production)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		//only allow GET, POST, and OPTIONS
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//Since I was building a REST API that returned JSON, I set the content type to JSON here.
		w.Header().Set("Content-Type", "application/json")
		//Allow requests to have the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		//if it's just an OPTIONS request, nothing other than the headers in the response is needed.
		//This is essential because you don't need to handle the OPTIONS requests in your handlers now
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func DecodeBody(body io.ReadCloser, out interface{}) error {
	var unmarshalErr *json.UnmarshalFieldError
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&out)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Bad Request. Wrong Type provided " + unmarshalErr.Field.Name)
		} else {
			return errors.New("Bad Request. " + err.Error())
		}
	}
	return nil
}

func HandleResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func HandleResponseWithObject(w http.ResponseWriter, object interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(object)
}

func DecodeParamsMiddleware(refType interface{}) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsedData := reflect.New(reflect.TypeOf(refType).Elem()).Interface()
			if err := SchemaDecoder.Decode(parsedData, r.URL.Query()); err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				ctxWithUser := context.WithValue(r.Context(), DecodeParamsContextKey, parsedData)
				next.ServeHTTP(w, r.WithContext(ctxWithUser))
			}
		})
	}
}

func DecodeBodyMiddleware(refType interface{}) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsedData := reflect.New(reflect.TypeOf(refType).Elem()).Interface()
			if err := DecodeBody(r.Body, parsedData); err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				ctxWithData := context.WithValue(r.Context(), DecodeBodyContextKey, parsedData)
				next.ServeHTTP(w, r.WithContext(ctxWithData))
			}
		})
	}
}

func SanitizeDataMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := validator.Validate(r.Context().Value(DecodeBodyContextKey))
			if err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func SanitizeParamsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := validator.Validate(r.Context().Value(DecodeBodyContextKey))
			if err != nil {
				HandleResponse(w, err.Error(), http.StatusBadRequest)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func DBCreateHandleFunc(db *gorm.DB, model interface{}, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(DecodeBodyContextKey)
		result := db.Model(model).Omit("ID").Create(data)
		if result.Error != nil {
			if update {
				DBUpdateHandleFunc(db, model).ServeHTTP(w, r)
			} else {
				HandleResponse(w, "Already Exists", http.StatusBadRequest)
			}
		} else {
			HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func DBUpdateHandleFunc(db *gorm.DB, model interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(DecodeBodyContextKey)
		result := db.Model(model).Updates(data)
		if result.Error != nil {
			HandleResponse(w, "Failed", http.StatusInternalServerError)
		} else {
			HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func DBCreateMiddleware(db *gorm.DB, model interface{}, update bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(DBCreateHandleFunc(db, model, update))
	}
}

func DBGetFromDataBody(db *gorm.DB, model interface{}, arrayRefType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := GetPagination(r)
		pagination.Rows = reflect.New(reflect.TypeOf(arrayRefType).Elem()).Interface()
		data := r.Context().Value(DecodeBodyContextKey)
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

func DBGetFromDataParams(db *gorm.DB, model interface{}, arrayRefType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := GetPagination(r)
		pagination.Rows = reflect.New(reflect.TypeOf(arrayRefType).Elem()).Interface()
		data := r.Context().Value(DecodeParamsContextKey)
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

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetStudentEnrollments(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(JWTClaimContextKey).(*models.Student)
		db.Preload("modules").Joins("inner join enrollments e on modules.id = e.module_id").Where("e.student_id = ?", user.ID).Find(&modules)
		HandleResponseWithObject(w, modules, http.StatusOK)
	})
}

func GetStaffSupervisions(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(JWTClaimContextKey).(*models.Staff)
		db.Preload("modules").Joins("inner join supervisions e on modules.id = e.module_id").Where("e.staff_id = ?", user.ID).Find(&modules)
		HandleResponseWithObject(w, modules, http.StatusOK)
	})
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

func GetAssignedPairingsHandlerFunc(db *gorm.DB, muxVarKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pairings []models.Pairing
		assignmentId := mux.Vars(r)[muxVarKey]
		student := r.Context().Value(JWTClaimContextKey).(*models.Student)
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

func RandToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func ModulePermCheckMiddleware(db *gorm.DB, muxVarKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			moduleId := mux.Vars(r)[muxVarKey]
			var count int64

			switch data := r.Context().Value(JWTClaimContextKey).(type) {
			case *models.Student:
				db.Model(&models.Enrollment{}).Where("student_id = ? and module_id = ?", data.ID, moduleId).Count(&count)
			case *models.Staff:
				db.Model(&models.Supervision{}).Where("staff_id = ? and module_id = ?", data.ID, moduleId).Count(&count)
			}
			if count == 0 {
				HandleResponse(w, "Not enrolled in module", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
