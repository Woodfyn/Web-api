package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Woodfyn/Web-api/internal/handler/rest"
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	h := rest.NewHandler(&service.Services{})

	require.IsType(t, &rest.Handler{}, h)
}

func TestNewHandler_InitRoutes(t *testing.T) {
	h := rest.NewHandler(&service.Services{})

	router := h.InitRoutes()

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
