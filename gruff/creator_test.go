package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCreator(t *testing.T) {
	assert.True(t, IsCreator(reflect.TypeOf(&User{})))
	assert.True(t, IsCreator(reflect.TypeOf(&Claim{})))
	assert.True(t, IsCreator(reflect.TypeOf(&Argument{})))
	assert.True(t, IsCreator(reflect.TypeOf(&Inference{})))
	assert.True(t, IsCreator(reflect.TypeOf(&BaseClaimEdge{})))
	assert.True(t, IsCreator(reflect.TypeOf(&PremiseEdge{})))
	assert.True(t, IsCreator(reflect.TypeOf(&Context{})))
	assert.True(t, IsCreator(reflect.TypeOf(&ContextEdge{})))
}
