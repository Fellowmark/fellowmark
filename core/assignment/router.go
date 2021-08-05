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
	ar.CreatePrivilegedRouters(route.NewRoute().Subrouter())
}

func (ar AssignmentRoute) CreatePrivilegedRouters(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Staff", "claims"))
	}

	ar.CreateAssigmentRouter(route)
	ar.CreateQuestionsRouter(route.PathPrefix("/question").Subrouter())
	ar.CreateRubricsRoute(route.PathPrefix("/rubric").Subrouter())
}

func (ar AssignmentRoute) CreateAssigmentRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment"))
	route.Use(utils.SanitizeDataMiddleware("assignment"))
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Assignment{}, "assigment", true)).Methods(http.MethodGet)
}

func (ar AssignmentRoute) CreateQuestionsRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Question{}, "question"))
	route.Use(utils.SanitizeDataMiddleware("question"))
	route.HandleFunc("/question", utils.DBCreateHandleFunc(ar.DB, &models.Question{}, "question", true)).Methods(http.MethodGet)
}

func (ar AssignmentRoute) CreateRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Rubric{}, "rubric"))
	route.Use(utils.SanitizeDataMiddleware("rubric"))
	route.HandleFunc("/rubric", utils.DBCreateHandleFunc(ar.DB, &models.Rubric{}, "rubric", true)).Methods(http.MethodGet)
}
