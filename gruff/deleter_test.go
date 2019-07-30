package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDeleter(t *testing.T) {
	assert.True(t, IsDeleter(reflect.TypeOf(&User{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&Claim{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&Argument{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&Inference{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&BaseClaimEdge{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&PremiseEdge{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&Context{})))
	assert.True(t, IsDeleter(reflect.TypeOf(&ContextEdge{})))
}
