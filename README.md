# github.com/countriesdb/validator-go

**Backend validation package for CountriesDB.** Provides server-side validation for country and subdivision codes using ISO 3166-1 and ISO 3166-2 standards.

[![Go Reference](https://pkg.go.dev/badge/github.com/countriesdb/validator-go.svg)](https://pkg.go.dev/github.com/countriesdb/validator-go)

📖 **[Full Documentation](https://countriesdb.com/docs/backend-api)** | 🌐 **[Website](https://countriesdb.com)** | 📦 **[Package](https://github.com/countriesdb/validator-go)**

**Important**: This package only provides validation methods. Data fetching for frontend widgets must be done through frontend packages ([`@countriesdb/widget-core`](https://www.npmjs.com/package/@countriesdb/widget-core), [`@countriesdb/widget`](https://www.npmjs.com/package/@countriesdb/widget)).

## Getting Started

**⚠️ API Key Required:** This package requires a CountriesDB **private** API key to function. You must create an account at [countriesdb.com](https://countriesdb.com) to obtain your private API key. Test accounts are available with limited functionality.

- 🔑 [Get your API key](https://countriesdb.com) - Create an account and get your private key
- 📚 [View documentation](https://countriesdb.com/docs/backend-api) - Complete API reference and examples
- 💬 [Support](https://countriesdb.com) - Get help and support

## Features

- ✅ **ISO 3166 Compliant** - Validates ISO 3166-1 (countries) and ISO 3166-2 (subdivisions) codes
- ✅ **Multiple Validation Options** - Support for `follow_upward`, `follow_related`, and `allow_parent_selection`
- ✅ **Batch Validation** - Validate multiple countries or subdivisions in a single request
- ✅ **Context Support** - Full support for Go's `context.Context` for cancellation and timeouts
- ✅ **Detailed Error Messages** - Returns specific error messages from the CountriesDB API

## Installation

```bash
go get github.com/countriesdb/validator-go
```

**Package:** [github.com/countriesdb/validator-go](https://github.com/countriesdb/validator-go)

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

## Error Handling

### Single-Value Methods

Single-value methods (`ValidateCountry`, `ValidateSubdivision`) return validation results with `Valid: false` for invalid input. They only return errors on network failures:

```go
result, err := validator.ValidateCountry(ctx, "US", validator.CountryOptions{})
if err != nil {
    log.Fatal(err) // Network error
}
if !result.Valid {
    fmt.Printf("Validation failed: %s\n", result.Message)
}
```

### Multi-Value Methods

Multi-value methods (`ValidateCountries`, `ValidateSubdivisions`) return per-item results. Invalid codes are included in the results slice with `Valid: false`. They only return errors on network failures or invalid input types:

```go
results, err := validator.ValidateCountries(ctx, []string{"US", "BAD", "CA"}, validator.CountryOptions{})
if err != nil {
    log.Fatal(err) // Network error
}

for _, result := range results {
    if !result.Valid {
        fmt.Printf("Code %s failed: %s\n", result.Code, result.Message)
    }
}
```

**Note:** 
- Empty slices return empty results slices (not an error)
- Basic type checks are performed client-side (e.g., ensuring country is a non-empty string)
- Format validation (e.g., 2-character country codes) is handled by the backend and included in results with appropriate error messages
- Invalid format codes or invalid country codes are returned in the results slice with `Valid: false` rather than returning errors

## Examples

Runnable examples using this package are available in the [countriesdb/examples](https://github.com/countriesdb/examples) repository:

- [`go/backend-validator`](https://github.com/countriesdb/examples/tree/main/go/backend-validator) – Examples using the `validator-go` package
- [`go/backend-http`](https://github.com/countriesdb/examples/tree/main/go/backend-http) – HTTP client examples using Go's standard `net/http` package (raw HTTP calls)

## Requirements

- Go 1.22+
- Valid CountriesDB API key

## License

Proprietary (NAYEE LLC)

Copyright (c) NAYEE LLC. All rights reserved.

This software is the proprietary property of NAYEE LLC. For licensing inquiries, please contact [NAYEE LLC](https://nayee.net).


