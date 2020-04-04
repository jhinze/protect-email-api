package main

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"hinze.dev/home/config"
	"hinze.dev/home/models"
	"hinze.dev/home/routes"
	"hinze.dev/home/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockGoogleRecaptcha struct {
	secret       string
	mockResponse *models.RecaptchaResponse
	err          *services.RecaptchaError
}

func (g *mockGoogleRecaptcha) SiteVerify(response string, remoteIp string) (resp *models.RecaptchaResponse, err error) {
	resp = g.mockResponse
	if g.err != nil {
		err = *g.err
	}
	return
}

func Test(t *testing.T) {

	config.Email = "foo@bar"

	t.Run("Recaptcha service success", func(t *testing.T) {
		services.Recaptcha = &services.GoogleRecaptcha{Secret: "foo"}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockResponse, _ := httpmock.NewJsonResponder(200, models.RecaptchaResponse{
			Success:     true,
			ChallengeTS: time.Now(),
			Hostname:    "localhost",
			ErrorCodes:  nil,
		})

		httpmock.RegisterResponder("POST", `=~^https://www\.google\.com/recaptcha/api/siteverify.*`, mockResponse)

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, fmt.Sprintf("{\"email\":\"%s\"}", config.Email), w.Body.String())

	})

	t.Run("Recaptcha service verify false", func(t *testing.T) {
		services.Recaptcha = &services.GoogleRecaptcha{Secret: "foo"}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockResponse, _ := httpmock.NewJsonResponder(200, models.RecaptchaResponse{
			Success:     false,
			ChallengeTS: time.Now(),
			Hostname:    "localhost",
			ErrorCodes:  nil,
		})

		httpmock.RegisterResponder("POST", `=~^https://www\.google\.com/recaptcha/api/siteverify.*`, mockResponse)

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("Recaptcha service decode error", func(t *testing.T) {
		services.Recaptcha = &services.GoogleRecaptcha{Secret: "foo"}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", `=~^https://www\.google\.com/recaptcha/api/siteverify.*`,
			httpmock.NewStringResponder(200, `I'm not JSON'`))

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("Recaptcha service unavailable", func(t *testing.T) {
		services.Recaptcha = &services.GoogleRecaptcha{Secret: "foo"}

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", `=~^https://www\.google\.com/recaptcha/api/siteverify.*`,
			httpmock.NewStringResponder(500, "Some sort of internal server error"))

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("GET email success", func(t *testing.T) {
		services.Recaptcha = &mockGoogleRecaptcha{
			secret: "foo",
			mockResponse: &models.RecaptchaResponse{
				Success:     true,
				ChallengeTS: time.Now(),
				Hostname:    "foo.bar.com",
				ErrorCodes:  nil,
			},
		}

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, fmt.Sprintf("{\"email\":\"%s\"}", config.Email), w.Body.String())
	})

	t.Run("GET email token verify fail", func(t *testing.T) {
		services.Recaptcha = &mockGoogleRecaptcha{
			secret: "foo",
			mockResponse: &models.RecaptchaResponse{
				Success:     false,
				ChallengeTS: time.Now(),
				Hostname:    "foo.bar.com",
				ErrorCodes:  nil,
			},
		}

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("reCAPTCHA verify error", func(t *testing.T) {
		services.Recaptcha = &mockGoogleRecaptcha{
			secret: "foo",
			err: &services.RecaptchaError{
				Problem: "some problem",
			},
		}

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email?token=someValidToken", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("token missing from query", func(t *testing.T) {
		services.Recaptcha = &mockGoogleRecaptcha{
			secret: "foo",
			err: &services.RecaptchaError{
				Problem: "some problem",
			},
		}

		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/email", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

	t.Run("health check", func(t *testing.T) {
		router := routes.SetupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, `{"status":"good"}`, w.Body.String())
	})

}
