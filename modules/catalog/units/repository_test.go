package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitEntity_Fields(t *testing.T) {
	u := Unit{
		ID:            1,
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	assert.Equal(t, int64(1), u.ID)
	assert.Equal(t, "Kilogramo", u.Name)
	assert.Equal(t, "kg", u.Abbreviation)
	assert.True(t, u.Active)
	assert.True(t, u.AllowDecimals)
	assert.Equal(t, int64(1), u.EnterpriseID)
}

func TestUnitEntity_DefaultValues(t *testing.T) {
	u := Unit{
		Name:          "Litro",
		Abbreviation:  "L",
		EnterpriseID:  1,
	}

	assert.Equal(t, "Litro", u.Name)
	assert.Equal(t, "L", u.Abbreviation)
}

func TestUnitList_Fields(t *testing.T) {
	ul := UnitList{
		Id:           1,
		Name:         "Unit List Test",
		Abbreviation: "UT",
	}

	assert.Equal(t, int64(1), ul.Id)
	assert.Equal(t, "Unit List Test", ul.Name)
	assert.Equal(t, "UT", ul.Abbreviation)
}