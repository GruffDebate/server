package gruff

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsIdentifier(t *testing.T) {
	assert.False(t, IsIdentifier(reflect.TypeOf(User{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(Context{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ClaimOpinion{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ArgumentOpinion{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Claim{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Argument{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Link{})))
}

func TestGetIdentifier(t *testing.T) {
	u := User{}
	c := Claim{}
	c.ID = "id"
	c.Key = "key"

	_, err := GetIdentifier(u)
	assert.Equal(t, "Item is not an Identifier", err.Error())

	id, err := GetIdentifier(c)
	assert.NoError(t, err)
	assert.Equal(t, "id", id.ID)
	assert.Equal(t, "key", id.Key)
}
