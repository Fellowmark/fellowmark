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

	"github.com/alexedwards/argon2id"
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

func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func SetHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request (not recommended in production)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		//only allow GET, POST, and OPTIONS
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
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

func DecodeParams(r *http.Request, params interface{}) error {
	return SchemaDecoder.Decode(params, r.URL.Query())
}

func generatePointerWithSameType(ptr interface{}) interface{} {
	return reflect.New(reflect.TypeOf(ptr).Elem()).Interface()
}

func DecodeParamsMiddleware(refType interface{}) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parsedData := generatePointerWithSameType(refType)
			if err := DecodeParams(r, parsedData); err != nil {
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

func UserCreateHandleFunc(db *gorm.DB, userModel interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(DecodeBodyContextKey).(*models.User)
		result := db.Model(userModel).Create(data)
		if result.Error != nil {
			HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
		HandleResponseWithObject(w, data, http.StatusOK)
	}
}

func UserPasswordHashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(DecodeBodyContextKey).(*models.User)
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			HandleResponse(w, err.Error(), http.StatusInternalServerError)
		} else {
			user.Password = hash
			ctxWithUser := context.WithValue(r.Context(), DecodeBodyContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

func DBCreateHandleFunc(db *gorm.DB, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(DecodeBodyContextKey)
		db.Model(generatePointerWithSameType(data)).Clauses(
			clause.OnConflict{UpdateAll: true}).Create(data)
		output := generatePointerWithSameType(data)
		db.Model(data).Where(data).First(output)
		HandleResponseWithObject(w, output, http.StatusOK)
	}
}

func DBCreateMiddleware(db *gorm.DB, model interface{}, update bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(DBCreateHandleFunc(db, update))
	}
}

func DBGetFromDataBody(db *gorm.DB, model interface{}, arrayRefType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := GetPagination(r)
		pagination.Rows = reflect.New(reflect.TypeOf(arrayRefType).Elem()).Interface()
		data := r.Context().Value(DecodeBodyContextKey)
		scope := Paginate(db, func(tx *gorm.DB) *gorm.DB {
			return tx.Model(model)
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

func ModelDBScope(model interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(model)
	}
}
