package validator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.countriesdb.com"

// Validator validates country and subdivision codes via the CountriesDB backend API.
type Validator struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option customizes the Validator.
type Option func(*Validator)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) Option {
	return func(v *Validator) {
		if baseURL != "" {
			v.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// WithHTTPClient provides a custom http.Client (otherwise a sane default is used).
func WithHTTPClient(h *http.Client) Option {
	return func(v *Validator) {
		if h != nil {
			v.httpClient = h
		}
	}
}

// NewValidator creates a CountriesDB validator.
func NewValidator(apiKey string, opts ...Option) (*Validator, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("countriesdb: api key is required")
	}

	validator := &Validator{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(validator)
	}

	return validator, nil
}

// ValidateCountry validates a single country code.
func (v *Validator) ValidateCountry(ctx context.Context, code string, opts CountryOptions) (ValidationResult, error) {
	if len(code) != 2 {
		return ValidationResult{Valid: false, Message: "Invalid country code."}, nil
	}

	var result ValidationResult
	err := v.post(ctx, "/api/validate/country", map[string]any{
		"code":          strings.ToUpper(code),
		"follow_upward": opts.FollowUpward,
	}, &result)

	return result, err
}

// ValidateCountries validates multiple country codes.
func (v *Validator) ValidateCountries(ctx context.Context, codes []string, opts CountryOptions) ([]ValidationResult, error) {
	if len(codes) == 0 {
		return []ValidationResult{}, nil
	}

	// Validate format
	for i, code := range codes {
		if len(code) != 2 {
			return nil, fmt.Errorf("invalid country code format. All codes must be 2-character strings")
		}
		codes[i] = strings.ToUpper(code)
	}

	var response multiResult
	err := v.post(ctx, "/api/validate/country", map[string]any{
		"code":          codes,
		"follow_upward": false, // Disabled for multi-select
	}, &response)

	return response.Results, err
}

// ValidateSubdivision validates a single subdivision for a given country.
func (v *Validator) ValidateSubdivision(ctx context.Context, code string, country string, opts SubdivisionOptions) (ValidationResult, error) {
	if len(country) != 2 {
		return ValidationResult{Valid: false, Message: "Invalid country code."}, nil
	}

	var result ValidationResult
	err := v.post(ctx, "/api/validate/subdivision", map[string]any{
		"code":                   code,
		"country":                strings.ToUpper(country),
		"follow_related":         opts.FollowRelated,
		"allow_parent_selection": opts.AllowParentSelection,
	}, &result)

	return result, err
}

// ValidateSubdivisions validates multiple subdivisions for the same country.
func (v *Validator) ValidateSubdivisions(ctx context.Context, codes []string, country string, opts SubdivisionOptions) ([]ValidationResult, error) {
	if len(country) != 2 {
		return nil, errors.New("invalid country code")
	}

	if len(codes) == 0 {
		return []ValidationResult{}, nil
	}

	payloadCodes := make([]string, len(codes))
	for i, code := range codes {
		if code == "" {
			payloadCodes[i] = ""
			continue
		}
		payloadCodes[i] = code
	}

	var response multiResult
	err := v.post(ctx, "/api/validate/subdivision", map[string]any{
		"code":                   payloadCodes,
		"country":                strings.ToUpper(country),
		"follow_related":         false, // Disabled for multi-select
		"allow_parent_selection": opts.AllowParentSelection,
	}, &response)

	return response.Results, err
}

func (v *Validator) post(ctx context.Context, path string, payload map[string]any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+v.apiKey)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr apiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil || apiErr.Message == "" {
			return fmt.Errorf("countriesdb: http %d", resp.StatusCode)
		}
		return errors.New(apiErr.Message)
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}


