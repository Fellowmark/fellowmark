package assignment

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AssignmentController struct {
	DB *gorm.DB
}

func (controller AssignmentController) CreateRouters(route *mux.Router) {
	controller.CreatePrivilegedRouters(route.NewRoute().Subrouter())

	controller.GetQuestionsRoute(route.PathPrefix("/question").Subrouter())
	controller.GetRubricsRoute(route.PathPrefix("/rubric").Subrouter())
	controller.GetStudentPairingsRoute(route.PathPrefix("/{assignmentId}/pairs").Subrouter())
	controller.GetSubmissionRoute(route.PathPrefix("/submission").Subrouter())
}

func (controller AssignmentController) CreatePrivilegedRouters(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())

	controller.GetAssignmentsRoute(route.NewRoute().Subrouter())
	controller.CreateAssignmentRouter(route.NewRoute().Subrouter())
	controller.CreateQuestionsRouter(route.PathPrefix("/question").Subrouter())
	controller.CreateRubricsRoute(route.PathPrefix("/rubric").Subrouter())
}

func (controller AssignmentController) CreateAssignmentRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.Use(controller.CreateAssignmentPermissionCheck())
	route.HandleFunc("", utils.DBCreateHandleFunc(controller.DB, true)).Methods(http.MethodPost)
}

func (controller AssignmentController) CreateQuestionsRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Question{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.Use(controller.CreateAssignmentPermissionCheck())
	route.HandleFunc("", utils.DBCreateHandleFunc(controller.DB, true)).Methods(http.MethodPost)
}

func (controller AssignmentController) CreateRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Rubric{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", utils.DBCreateHandleFunc(controller.DB, true)).Methods(http.MethodPost)
}

func (controller AssignmentController) GetAssignmentsRoute(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.DecodeParamsMiddleware(&models.Assignment{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Assignment{}, &[]models.Assignment{})).Methods(http.MethodGet)
}

func (controller AssignmentController) GetQuestionsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Question{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Question{}, &[]models.Question{})).Methods(http.MethodGet)
}

func (controller AssignmentController) GetRubricsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Rubric{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Rubric{}, &[]models.Rubric{})).Methods(http.MethodGet)
}

func (controller AssignmentController) GetStudentPairingsRoute(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.ValidateAssignmentIdMiddlware(controller.DB, "assignmentId", "moduleId"))
	route.Use(utils.EnrollmentCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.GetAssignedPairingsHandlerFunc(controller.DB, "assignmentId")).Methods(http.MethodGet)
}

func (controller AssignmentController) GetSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Submission{}, &[]models.Submission{})).Methods(http.MethodGet)
}

func (controller AssignmentController) GetAllAssignmentPairings(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Pairing{}))
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.ValidateAssignmentIdMiddlware(controller.DB, "assignmentId", "moduleId"))
	// route.Use(utils.SupervisionCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Pairing{}, &[]models.Pairing{})).Methods(http.MethodGet)
}
