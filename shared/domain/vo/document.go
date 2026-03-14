package vo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloud-tech-develop/aura-back/shared/errors"
)

var dvRegex = regexp.MustCompile(`^[\w-]{5,20}$`)

// Document represents an identification document (e.g., NIT, DV).
type Document string

// ParseDocument parses and validates a document string.
func ParseDocument(s string) (Document, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return "", fmt.Errorf("%s", errors.ErrDocumentRequired)
	}

	if !dvRegex.MatchString(s) {
		return "", fmt.Errorf("%s", errors.ErrInvalidDocument)
	}

	return Document(s), nil
}

func (d Document) String() string {
	return string(d)
}

// UnmarshalJSON parses JSON bytes into a valid Document.
func (d *Document) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	doc, err := ParseDocument(s)
	if err != nil {
		return err
	}

	*d = doc
	return nil
}

// MarshalJSON returns the JSON encoding of Document.
func (d Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// Value implements the driver.Valuer interface.
func (d Document) Value() (driver.Value, error) {
	return string(d), nil
}

// Scan implements the sql.Scanner interface.
func (d *Document) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*d = Document(v)
	case []byte:
		*d = Document(string(v))
	default:
		return fmt.Errorf("cannot scan type %T into Document", value)
	}
	return nil
}
