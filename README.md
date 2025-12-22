# github.com/countriesdb/validator-go

Backend validation package for CountriesDB. Provides server-side validation for country and subdivision codes.

**Important**: This package only provides validation methods. Data fetching is frontend-only and must be done through frontend packages.

## Installation

```bash
go get github.com/countriesdb/validator-go
```

## Usage

### Standalone Validator

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/countriesdb/validator-go"
)

func main() {
	validator, err := validator.NewValidator("YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Validate a single country
	result, err := validator.ValidateCountry(ctx, "US", validator.CountryOptions{
		FollowUpward: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	if result.Valid {
		fmt.Println("Valid country")
	} else {
		fmt.Printf("Invalid: %s\n", result.Message)
	}

	// Validate a single subdivision
	subdivisionResult, err := validator.ValidateSubdivision(ctx, "US-CA", "US", validator.SubdivisionOptions{
		FollowRelated:        false,
		AllowParentSelection: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	if subdivisionResult.Valid {
		fmt.Println("Valid subdivision")
	}

	// Validate multiple countries
	codes := []string{"US", "CA", "MX"}
	results, err := validator.ValidateCountries(ctx, codes, validator.CountryOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range results {
		fmt.Printf("%s: %s\n", r.Code, map[bool]string{true: "Valid", false: "Invalid"}[r.Valid])
	}

	// Validate multiple subdivisions
	subdivisionCodes := []string{"US-CA", "US-NY", "US-TX"}
	subdivisionResults, err := validator.ValidateSubdivisions(ctx, subdivisionCodes, "US", validator.SubdivisionOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range subdivisionResults {
		fmt.Printf("%s: %s\n", r.Code, map[bool]string{true: "Valid", false: "Invalid"}[r.Valid])
	}
}
```

### Configuration

```go
validator, err := validator.NewValidator(
	os.Getenv("COUNTRIESDB_PRIVATE_KEY"),
	validator.WithBaseURL(os.Getenv("COUNTRIESDB_BASE_URL")),
	validator.WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
)
```

## API Reference

### `NewValidator(apiKey, opts ...Option)`

Creates a new CountriesDB validator.

**Parameters:**
- `apiKey` (required): Your CountriesDB API key
- `opts` (optional): Configuration options:
  - `WithBaseURL(baseURL)`: Override the default API base URL (defaults to `https://api.countriesdb.com`)
  - `WithHTTPClient(client)`: Provide a custom `http.Client` (defaults to 10s timeout)

**Returns:** `*Validator`, `error`

### `ValidateCountry(ctx, code, opts)`

Validate a single country code.

**Parameters:**
- `ctx`: Context for request cancellation/timeout
- `code`: ISO 3166-1 alpha-2 country code
- `opts`: `CountryOptions` with `FollowUpward` boolean

**Returns:** `ValidationResult`, `error`

### `ValidateCountries(ctx, codes, opts)`

Validate multiple country codes.

**Parameters:**
- `ctx`: Context for request cancellation/timeout
- `codes`: Slice of ISO 3166-1 alpha-2 country codes
- `opts`: `CountryOptions` (FollowUpward is always false for multi-select)

**Returns:** `[]ValidationResult`, `error`

### `ValidateSubdivision(ctx, code, country, opts)`

Validate a single subdivision code.

**Parameters:**
- `ctx`: Context for request cancellation/timeout
- `code`: Subdivision code (e.g., 'US-CA') or empty string
- `country`: ISO 3166-1 alpha-2 country code
- `opts`: `SubdivisionOptions` with `FollowRelated` and `AllowParentSelection` booleans

**Returns:** `ValidationResult`, `error`

### `ValidateSubdivisions(ctx, codes, country, opts)`

Validate multiple subdivision codes.

**Parameters:**
- `ctx`: Context for request cancellation/timeout
- `codes`: Slice of subdivision codes or empty strings
- `country`: ISO 3166-1 alpha-2 country code
- `opts`: `SubdivisionOptions` (FollowRelated is always false for multi-select)

**Returns:** `[]ValidationResult`, `error`

### `ValidationResult`

```go
type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
```

## Examples

See `/examples/backend-go-http` and `/examples/backend-go-resty` for runnable demos that mirror the documentation and use this package/API.

## Requirements

- Go 1.22+
- Valid CountriesDB API key

## License

MIT


