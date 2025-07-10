package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/goswag"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/infra/auth"
	"github.com/diegoclair/leaderpro/infra/configmock"
	"github.com/diegoclair/leaderpro/infra/contract"
	infraMocks "github.com/diegoclair/leaderpro/infra/mocks"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/authroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/companyroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/personroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routes/userroute"
	"github.com/diegoclair/leaderpro/internal/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/leaderpro/internal/transport/rest/serverMiddleware"
	"github.com/diegoclair/leaderpro/mocks"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"go.uber.org/mock/gomock"
)

type SvcMocks struct {
	UserAppMock    *mocks.MockUserApp
	AuthAppMock    *mocks.MockAuthApp
	PersonAppMock  *mocks.MockPersonApp
	CompanyAppMock *mocks.MockCompanyApp
	AuthTokenMock  *infraMocks.MockAuthToken
	CacheMock      *mocks.MockCacheManager
}

func GetServerTest(t *testing.T) (m SvcMocks, server goswag.Echo, ctrl *gomock.Controller) {
	t.Helper()

	ctrl = gomock.NewController(t)
	m = SvcMocks{
		UserAppMock:    mocks.NewMockUserApp(ctrl),
		AuthAppMock:    mocks.NewMockAuthApp(ctrl),
		PersonAppMock:  mocks.NewMockPersonApp(ctrl),
		CompanyAppMock: mocks.NewMockCompanyApp(ctrl),
		AuthTokenMock:  infraMocks.NewMockAuthToken(ctrl),
		CacheMock:      mocks.NewMockCacheManager(ctrl),
	}

	server = goswag.NewEcho()
	server.Echo().HTTPErrorHandler = func(err error, c echo.Context) {
		_ = routeutils.HandleError(c, err)
	}
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(getTestTokenMaker(t), m.CacheMock),
	)

	g := &routeutils.EchoGroups{
		AppGroup:     appGroup,
		PrivateGroup: privateGroup,
	}

	userHandler := userroute.NewHandler(m.UserAppMock)
	userRoute := userroute.NewRouter(userHandler)
	authHandler := authroute.NewHandler(m.AuthAppMock, m.AuthTokenMock)
	authRoute := authroute.NewRouter(authHandler)
	personHandler := personroute.NewHandler(m.PersonAppMock)
	personRoute := personroute.NewRouter(personHandler)
	companyHandler := companyroute.NewHandler(m.CompanyAppMock)
	companyRoute := companyroute.NewRouter(companyHandler)

	userRoute.RegisterRoutes(g)
	authRoute.RegisterRoutes(g)
	personRoute.RegisterRoutes(g)
	companyRoute.RegisterRoutes(g)
	return
}

var (
	tokenMaker contract.AuthToken
	onceToken  sync.Once
)

func getTestTokenMaker(t *testing.T) contract.AuthToken {
	t.Helper()

	onceToken.Do(func() {
		cfg := configmock.New()
		var err error

		cfg.Auth.AccessTokenDuration = 2 * time.Second
		cfg.Auth.RefreshTokenDuration = 2 * time.Second

		tokenMaker, err = auth.NewAuthToken(cfg.Auth.AccessTokenDuration,
			cfg.Auth.RefreshTokenDuration,
			cfg.Auth.PasetoSymmetricKey,
			cfg.GetLogger(),
		)
		require.NoError(t, err)
		require.NotEmpty(t, tokenMaker)
	})
	return tokenMaker
}

var (
	userUUID    = uuid.NewV4().String()
	sessionUUID = uuid.NewV4().String()
)

func AddAuthorization(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
	t.Helper()

	token := addAuthorizationWithNoCache(ctx, t, req)
	m.CacheMock.EXPECT().GetString(gomock.Any(), token).Return("", nil).Times(1)
}

func addAuthorizationWithNoCache(ctx context.Context, t *testing.T, req *http.Request) (token string) {
	t.Helper()

	tokenMaker := getTestTokenMaker(t)

	token, _, err := tokenMaker.CreateAccessToken(ctx, contract.TokenPayloadInput{UserUUID: userUUID, SessionUUID: sessionUUID})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	req.Header.Set(infra.TokenKey.String(), token)
	return token
}

func GetTestContext(t *testing.T, req *http.Request, w http.ResponseWriter, authEndpoint bool) context.Context {
	t.Helper()

	c := echo.New().NewContext(req, w)
	if authEndpoint {
		c.Set(infra.UserUUIDKey.String(), userUUID)
		c.Set(infra.SessionKey.String(), sessionUUID)
	}
	return routeutils.GetContext(c)
}

type PrivateEndpointTest struct {
	Name          string
	Body          any
	SetupAuth     func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks)
	BuildMocks    func(ctx context.Context, m SvcMocks, body any)
	CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

var PrivateEndpointValidations = []PrivateEndpointTest{
	{
		Name: "Should return error when token is invalid",
		SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
			req.Header.Set(infra.TokenKey.String(), "invalid token")
		},
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			if recorder.Code != http.StatusUnauthorized {
				t.Errorf("Expected status code %d. Got %d", http.StatusUnauthorized, recorder.Code)
				t.Errorf("Response body: %s", recorder.Body.String())
			}
		},
	},
	{
		Name: "Should return error when token is missing",
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			if recorder.Code != http.StatusUnauthorized {
				t.Errorf("Expected status code %d. Got %d", http.StatusUnauthorized, recorder.Code)
				t.Errorf("Response body: %s", recorder.Body.String())
			}
		},
	},
	{
		Name: "Should return error when token is invalid",
		SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
			addAuthorizationWithNoCache(ctx, t, req)
		},
		BuildMocks: func(ctx context.Context, m SvcMocks, body any) {
			m.CacheMock.EXPECT().GetString(gomock.Any(), gomock.Any()).Return("invalid", nil).Times(1)
		},
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	},
}
