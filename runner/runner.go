package runner

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/codingconcepts/qapi/models"
	"github.com/codingconcepts/qapi/text"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

// Runner holds the runtime configuration for the application.
type Runner struct {
	models.Config
	client              *http.Client
	duration            time.Duration
	events              chan models.RequestResult
	assertionValidators map[string]models.Validator
	logger              *zerolog.Logger
}

func New(c models.Config, d time.Duration, e chan models.RequestResult, logger *zerolog.Logger) *Runner {
	r := Runner{
		Config: c,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		duration: d,
		events:   e,
		logger:   logger,
	}

	if r.Variables == nil {
		r.Variables = map[string]any{}
	}

	r.assertionValidators = map[string]models.Validator{
		"is_not_null":  models.ValidateIsNotNull,
		"is_not_empty": models.ValidateIsNotEmpty,
		"is_uuid":      models.ValidateIsUUID,
	}

	return &r
}

func (r *Runner) Start() error {
	// Run setup requests first.
	for _, req := range r.SetupRequests {
		r.logger.Debug().Str("name", req.Name).Msg("setup_request")

		if err := r.runRequest(req); err != nil {
			return fmt.Errorf("running request: %w", err)
		}
	}

	// Then go onto regular requests (ignoring errors this time).
	for range time.Tick(r.RequestFrequency) {
		for _, req := range r.Requests {
			r.logger.Debug().Str("name", req.Name).Msg("request")

			if err := r.runRequest(req); err != nil {
				r.logger.Warn().Str("request", req.Name).Err(err).Msg("error")
			}
		}
	}

	return nil
}

func (r *Runner) runRequest(req models.Request) (err error) {
	defer func() {
		if err != nil {
			// Publish an error response to the logger.
			r.events <- models.RequestResult{
				StatusCode: 999,
			}

		}
	}()

	p := req.Path
	b := req.Body

	// Substitute variables.
	p = text.AddVariables(r.Variables, p)
	b = text.AddVariables(r.Variables, b)

	// Substitute generators.
	p = text.GenerateVariable(p)
	b = text.GenerateVariable(b)

	u, err := url.JoinPath(r.Environment.BaseURL, p)
	if err != nil {
		return fmt.Errorf("forming request path: %w", err)
	}

	request, err := http.NewRequest(req.Method, u, strings.NewReader(b))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	headers := map[string][]string{}
	for k, v := range req.Headers {
		headers[k] = []string{text.AddVariables(r.Variables, v)}
	}
	request.Header = headers

	r.logger.Debug().Str("path", p).Str("body", b).Any("headers", headers).Msg("request")

	resp, err := r.client.Do(request)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}

	body, err := readResponse(resp)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	// Publish the response to the logger.
	r.events <- models.RequestResult{
		StatusCode: resp.StatusCode,
	}

	// Return early in the result of a failure but don't report as a failure,
	// as it might be transient.
	if resp.StatusCode >= http.StatusBadRequest {
		r.logger.Warn().Int("status_code", resp.StatusCode).Str("body", body).Msg("response")
		return nil
	}

	if err = r.extractVariables(req.Extractors, body); err != nil {
		return fmt.Errorf("extracting variables: %w", err)
	}

	if err = r.assertVariables(req.Assertions); err != nil {
		return fmt.Errorf("asserting variables: %w", err)
	}

	return nil
}

func readResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	return strings.TrimSuffix(string(bodyBytes), "\n"), nil
}

func (r *Runner) extractVariables(extractors []models.Extractor, body string) error {
	for _, e := range extractors {
		switch strings.ToLower(e.Type) {
		case "json":
			return r.extractVariablesBodyJSON(e, body)
		}
	}

	return nil
}

func (r *Runner) extractVariablesBodyJSON(extractor models.Extractor, body string) error {
	r.logger.Debug().Str("body", body).Msg("response")

	for k, v := range extractor.Selectors {
		value := gjson.Get(body, v)
		r.logger.Debug().Str("key", k).Str("body value", value.Raw).Msg("json extract")

		if !value.Exists() {
			return nil
		}

		// Support for array responses.
		if value.IsArray() {
			r.Variables[k] = lo.Map(value.Array(), func(r gjson.Result, _ int) any {
				return parseScalar(r)
			})

			return nil
		}

		// Value is a scalar.
		r.Variables[k] = parseScalar(value)
	}

	return nil
}

func parseScalar(r gjson.Result) any {
	switch r.Type {
	case gjson.Number:
		// Special handling for floats and integers.
		if strings.Contains(r.Raw, ".") {
			return r.Num
		} else {
			return r.Int()
		}
	case gjson.String:
		return r.Str
	default:
		return r.Value()
	}
}

func (r *Runner) assertVariables(assertions []models.Assertion) error {
	for _, a := range assertions {
		validator, ok := r.assertionValidators[a.Type]
		if !ok {
			return fmt.Errorf("missing validator for assertion: %q", a.Type)
		}

		variable := strings.Trim(a.Variable, "{}")
		if err := validator(variable, r.Variables); err != nil {
			return fmt.Errorf("assertion failed: %w", err)
		}
	}

	return nil
}
