package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"tuiter.com/api/api"
	"tuiter.com/api/pkg"
	"tuiter.com/api/pkg/user"
)

var errMock = errors.New("mock error")

type mockTuiterTime struct {
	mock.Mock
}

func (m *mockTuiterTime) Now() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

type mockUseCases struct {
	mock.Mock
}

func (m *mockUseCases) FindUserByID(_ context.Context, id string) (*user.User, error) {
	args := m.Called(id)

	return args.Get(0).(*user.User), args.Error(1) //nolint:wrapcheck
}

func (m *mockUseCases) Create(_ context.Context, u *user.User) (*user.User, error) {
	args := m.Called(u)

	return args.Get(0).(*user.User), args.Error(1) //nolint:wrapcheck
}

func (m *mockUseCases) Search(_ context.Context, query map[string]interface{}) ([]*user.User, error) {
	args := m.Called(query)

	return args.Get(0).([]*user.User), args.Error(1) //nolint:wrapcheck
}

type UserHTTPSuite struct {
	suite.Suite
	writer       *httptest.ResponseRecorder
	request      *http.Request
	useCasesMock *mockUseCases
	tt           pkg.Time
}

func (suite *UserHTTPSuite) SetupTest() {
	suite.writer = httptest.NewRecorder()
	suite.request = httptest.NewRequest(http.MethodGet, "/", nil)
	suite.useCasesMock = &mockUseCases{}
	suite.tt = &mockTuiterTime{}
}

func (suite *UserHTTPSuite) TestRouterFindUserOK() {
	// Given
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", "username")
	request := suite.request.WithContext(context.WithValue(suite.request.Context(), chi.RouteCtxKey, chiContext))
	expected := &user.User{}
	suite.useCasesMock.On("FindUserByID", "username").Return(expected, nil)
	subject := api.NewUserRouter(suite.useCasesMock)
	// When
	subject.FindUserByID(suite.writer, request)
	// Then
	suite.useCasesMock.AssertExpectations(suite.T())

	res := suite.writer.Result()
	defer res.Body.Close()
	assert.Equal(suite.T(), 200, res.StatusCode)
}

func (suite *UserHTTPSuite) TestRouterFindUserNotFound() {
	// Given
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", "username")
	request := suite.request.WithContext(context.WithValue(suite.request.Context(), chi.RouteCtxKey, chiContext))
	suite.useCasesMock.On("FindUserByID", "username").Return(&user.User{}, errMock)
	subject := api.NewUserRouter(suite.useCasesMock)
	// When
	subject.FindUserByID(suite.writer, request)
	// Then
	suite.useCasesMock.AssertExpectations(suite.T())

	res := suite.writer.Result()
	defer res.Body.Close()
	assert.Equal(suite.T(), 400, res.StatusCode)
}

func TestHtpFindSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserHTTPSuite))
}
