package domain

import "encoding/json"

// PageParams holds pagination request parameters (legacy format)
type PageParams struct {
	First  int64  `json:"first"` // offset (starting from 1)
	Rows   int64  `json:"rows"`  // limit
	Search string `json:"search"`
}

// PageRequest holds the new pagination request parameters with generic params
type PageRequest struct {
	Page   int64          `json:"page"`
	Limit  int64          `json:"limit"`
	Search string         `json:"search"`
	Sort   string         `json:"sort"`
	Order  string         `json:"order"`
	Params map[string]any `json:"params"`
}

// ParsePageRequest parses a PageRequest from JSON, handling both new and legacy formats
func ParsePageRequest(data json.RawMessage) (*PageRequest, error) {
	if len(data) == 0 {
		return DefaultPageRequest(), nil
	}

	var req PageRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	// Apply defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Sort == "" {
		req.Sort = "id"
	}
	if req.Order == "" {
		req.Order = "asc"
	}

	return &req, nil
}

// DefaultPageRequest returns a default page request
func DefaultPageRequest() *PageRequest {
	return &PageRequest{
		Page:   1,
		Limit:  10,
		Search: "",
		Sort:   "id",
		Order:  "asc",
		Params: make(map[string]any),
	}
}

// PageResult holds a paginated result with items and pagination metadata
// The fields are flattened to match the expected JSON structure:
// {"items": [...], "total": 15, "page": 1, "limit": 10, "totalPages": 2}
type PageResult struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int64       `json:"page"`
	Limit      int64       `json:"limit"`
	TotalPages int64       `json:"totalPages"`
}
