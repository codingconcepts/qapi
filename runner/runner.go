package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/codingconcepts/qapi/models"
	"github.com/codingconcepts/qapi/text"
	"github.com/tidwall/gjson"
)

// Runner holds the runtime configuration for the application.
type Runner struct {
	*models.Config
	Client              *http.Client
	assertionValidators map[string]models.Validator
}

func New(c *models.Config) *Runner {
	r := Runner{
		Config: c,
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	if r.Variables == nil {
		r.Variables = map[string]string{}
	}

	r.assertionValidators = map[string]models.Validator{
		"is_not_null": models.ValidateIsNotNull,
		"is_uuid":     models.ValidateIsUUID,
	}

	return &r
}

// Start making requests.
func (r *Runner) Start() error {
	for _, req := range r.Requests {
		log.Printf("[request] %s", req.Name)
		if err := r.runRequest(req); err != nil {
			return fmt.Errorf("making request: %w", err)
		}
	}

	return nil
}

func (r *Runner) runRequest(req models.Request) error {
	p := text.AddVariables(r.Variables, req.Path)
	b := text.AddVariables(r.Variables, req.Body)
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

	resp, err := r.Client.Do(request)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}

	if err = r.extractVariables(req.Extractors, resp); err != nil {
		return fmt.Errorf("extracting variables: %w", err)
	}

	if err = r.assertVariables(req.Assertions); err != nil {
		return fmt.Errorf("asserting variables: %w", err)
	}

	return nil
}

func (r *Runner) extractVariables(extractors []models.Extractor, resp *http.Response) error {
	for _, e := range extractors {
		for k, v := range e.Selectors {
			log.Printf("\t[extract] %v from %s", k, v)
		}

		switch strings.ToLower(e.Type) {
		case "json":
			return r.extractVariablesBodyJSON(e, resp)
		}
	}

	return nil
}

func (r *Runner) extractVariablesBodyJSON(extractor models.Extractor, resp *http.Response) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	body := string(bodyBytes)
	for k, v := range extractor.Selectors {
		value := gjson.Get(body, v)
		if value.Exists() {
			r.Variables[k] = value.String()
		}
	}

	return nil
}

func (r *Runner) assertVariables(assertions []models.Assertion) error {
	for _, a := range assertions {
		log.Printf("\t[assert] %s %s", a.Variable, a.Type)

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
