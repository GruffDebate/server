package gruff

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var ZERO_UUID uuid.UUID

func TestIsIdentifier(t *testing.T) {
	assert.False(t, IsIdentifier(reflect.TypeOf(Tag{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(Value{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(User{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(Context{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ClaimOpinion{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ArgumentOpinion{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Claim{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Argument{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Link{})))
}

func TestIdentifierGenerateUUID(t *testing.T) {
	d := Claim{}

	assert.Equal(t, ZERO_UUID, d.ID)

	assert.True(t, (&d).GenerateUUID() != ZERO_UUID)
	assert.True(t, d.ID != ZERO_UUID)
}

func TestSetCreatedBy(t *testing.T) {
	d := Claim{}
	SetCreatedByID(&d, 44)

	assert.Equal(t, uint64(44), d.CreatedByID)
}
