package domain

// PageParams holds pagination request parameters
type PageParams struct {
	First  int64  `json:"first"` // offset (starting from 1)
	Rows   int64  `json:"rows"`  // limit
	Search string `json:"search"`
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
