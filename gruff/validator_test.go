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
	assert.False(t, IsValidator(reflect.TypeOf(&Inference{})))
	assert.False(t, IsValidator(reflect.TypeOf(&BaseClaimEdge{})))
	assert.False(t, IsValidator(reflect.TypeOf(&PremiseEdge{})))
	assert.False(t, IsValidator(reflect.TypeOf(&ContextEdge{})))
}
