package submissions

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type FileserverRoute struct {
	DB            *gorm.DB
	UploadPath    string
	MaxUploadSize int64
}

func (fr FileserverRoute) CreateRouters(route *mux.Router) {
	fr.CreatePrivilegedRoute(route.NewRoute().Subrouter())
}

func (fr FileserverRoute) CreatePrivilegedRoute(route *mux.Router) {
	// if os.Getenv("RUN_ENV") == "production" {
	// route.Use(utils.ValidateJWTMiddleware("Student", "claims", &models.Student{}))
	// route.Use(utils.ModulePermCheckMiddleware(fr.DB, "claims", "moduleId"))
	// }

	fr.CreateDownloadRoute(route.NewRoute().Subrouter())
	fr.CreateUploadRoute(route.NewRoute().Subrouter())
}

func (fr FileserverRoute) CreateDownloadRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}, "submission"))
	route.Use(fr.FileAuthMiddleware("claims", "submission", false))
	route.Use(UpdateFilePathMiddleware(fr.DB, "submission", "file"))
	route.HandleFunc("", utils.DownloadHandlerFunc("submission", "file")).Methods(http.MethodGet)
}

func (fr FileserverRoute) CreateUploadRoute(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", "claims", &models.Student{}))
	route.Use(utils.EnrollmentCheckMiddleware(fr.DB, "claims", "moduleId"))
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}, "submission"))
	route.Use(fr.FileAuthMiddleware("claims", "submission", true))
	route.Use(fr.GetSubmissionMiddleware("submission"))
	route.Use(utils.UploadMiddleware(fr.UploadPath, "file", fr.MaxUploadSize))
	route.Use(UpdateSubmissionContentFile("submission", "file"))
	route.HandleFunc("", utils.DBCreateHandleFunc(fr.DB, &models.Submission{}, "submission", true)).Methods(http.MethodPost)
}

func UpdateFilePathMiddleware(db *gorm.DB, dbDataContextInKey string, filePathContextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(dbDataContextInKey).(*models.Submission)
			db.Where(data).First(&data)
			ctxWithPath := context.WithValue(r.Context(), filePathContextOutKey, data.ContentFile)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func UpdateSubmissionContentFile(dbDataContextInKey string, filePathContextInKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(dbDataContextInKey).(*models.Submission)
			filePath := r.Context().Value(filePathContextInKey).(string)
			data.ContentFile = filePath
			ctxWithPath := context.WithValue(r.Context(), dbDataContextInKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (fr FileserverRoute) FileAuthMiddleware(userDataContextInKey string, dbDataContextOutKey string, isMarkee bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(dbDataContextOutKey).(*models.Submission)
			user := r.Context().Value(userDataContextInKey)

			switch marker := user.(type) {
			case *models.Student:
				if marker.ID != data.StudentID {
					if isMarkee {
						utils.HandleResponse(w, "Bad request. Incorrect permissions", http.StatusUnauthorized)
						return
					}
					var isMarker int64
					fr.DB.Model(&models.Submission{}).
						Joins("questions").Joins("pairings").
						Where("pairings.marker_id = ? AND questions.id = ? AND pairings.student_id = ?", marker.ID, data.QuestionID, data.StudentID).
						Count(&isMarker)
					if isMarker != 0 {
						utils.HandleResponse(w, "Bad request. Incorrect permissions", http.StatusUnauthorized)
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (fr FileserverRoute) GetSubmissionMiddleware(dbDataContextOutKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(dbDataContextOutKey).(*models.Submission)
			fr.DB.Model(&models.Submission{}).Where("student_id = ? AND question_id = ?", data.StudentID, data.QuestionID).Find(data)
			ctxWithPath := context.WithValue(r.Context(), dbDataContextOutKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}
