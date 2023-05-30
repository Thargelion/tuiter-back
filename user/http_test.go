package user

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tuiter.com/api/kit"
)

type mockTuiterTime struct {
	mock.Mock
}

func (m *mockTuiterTime) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindUserByID(ctx context.Context, ID string) (*User, error) {
	args := m.Called(ID)
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepository) Create(ctx context.Context, user *User) (*User, error) {
	args := m.Called(user)
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepository) Search(ctx context.Context, query map[string]interface{}) ([]*User, error) {
	args := m.Called(query)
	return args.Get(0).([]*User), args.Error(1)
}

type UserHttpSuite struct {
	suite.Suite
	writer  *httptest.ResponseRecorder
	request *http.Request
	repo    *mockRepository
	tt      kit.Time
}

func (suite *UserHttpSuite) SetupTest() {
	suite.writer = httptest.NewRecorder()
	suite.request = httptest.NewRequest("GET", "/", nil)
	suite.repo = &mockRepository{}
	suite.tt = &mockTuiterTime{}
}

func (suite *UserHttpSuite) TestRouterFindUserOK() {
	// Given
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", "username")
	request := suite.request.WithContext(context.WithValue(suite.request.Context(), chi.RouteCtxKey, chiContext))
	expected := &User{}
	suite.repo.On("FindUserByID", "username").Return(expected, nil)
	subject := NewUserRouter(suite.tt, suite.repo)
	// When
	subject.FindUserByID(suite.writer, request)
	// Then
	suite.repo.AssertExpectations(suite.T())
	res := suite.writer.Result()
	assert.Equal(suite.T(), 200, res.StatusCode)
}

func (suite *UserHttpSuite) TestRouterFindUserNotFound() {
	// Given
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", "username")
	request := suite.request.WithContext(context.WithValue(suite.request.Context(), chi.RouteCtxKey, chiContext))
	suite.repo.On("FindUserByID", "username").Return(&User{}, errors.New("x.x"))
	subject := NewUserRouter(suite.tt, suite.repo)
	// When
	subject.FindUserByID(suite.writer, request)
	// Then
	suite.repo.AssertExpectations(suite.T())
	res := suite.writer.Result()
	assert.Equal(suite.T(), 400, res.StatusCode)
}

func TestHtpFindSuite(t *testing.T) {
	suite.Run(t, new(UserHttpSuite))
}
