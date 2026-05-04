package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryInterface(t *testing.T) {
	var _ Repository = (*repository)(nil)
}

func TestServiceInterface(t *testing.T) {
	var _ Service = (*service)(nil)
}

func TestProductTypeConstants(t *testing.T) {
	expectedTypes := []string{"STANDARD", "WEIGHTABLE", "KIT", "SERVICE"}

	assert.Equal(t, len(expectedTypes), len(ValidProductTypes))

	for i, expected := range expectedTypes {
		assert.Equal(t, expected, ValidProductTypes[i])
	}
}

func TestIsValidProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		want       bool
	}{
		{"valid STANDARD", "STANDARD", true},
		{"valid WEIGHTABLE", "WEIGHTABLE", true},
		{"valid KIT", "KIT", true},
		{"valid SERVICE", "SERVICE", true},
		{"invalid lowercase", "standard", false},
		{"invalid empty", "", false},
		{"invalid INVALIDO", "INVALIDO", false},
		{"invalid WEIGHABLE (wrong)", "WEIGHABLE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidProductType(tt.productType)
			assert.Equal(t, tt.want, got, "IsValidProductType(%q) = %v, want %v", tt.productType, got, tt.want)
		})
	}
}