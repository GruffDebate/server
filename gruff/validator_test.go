package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidator(t *testing.T) {
	assert.True(t, IsValidator(reflect.TypeOf(User{})))
	assert.True(t, IsValidator(reflect.TypeOf(Claim{})))
	assert.True(t, IsValidator(reflect.TypeOf(Argument{})))
	assert.True(t, IsValidator(reflect.TypeOf(Link{})))
	assert.True(t, IsValidator(reflect.TypeOf(Context{})))
	assert.True(t, IsValidator(reflect.TypeOf(Value{})))
	assert.True(t, IsValidator(reflect.TypeOf(Tag{})))
}
