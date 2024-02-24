package rest

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Woodfyn/Web-api/internal/service"
	mock_service "github.com/Woodfyn/Web-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_middlewareUserIdentity(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, token string)

	tests := []struct {
		name               string
		headerName         string
		headerValue        string
		token              string
		mockBehavior       mockBehavior
		expectedStatusCode int
		errWant            bool
		errMessage         string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUsers, token string) {
				r.EXPECT().ParseToken(token).Return("1", nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "invalide header name",
			headerName:         "",
			headerValue:        "Bearer token",
			token:              "token",
			mockBehavior:       func(r *mock_service.MockUsers, token string) {},
			expectedStatusCode: 401,
			errWant:            true,
			errMessage:         `{"message":"empty auth header"}`,
		},
		{
			name:               "invalid header value",
			headerName:         "Authorization",
			headerValue:        "Bearr token",
			token:              "token",
			mockBehavior:       func(r *mock_service.MockUsers, token string) {},
			expectedStatusCode: 401,
			errWant:            true,
			errMessage:         `{"message":"invalid auth header"}`,
		},
		{
			name:               "empty token",
			headerName:         "Authorization",
			headerValue:        "Bearer ",
			token:              "token",
			mockBehavior:       func(r *mock_service.MockUsers, token string) {},
			expectedStatusCode: 401,
			errWant:            true,
			errMessage:         `{"message":"token is empty"}`,
		},
		{
			name:        "parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUsers, token string) {
				r.EXPECT().ParseToken(token).Return("0", errors.New("invalid token"))
			},
			errWant:            true,
			expectedStatusCode: 401,
			errMessage:         `{"message":"invalid token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mock_service.NewMockUsers(ctrl)

			tt.mockBehavior(mockUsers, tt.token)

			service := service.Services{Users: mockUsers}
			handler := Handler{&service}

			// Init Endpoint
			r := gin.New()
			r.GET("/api", handler.userIdentity())

			// Init Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api", nil)
			req.Header.Add(tt.headerName, tt.headerValue)

			// Make Request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)

			if tt.errWant {
				assert.Equal(t, w.Body.String(), tt.errMessage)
			}
		})
	}
}
