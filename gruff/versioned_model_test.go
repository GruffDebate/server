package gruff

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsVersionedModel(t *testing.T) {
	assert.False(t, IsVersionedModel(reflect.TypeOf(User{})))
	assert.False(t, IsVersionedModel(reflect.TypeOf(Context{})))
	assert.False(t, IsVersionedModel(reflect.TypeOf(ClaimOpinion{})))
	assert.False(t, IsVersionedModel(reflect.TypeOf(ArgumentOpinion{})))
	assert.True(t, IsVersionedModel(reflect.TypeOf(Claim{})))
	assert.True(t, IsVersionedModel(reflect.TypeOf(Argument{})))
	assert.True(t, IsVersionedModel(reflect.TypeOf(Link{})))
}

func TestGetVersionedModel(t *testing.T) {
	u := User{}
	c := Claim{}
	c.ID = "id"
	c.Key = "key"

	_, err := GetVersionedModel(u)
	assert.Equal(t, "Item is not a VersionedModel", err.Error())

	id, err := GetVersionedModel(c)
	assert.NoError(t, err)
	assert.Equal(t, "id", id.ID)
	assert.Equal(t, "key", id.Key)
}

func TestIsVersioner(t *testing.T) {
	assert.False(t, IsVersioner(reflect.TypeOf(&User{})))
	assert.False(t, IsVersioner(reflect.TypeOf(&Context{})))
	assert.False(t, IsVersioner(reflect.TypeOf(&ClaimOpinion{})))
	assert.False(t, IsVersioner(reflect.TypeOf(&ArgumentOpinion{})))
	assert.True(t, IsVersioner(reflect.TypeOf(&Claim{})))
	assert.True(t, IsVersioner(reflect.TypeOf(&Argument{})))
	assert.False(t, IsVersioner(reflect.TypeOf(&Link{})))
}
