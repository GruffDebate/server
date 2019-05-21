package gruff

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsIdentifier(t *testing.T) {
	assert.False(t, IsIdentifier(reflect.TypeOf(Tag{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(User{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(Context{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ClaimOpinion{})))
	assert.False(t, IsIdentifier(reflect.TypeOf(ArgumentOpinion{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Claim{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Argument{})))
	assert.True(t, IsIdentifier(reflect.TypeOf(Link{})))
}
