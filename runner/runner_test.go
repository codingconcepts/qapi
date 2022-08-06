package runner

import (
	"net/http"
	"testing"

	"github.com/codingconcepts/qapi/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRunner_Start(t *testing.T) {
	cases := []struct {
		name      string
		variables map[string]string
		request   models.Request
		// TODO: Add assertion for post-request variables.
	}{
		{
			name: "request with headers",
			variables: map[string]string{
				"token": "some_token",
			},
			request: models.Request{
				Method: http.MethodGet,
				Path:   "/headers",
				Headers: map[string]string{
					"Authorization": "Bearer {{token}}",
				},
			},
		},
		{
			name: "request with url params",
			variables: map[string]string{
				"id": "some_id",
			},
			request: models.Request{
				Method: http.MethodGet,
				Path:   "/params/{{id}}",
			},
		},
		{
			name: "request with body",
			variables: map[string]string{
				"username": "un",
				"password": "pw",
			},
			request: models.Request{
				Method: http.MethodPost,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Path: "/body",
				Body: `{"username": "{{username}}", "password": "{{password}}"}`,
			},
		},
	}

	defer gock.Off()

	gock.New("http://localhost:8080/test").
		Get("/headers").
		MatchHeader("Authorization", "Bearer some_token").
		Reply(200)

	gock.New("http://localhost:8080/test").
		Get("/params/some_id").
		Reply(200)

	gock.New("http://localhost:8080/test").
		Post("/body").
		MatchType("json").
		JSON(map[string]string{"username": "un", "password": "pw"}).
		Reply(200)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runner := createRunner(c.variables, c.request)

			err := runner.Start()
			assert.NoError(t, err)
		})
	}
}

func createRunner(variables map[string]string, req models.Request) *Runner {
	return &Runner{
		Environment: models.Environment{
			BaseURL: "http://localhost:8080/test",
		},
		Variables: variables,
		Requests:  []models.Request{req},
	}
}