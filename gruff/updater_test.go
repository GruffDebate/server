package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUpdater(t *testing.T) {
	assert.True(t, IsUpdater(reflect.TypeOf(&User{})))
	assert.False(t, IsUpdater(reflect.TypeOf(&Claim{})))
	assert.False(t, IsUpdater(reflect.TypeOf(&Argument{})))
	assert.False(t, IsUpdater(reflect.TypeOf(&Inference{})))
	assert.False(t, IsUpdater(reflect.TypeOf(&BaseClaimEdge{})))
	assert.False(t, IsUpdater(reflect.TypeOf(&PremiseEdge{})))
}