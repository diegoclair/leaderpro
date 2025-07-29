package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/goswag"
	"github.com/diegoclair/leaderpro/infra/config"
	infraContract "github.com/diegoclair/leaderpro/infra/contract"
	"github.com/diegoclair/leaderpro/internal/application/service"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/airoute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/authroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/companyroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/dashboardroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/personroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/pingroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/shared"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/swaggerroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/userroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/leaderpro/internal/transport/rest/serverMiddleware"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	routes []routeutils.IRoute
	Router goswag.Echo
	cache  contract.CacheManager
}

func StartRestServer(ctx context.Context, cfg *config.Config, infra domain.Infrastructure, services *service.Apps, appName, port string) *Server {
	server := NewRestServer(services, cfg.GetAuthToken(), infra, appName)
	if port == "" {
		port = "5000"
	}

	infra.Logger().Infof(ctx, "About to start the application on port: %s...", port)

	go func() {
		if err := server.Start(port); err != nil {
			if err == http.ErrServerClosed {
				infra.Logger().Infof(ctx, "Server stopped")
			} else {
				infra.Logger().Errorf(ctx, "Server error: %v", err)
			}
		}
	}()

	return server
}

func NewRestServer(services *service.Apps, authToken infraContract.AuthToken, infra domain.Infrastructure, appName string) *Server {
	router := goswag.NewEcho(resterrors.GoSwagDefaultResponseErrors()...)
	router.Echo().Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	router.Echo().HTTPErrorHandler = func(err error, c echo.Context) {
		_ = routeutils.HandleError(c, err)
	}

	// shared auth helper
	authHelper := shared.NewAuthHelper(services.Auth, services.User, authToken)

	pingHandler := pingroute.NewHandler()
	authHandler := authroute.NewHandler(services.Auth, authToken, authHelper, infra.Logger())
	aiHandler := airoute.NewHandler(services.AI)
	companyHandler := companyroute.NewHandler(services.Company)
	dashboardHandler := dashboardroute.NewHandler(services.Dashboard)
	personHandler := personroute.NewHandler(services.Person)
	userHandler := userroute.NewHandler(services.User, authHelper)

	pingRoute := pingroute.NewRouter(pingHandler)
	authRoute := authroute.NewRouter(authHandler)
	aiRoute := airoute.NewRouter(aiHandler)
	companyRoute := companyroute.NewRouter(companyHandler)
	dashboardRoute := dashboardroute.NewRouter(dashboardHandler)
	personRoute := personroute.NewRouter(personHandler)
	userRoute := userroute.NewRouter(userHandler)

	swaggerRoute := swaggerroute.NewRouter(router.Echo())

	server := &Server{Router: router, cache: infra.CacheManager()}
	server.addRouters(authRoute)
	server.addRouters(aiRoute)
	server.addRouters(companyRoute)
	server.addRouters(dashboardRoute)
	server.addRouters(personRoute)
	server.addRouters(pingRoute)
	server.addRouters(swaggerRoute)
	server.addRouters(userRoute)
	server.registerAppRouters(authToken, services.Company)

	server.setupPrometheus(appName)

	return server
}

func (r *Server) addRouters(router routeutils.IRoute) {
	r.routes = append(r.routes, router)
}

func (r *Server) registerAppRouters(authToken infraContract.AuthToken, companyService contract.CompanyApp) {
	g := &routeutils.EchoGroups{}
	g.AppGroup = r.Router.Group("/")
	g.PrivateGroup = g.AppGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken, r.cache),
	)
	g.CompanyGroup = g.PrivateGroup.Group("",
		servermiddleware.CompanyOwnershipMiddleware(companyService),
	)

	for _, appRouter := range r.routes {
		appRouter.RegisterRoutes(g)
	}
}

func (r *Server) setupPrometheus(appName string) {
	p := echoprometheus.NewMiddleware(appName)
	r.Router.Echo().Use(p)
}

func (r *Server) Start(port string) error {
	return r.Router.Echo().Start(fmt.Sprintf(":%s", port))
}
