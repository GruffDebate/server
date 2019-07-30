package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRestrictor(t *testing.T) {
	assert.True(t, IsRestrictor(reflect.TypeOf(&User{})))
	assert.True(t, IsRestrictor(reflect.TypeOf(&Claim{})))
	assert.True(t, IsRestrictor(reflect.TypeOf(&Argument{})))
	assert.False(t, IsRestrictor(reflect.TypeOf(&Inference{})))
	assert.False(t, IsRestrictor(reflect.TypeOf(&BaseClaimEdge{})))
	assert.False(t, IsRestrictor(reflect.TypeOf(&PremiseEdge{})))
	assert.True(t, IsRestrictor(reflect.TypeOf(&Context{})))
	assert.False(t, IsRestrictor(reflect.TypeOf(&ContextEdge{})))
}
