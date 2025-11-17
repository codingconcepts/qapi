package runner

import (
	"net/http"
	"testing"
	"time"

	"github.com/codingconcepts/qapi/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRunner_Start(t *testing.T) {
	cases := []struct {
		name              string
		variables         map[string]any
		request           models.Request
		expectedVariables map[string]string
	}{
		{
			name: "request with headers",
			variables: map[string]any{
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
			variables: map[string]any{
				"id": "some_id",
			},
			request: models.Request{
				Method: http.MethodGet,
				Path:   "/params/{{id}}",
			},
		},
		{
			name: "request with body",
			variables: map[string]any{
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
		{
			name: "request with json body extraction",
			request: models.Request{
				Method: http.MethodGet,
				Path:   "/extract",
				Extractors: []models.Extractor{
					{
						Type: "json",
						Selectors: map[string]string{
							"value": "a.b.c",
						},
					},
				},
			},
			expectedVariables: map[string]string{
				"value": "hello",
			},
		},
		{
			name: "request with json array body extraction",
			request: models.Request{
				Method: http.MethodGet,
				Path:   "/extract",
				Extractors: []models.Extractor{
					{
						Type: "json",
						Selectors: map[string]string{
							"value": "a.b.c",
						},
					},
				},
			},
			expectedVariables: map[string]string{
				"value": "hello",
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
		Get("/extract").
		Reply(200).
		BodyString(`{"a": {"b": {"c": "hello"}}}`)

	gock.New("http://localhost:8080/test").
		Post("/body").
		MatchType("json").
		JSON(map[string]string{"username": "un", "password": "pw"}).
		Reply(200)

	events := make(chan models.RequestResult, 1)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runner := New(models.Config{
				Environment: models.Environment{
					BaseURL: "http://localhost:8080/test",
				},
				Variables: c.variables,
				Requests: []models.Request{
					c.request,
				},
			}, time.Second, events, &zerolog.Logger{})

			err := runner.Start()
			assert.NoError(t, err)

			if c.expectedVariables != nil {
				for k, v := range c.expectedVariables {
					assert.Equal(t, v, runner.Variables[k])
				}
			}

			assert.Equal(t, 200, <-events)
		})
	}
}
