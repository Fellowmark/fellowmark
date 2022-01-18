package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StaffRoute struct {
	DB *gorm.DB
}

func (ur StaffRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
	ur.CreatePrivilegedRouter(route.PathPrefix("/module/{moduleId}").Subrouter())
}

func (ur StaffRoute) CreateAuthRouter(route *mux.Router) {
	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(utils.DecodeParamsMiddleware(&models.Staff{}))
	loginRoute.Use(ur.StaffLoginMiddleware)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Staff")).Methods(http.MethodGet)

	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(utils.DecodeBodyMiddleware(&models.Staff{}))
	signUpRoute.Use(utils.SanitizeDataMiddleware())
	signUpRoute.Use(ur.PasswordHash)
	signUpRoute.HandleFunc("/signup", utils.DBCreateHandleFunc(ur.DB, &models.Staff{}, false)).Methods(http.MethodPost)
}

func (ur StaffRoute) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Staff", &models.Staff{}))
	route.Use(utils.SupervisionCheckMiddleware(ur.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	ur.CreatePairingsRoute(route.PathPrefix("/pairing").Subrouter())
	ur.GetPairingsRoute(route.PathPrefix("/pairing").Subrouter())
}

func (ur StaffRoute) CreatePairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}))
	route.HandleFunc("/initialize", ur.InitializePairings).Methods(http.MethodPost)
	route.HandleFunc("/assign", ur.AssignPairings).Methods(http.MethodPost)
}

func (ur StaffRoute) GetPairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Pairing{}))
	route.HandleFunc("", utils.DBGetFromDataParams(ur.DB, &models.Pairing{}, &[]models.Pairing{})).Methods(http.MethodGet)
}
