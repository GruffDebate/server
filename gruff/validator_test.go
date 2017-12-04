package gruff

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsValidator(t *testing.T) {
	assert.True(t, IsValidator(reflect.TypeOf(User{})))
	assert.False(t, IsValidator(reflect.TypeOf(Claim{})))
	assert.True(t, IsValidator(reflect.TypeOf(Argument{})))
	assert.False(t, IsValidator(reflect.TypeOf(Link{})))
	assert.False(t, IsValidator(reflect.TypeOf(Context{})))
	assert.False(t, IsValidator(reflect.TypeOf(Value{})))
	assert.False(t, IsValidator(reflect.TypeOf(Tag{})))
}
