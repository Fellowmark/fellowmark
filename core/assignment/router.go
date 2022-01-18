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

func (ar AssignmentRoute) CreateRouters(route *mux.Router) {
	ar.GetAssignmentsRoute(route.NewRoute().Subrouter())
	ar.GetQuestionsRoute(route.PathPrefix("/question").Subrouter())
	ar.GetRubricsRoute(route.PathPrefix("/rubric").Subrouter())
	ar.CreatePrivilegedRouters(route.NewRoute().Subrouter())
	ar.GetStudentPairingsRoute(route.PathPrefix("/{assignmentId}/pairs").Subrouter())
	ar.GetSubmissionRoute(route.PathPrefix("/submission").Subrouter())
}

func (ar AssignmentRoute) CreatePrivilegedRouters(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Staff", &models.Staff{}))

	ar.CreateAssignmentRouter(route.NewRoute().Subrouter())
	ar.CreateQuestionsRouter(route.PathPrefix("/question").Subrouter())
	ar.CreateRubricsRoute(route.PathPrefix("/rubric").Subrouter())
}

func (ar AssignmentRoute) CreateAssignmentRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}))
	route.Use(ar.AuthorizedToMutateAssignment(ar.resolveAssignmentIdFromContext))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Assignment{}, true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) CreateQuestionsRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Question{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Question{}, true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) CreateRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Rubric{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", utils.DBCreateHandleFunc(ar.DB, &models.Rubric{}, true)).Methods(http.MethodPost)
}

func (ar AssignmentRoute) GetAssignmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Assignment{}))
	route.HandleFunc("", utils.DBGetFromDataParams(ar.DB, &models.Assignment{}, &[]models.Assignment{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetQuestionsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Question{}))
	route.HandleFunc("", utils.DBGetFromDataParams(ar.DB, &models.Question{}, &[]models.Question{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Rubric{}))
	route.HandleFunc("", utils.DBGetFromDataParams(ar.DB, &models.Rubric{}, &[]models.Rubric{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetStudentPairingsRoute(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.ValidateAssignmentIdMiddlware(ar.DB, "assignmentId", "moduleId"))
	route.Use(utils.EnrollmentCheckMiddleware(ar.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.GetAssignedPairingsHandlerFunc(ar.DB, "assignmentId")).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.HandleFunc("", utils.DBGetFromDataParams(ar.DB, &models.Submission{}, &[]models.Submission{})).Methods(http.MethodGet)
}

func (ar AssignmentRoute) GetAllAssignmentPairings(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Pairing{}))
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.ValidateAssignmentIdMiddlware(ar.DB, "assignmentId", "moduleId"))
	route.Use(utils.SupervisionCheckMiddleware(ar.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.DBGetFromDataParams(ar.DB, &models.Pairing{}, &[]models.Pairing{})).Methods(http.MethodGet)
}
