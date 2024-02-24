package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/service"
	mock_service "github.com/Woodfyn/Web-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	normalName     = "test"
	normalEmail    = "test@example.com"
	normalPassword = "test1234"

	accessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDgyNzA0MDEsInN1YiI6IjUifQ.KARe4sjN4bofM-bqXHCDCLIaLqZyqv16EGP6-dK0SMw"
	refreshToken = "g29837rb8bf914vbg9f13bg981534gg98gvb13895vb198534bv"
)

func TestHandler_authSignUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, user domain.SignUpInput)

	tests := []struct {
		name         string
		requestBody  string
		serviceInput domain.SignUpInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:        "ok",
			requestBody: fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s"}`, normalName, normalEmail, normalPassword),
			serviceInput: domain.SignUpInput{
				Name:     normalName,
				Email:    normalEmail,
				Password: normalPassword,
			},
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {
				r.EXPECT().SignUp(user).Return(nil)
			},
			statusCode: 200,
		},
		{
			name:         "miss name",
			requestBody:  fmt.Sprintf(`{"email": "%s", "password": "%s"}`, normalEmail, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:         "invalid name",
			requestBody:  fmt.Sprintf(`{"name": "t", "email": "%s", "password": "%s"}`, normalEmail, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Name' Error:Field validation for 'Name' failed on the 'gte' tag"}`,
		},
		{
			name:         "miss email",
			requestBody:  fmt.Sprintf(`{"name": "%s", "password": "%s"}`, normalName, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:         "invalid email",
			requestBody:  fmt.Sprintf(`{"name": "%s", "email": "test", "password": "%s"}`, normalName, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{
			name:         "miss password",
			requestBody:  fmt.Sprintf(`{"name": "%s", "email": "%s"}`, normalName, normalEmail),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:         "invalid password",
			requestBody:  fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "test"}`, normalName, normalEmail),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignUpInput.Password' Error:Field validation for 'Password' failed on the 'gte' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mock_service.NewMockUsers(ctrl)

			tt.mockBehavior(mockUsers, tt.serviceInput)

			services := &service.Services{Users: mockUsers}
			handler := &Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)

			if tt.name != "ok" {
				assert.Equal(t, w.Body.String(), tt.responseBody)
			}
		})
	}
}

func TestHandler_authSignIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string)

	tests := []struct {
		name         string
		requestBody  string
		serviceInput domain.SignInInput
		accessToken  string
		refreshToken string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:        "ok",
			requestBody: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, normalEmail, normalPassword),
			serviceInput: domain.SignInInput{
				Email:    normalEmail,
				Password: normalPassword,
			},
			accessToken:  accessToken,
			refreshToken: refreshToken,
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string) {
				r.EXPECT().SignIn(user).Return(accessToken, refreshToken, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"token":"%s"}`, accessToken),
		},
		{
			name:         "miss email",
			requestBody:  fmt.Sprintf(`{"password": "%s"}`, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignInInput.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:         "invalid email",
			requestBody:  fmt.Sprintf(`{"email": "test", "password": "%s"}`, normalPassword),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignInInput.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{
			name:         "miss password",
			requestBody:  fmt.Sprintf(`{"email": "%s"}`, normalEmail),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:         "invalid password",
			requestBody:  fmt.Sprintf(`{"password": "1234", "email": "%s"}`, normalEmail),
			mockBehavior: func(r *mock_service.MockUsers, user domain.SignInInput, accessToken, refreshToken string) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'SignInInput.Password' Error:Field validation for 'Password' failed on the 'gte' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mock_service.NewMockUsers(ctrl)

			tt.mockBehavior(mockUsers, tt.serviceInput, tt.accessToken, tt.refreshToken)

			services := &service.Services{Users: mockUsers}
			handler := &Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/sign-in", handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/sign-in", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			if tt.name == "ok" {
				res := w.Result().Cookies()[0].Value
				headerParts := strings.Split(res, " ")
				assert.Equal(t, headerParts[1], tt.refreshToken)
			}

			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

// This test is nit working, because it does not return a cookie on response, but on rilly method all worked
func TestHandler_authRefresh(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, accessToken, refreshToken string)

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:         "ok",
			refreshToken: refreshToken,
			mockBehavior: func(r *mock_service.MockUsers, accessToken, refreshToken string) {
				r.EXPECT().RefreshTokens(refreshToken).Return(accessToken, refreshToken, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"token":"%s"}`, accessToken),
		},
		{
			name:         "cookie missing",
			mockBehavior: func(r *mock_service.MockUsers, accessToken, refreshToken string) {},
			statusCode:   400,
			responseBody: `{"message":"key: 'RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mock_service.NewMockUsers(ctrl)

			tt.mockBehavior(mockUsers, tt.accessToken, tt.refreshToken)

			services := &service.Services{Users: mockUsers}
			handler := &Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/refresh", handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/refresh", nil)
			req.AddCookie(&http.Cookie{
				Name:     "Authorization",
				Value:    tt.refreshToken,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
			})

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			if tt.refreshToken != "" {
				header := w.Result().Header
				cookie := header["Cookie"]
				token := strings.Split(cookie[0], "=")
				assert.Equal(t, token[1], tt.refreshToken)
			}

			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

// This test is nit working, because it does not return a cookie on response, but on rilly method all worked
func TestHandler_authLogOut(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, refreshToken string)

	tests := []struct {
		name         string
		refreshToken string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:         "ok",
			refreshToken: refreshToken,
			mockBehavior: func(r *mock_service.MockUsers, refreshToken string) {
				r.EXPECT().LogOut(refreshToken).Return(nil)
			},
			statusCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mock_service.NewMockUsers(ctrl)

			tt.mockBehavior(mockUsers, tt.refreshToken)

			services := &service.Services{Users: mockUsers}
			handler := &Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/log-out", handler.logOut)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/log-out", nil)

			req.AddCookie(&http.Cookie{
				Name:     "Authorization",
				Value:    tt.refreshToken,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
			})

			// Make Request
			r.ServeHTTP(w, req)

			//Assert
			if tt.refreshToken != "" {
				res := w.Result().Cookies()[0].Value
				assert.Equal(t, res, "")
			}

			assert.Equal(t, w.Body.String(), tt.responseBody)
			assert.Equal(t, w.Code, tt.statusCode)
		})
	}
}
