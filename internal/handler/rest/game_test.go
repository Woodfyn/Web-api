package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/Woodfyn/Web-api/internal/service"
	mock_service "github.com/Woodfyn/Web-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	normalTitle      = "test"
	normalGenre      = "test"
	normalEvaluation = 10
)

func TestHandler_gameAddGame(t *testing.T) {
	type mockBehavior func(*mock_service.MockGames, domain.Game)

	tests := []struct {
		name         string
		requestBody  string
		serviceInput domain.Game
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:        "ok",
			requestBody: fmt.Sprintf(`{"title": "%s", "genre": "%s", "evaluation": %d}`, normalTitle, normalGenre, normalEvaluation),
			serviceInput: domain.Game{
				Title:      normalTitle,
				Genre:      normalGenre,
				Evaluation: normalEvaluation,
			},
			mockBehavior: func(r *mock_service.MockGames, game domain.Game) {
				r.EXPECT().Create(game).Return(nil)
			},
			statusCode: 200,
		},
		{
			name:         "miss title",
			requestBody:  fmt.Sprintf(`{"genre": "%s", "evaluation": %d}`, normalGenre, normalEvaluation),
			mockBehavior: func(r *mock_service.MockGames, game domain.Game) {},
			statusCode:   400,
			responseBody: `{"message":"Key: 'Game.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGames := mock_service.NewMockGames(ctrl)
			tt.mockBehavior(mockGames, tt.serviceInput)

			service := service.Services{Games: mockGames}
			handler := Handler{services: &service}

			// Init Endpoint
			r := gin.New()
			r.POST("/game", handler.addGame)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/game", bytes.NewBufferString(tt.requestBody))

			// Write Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			if tt.name != "ok" {
				assert.Equal(t, tt.responseBody, w.Body.String())
			}
		})
	}
}

func TestHandler_gameGetAllGames(t *testing.T) {
	type mockBehavior func(*mock_service.MockGames, []domain.Game)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		responseData []domain.Game
		responseBody getAllGameResponse
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockGames, games []domain.Game) {
				r.EXPECT().GetAll().Return(games, nil)
			},
			statusCode: 200,
			responseData: []domain.Game{
				{
					Id:         1,
					Title:      "test",
					Genre:      "test",
					Evaluation: 10,
				},
				{
					Id:         2,
					Title:      "test",
					Genre:      "test",
					Evaluation: 10,
				},
			},
			responseBody: getAllGameResponse{
				Data: []domain.Game{
					{
						Id:         1,
						Title:      "test",
						Genre:      "test",
						Evaluation: 10,
					},
					{
						Id:         2,
						Title:      "test",
						Genre:      "test",
						Evaluation: 10,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGames := mock_service.NewMockGames(ctrl)

			tt.mockBehavior(mockGames, tt.responseData)

			services := service.Services{Games: mockGames}
			handler := Handler{services: &services}

			// Init Endpoint
			r := gin.New()
			r.GET("/game", handler.getAllGame)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/game", nil)

			// Write Request
			r.ServeHTTP(w, req)

			// Parse the actual response body into the expected struct
			var actualResponseBody getAllGameResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualResponseBody)
			if err != nil {
				t.Errorf("failed to unmarshal response body: %v", err)
				return
			}

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, actualResponseBody)
		})
	}
}

func TestHandler_gameGetById(t *testing.T) {
	type mockBehavior func(*mock_service.MockGames, int, domain.Game)

	tests := []struct {
		name         string
		id           int
		mockBehavior mockBehavior
		statusCode   int
		responseBody domain.Game
		responseErr  errorResponse
	}{
		{
			name: "ok",
			id:   1,
			mockBehavior: func(r *mock_service.MockGames, id int, game domain.Game) {
				r.EXPECT().GetById(id).Return(game, nil)
			},
			statusCode: 200,
			responseBody: domain.Game{
				Id:         1,
				Title:      "test",
				Genre:      "test",
				Evaluation: 10,
			},
		},
		{
			name:         "miss id param",
			mockBehavior: func(r *mock_service.MockGames, id int, game domain.Game) {},
			statusCode:   400,
			responseErr: errorResponse{
				Message: "invalid id param",
			},
		},
		{
			name:         "invalid id param -1",
			id:           -1,
			mockBehavior: func(r *mock_service.MockGames, id int, game domain.Game) {},
			statusCode:   400,
			responseErr: errorResponse{
				Message: "invalid id param",
			},
		},
		{
			name:         "invalid id param 0",
			id:           0,
			mockBehavior: func(r *mock_service.MockGames, id int, game domain.Game) {},
			statusCode:   400,
			responseErr: errorResponse{
				Message: "invalid id param",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGames := mock_service.NewMockGames(ctrl)

			tt.mockBehavior(mockGames, tt.id, tt.responseBody)

			services := service.Services{Games: mockGames}
			handler := Handler{services: &services}

			// Init Endpoint
			r := gin.New()
			r.GET("/game/:id", handler.getGameByID)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/game/%d", tt.id), nil)

			// Write Request
			r.ServeHTTP(w, req)

			// Parse the actual response body into the expected struct
			if tt.responseErr.Message == "" {
				var actualResponseBody domain.Game
				err := json.Unmarshal(w.Body.Bytes(), &actualResponseBody)
				if err != nil {
					t.Errorf("failed to unmarshal response body: %v", err)
				}

				// Assert
				assert.Equal(t, tt.statusCode, w.Code)
				assert.Equal(t, tt.responseBody, actualResponseBody)
			}

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, tt.responseBody)
		})
	}
}

func TestHandler_gameDeleteGameById(t *testing.T) {
	type mockBehavior func(*mock_service.MockGames, int)

	tests := []struct {
		name         string
		id           int
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name: "ok",
			id:   1,
			mockBehavior: func(r *mock_service.MockGames, id int) {
				r.EXPECT().Delete(id).Return(nil)
			},
			statusCode:   200,
			responseBody: `{"status":"ok"}`,
		},
		{
			name:         "miss id param",
			mockBehavior: func(r *mock_service.MockGames, id int) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:         "invalid id param 0",
			id:           0,
			mockBehavior: func(r *mock_service.MockGames, id int) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:         "invalid id param -1",
			id:           -1,
			mockBehavior: func(r *mock_service.MockGames, id int) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGames := mock_service.NewMockGames(ctrl)

			tt.mockBehavior(mockGames, tt.id)

			services := service.Services{Games: mockGames}
			handler := Handler{services: &services}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/game/:id", handler.deleteGameByID)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/game/%d", tt.id), nil)

			// Write Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
