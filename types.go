package validator

// ValidationResult mirrors the CountriesDB API response.
type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// CountryOptions toggles follow_upward logic.
type CountryOptions struct {
	FollowUpward bool
}

// SubdivisionOptions toggles follow_related / allow_parent_selection logic.
type SubdivisionOptions struct {
	FollowRelated        bool
	AllowParentSelection bool
}

type multiResult struct {
	Results []ValidationResult `json:"results"`
}

type apiError struct {
	Message string `json:"message"`
}


