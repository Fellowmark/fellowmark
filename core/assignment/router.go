package assignment

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AssignmentRoute struct {
	DB *gorm.DB
}

func (ur AssignmentRoute) CreateAssigmentRouter(route *mux.Router) {
	createAssignmentRoute := route.NewRoute().Subrouter()
	createAssignmentRoute.Use(utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment"))
	createAssignmentRoute.Use(utils.SanitizeDataMiddleware("assignment"))
	createAssignmentRoute.HandleFunc("", utils.DBCreateHandleFunc(ur.DB, &models.Assignment{}, "assigment", true)).Methods(http.MethodGet)

	createQuestionRoute := route.NewRoute().Subrouter()
	createQuestionRoute.Use(utils.DecodeBodyMiddleware(&models.Question{}, "question"))
	createQuestionRoute.Use(utils.SanitizeDataMiddleware("question"))
	createQuestionRoute.HandleFunc("/question", utils.DBCreateHandleFunc(ur.DB, &models.Question{}, "question", true)).Methods(http.MethodGet)

	createRubricRoute := route.NewRoute().Subrouter()
	createRubricRoute.Use(utils.DecodeBodyMiddleware(&models.Rubric{}, "rubric"))
	createRubricRoute.Use(utils.SanitizeDataMiddleware("rubric"))
	createRubricRoute.HandleFunc("/rubric", utils.DBCreateHandleFunc(ur.DB, &models.Rubric{}, "rubric", true)).Methods(http.MethodGet)
}
