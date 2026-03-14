package vo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	"github.com/cloud-tech-develop/aura-back/shared/errors"
)

// Email represents a validated email address.
type Email string

// ParseEmail parses and validates an email string.
func ParseEmail(s string) (Email, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return "", fmt.Errorf("%s", errors.ErrEmailRequired)
	}

	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", fmt.Errorf("%s", errors.ErrInvalidEmail)
	}

	return Email(addr.Address), nil
}

func (e Email) String() string {
	return string(e)
}

// UnmarshalJSON parses JSON bytes into a valid Email.
func (e *Email) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	email, err := ParseEmail(s)
	if err != nil {
		return err
	}

	*e = email
	return nil
}

// MarshalJSON returns the JSON encoding of Email.
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// Value implements the driver.Valuer interface.
func (e Email) Value() (driver.Value, error) {
	return string(e), nil
}

// Scan implements the sql.Scanner interface.
func (e *Email) Scan(value interface{}) error {
	if value == nil {
		*e = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*e = Email(v)
	case []byte:
		*e = Email(string(v))
	default:
		return fmt.Errorf("cannot scan type %T into Email", value)
	}
	return nil
}
