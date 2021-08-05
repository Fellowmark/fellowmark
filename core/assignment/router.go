package assignment

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AssignmentRoute struct {
	DB *gorm.DB
}

func (ar AssignmentRoute) CreateRouters(route *mux.Router) {
	ar.GetAssignmentsRoute(route.NewRoute().Subrouter())
	ar.GetQuestionsRoute(route.PathPrefix("/question").Subrouter())
	ar.GetRubricsRoute(route.PathPrefix("/rubric").Subrouter())
	ar.CreatePrivilegedRouters(route.NewRoute().Subrouter())
}

func (ar AssignmentRoute) CreatePrivilegedRouters(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Staff", "claims"))
	}

	ar.CreateAssignmentRouter(route.NewRoute().Subrouter())
	ar.CreateQuestionsRouter(route.PathPrefix("/question").Subrouter())
	ar.CreateRubricsRoute(route.PathPrefix("/rubric").Subrouter())
}

func (ar AssignmentRoute) CreateAssignmentRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment"))
	route.Use(utils.SanitizeDataMiddleware("assignment"))
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Assignment{}, "assignment", true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) CreateQuestionsRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Question{}, "question"))
	route.Use(utils.SanitizeDataMiddleware("question"))
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Question{}, "question", true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) CreateRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Rubric{}, "rubric"))
	route.Use(utils.SanitizeDataMiddleware("rubric"))
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Rubric{}, "rubric", true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) GetAssignmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment"))
	route.HandleFunc("", utils.DBGetFromData(ar.DB, &models.Assignment{}, "assignment", &[]models.Supervision{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetQuestionsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Question{}, "question"))
	route.HandleFunc("", utils.DBGetFromData(ar.DB, &models.Question{}, "question", &[]models.Question{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Rubric{}, "rubric"))
	route.HandleFunc("", utils.DBGetFromData(ar.DB, &models.Rubric{}, "rubric", &[]models.Rubric{})).Methods(http.MethodGet)
}
