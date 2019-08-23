package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLoader(t *testing.T) {
	assert.True(t, IsLoader(reflect.TypeOf(&User{})))
	assert.True(t, IsLoader(reflect.TypeOf(&Claim{})))
	assert.True(t, IsLoader(reflect.TypeOf(&Argument{})))
	assert.True(t, IsLoader(reflect.TypeOf(&Context{})))
	assert.False(t, IsLoader(reflect.TypeOf(&Inference{})))
	assert.False(t, IsLoader(reflect.TypeOf(&BaseClaimEdge{})))
	assert.False(t, IsLoader(reflect.TypeOf(&PremiseEdge{})))
	assert.False(t, IsLoader(reflect.TypeOf(&ContextEdge{})))
	assert.False(t, IsLoader(reflect.TypeOf(&UserScore{})))
}
