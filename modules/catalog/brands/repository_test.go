package brands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrandEntity_Fields(t *testing.T) {
	b := Brand{
		ID:           1,
		Name:         "Test Brand",
		Description: "Test Description",
		Active:       true,
		EnterpriseID: 1,
	}

	assert.Equal(t, int64(1), b.ID)
	assert.Equal(t, "Test Brand", b.Name)
	assert.Equal(t, "Test Description", b.Description)
	assert.True(t, b.Active)
	assert.Equal(t, int64(1), b.EnterpriseID)
}

func TestBrandEntity_DefaultActive(t *testing.T) {
	b := Brand{
		Name:         "Test Brand",
		EnterpriseID: 1,
	}

	assert.Equal(t, "Test Brand", b.Name)
	assert.Equal(t, int64(1), b.EnterpriseID)
}

func TestBrandList_Fields(t *testing.T) {
	bl := BrandList{
		ID:   1,
		Name: "Brand List Test",
	}

	assert.Equal(t, int64(1), bl.ID)
	assert.Equal(t, "Brand List Test", bl.Name)
}