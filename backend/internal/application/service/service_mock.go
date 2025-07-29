package service

import (
	"testing"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/infra/configmock"
	"github.com/diegoclair/leaderpro/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type allMocks struct {
	mockDataManager *mocks.MockDataManager

	mockAuthRepo *mocks.MockAuthRepo
	mockUserRepo *mocks.MockUserRepo

	mockCacheManager *mocks.MockCacheManager
	mockCrypto       *mocks.MockCrypto
	mockValidator    validator.Validator
	mockLogger       logger.Logger

	mockUserSvc *mocks.MockUserApp
	
	mockAIProvider *mocks.MockAIProvider

	mockDomain *mocks.MockInfrastructure
}

func newServiceTestMock(t *testing.T) (m allMocks, ctrl *gomock.Controller) {
	t.Helper()
	cfg := configmock.New()

	ctrl = gomock.NewController(t)

	dm := mocks.NewMockDataManager(ctrl)

	userRepo := mocks.NewMockUserRepo(ctrl)
	dm.EXPECT().User().Return(userRepo).AnyTimes()

	authRepo := mocks.NewMockAuthRepo(ctrl)
	dm.EXPECT().Auth().Return(authRepo).AnyTimes()

	cm := cfg.GetCacheManager(ctrl)
	crypto := cfg.GetCrypto(ctrl)
	log := cfg.GetLogger()
	v := cfg.GetValidator(t)

	userSvc := mocks.NewMockUserApp(ctrl)
	aiProvider := mocks.NewMockAIProvider(ctrl)

	domainMock := mocks.NewMockInfrastructure(ctrl)
	domainMock.EXPECT().DataManager().Return(dm).AnyTimes()
	domainMock.EXPECT().Logger().Return(log).AnyTimes()
	domainMock.EXPECT().CacheManager().Return(cm).AnyTimes()
	domainMock.EXPECT().Crypto().Return(crypto).AnyTimes()
	domainMock.EXPECT().Validator().Return(v).AnyTimes()

	m = allMocks{
		mockDataManager:  dm,
		mockUserRepo:     userRepo,
		mockCacheManager: cm,
		mockAuthRepo:     authRepo,
		mockCrypto:       crypto,
		mockUserSvc:      userSvc,
		mockAIProvider:   aiProvider,
		mockDomain:       domainMock,
		mockValidator:    v,
		mockLogger:       log,
	}

	// validate func New
	s, err := New(domainMock, aiProvider, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, s)

	return
}
