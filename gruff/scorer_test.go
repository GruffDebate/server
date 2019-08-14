package gruff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsScorer(t *testing.T) {
	assert.False(t, IsScorer(reflect.TypeOf(&User{})))
	assert.True(t, IsScorer(reflect.TypeOf(&Claim{})))
	assert.True(t, IsScorer(reflect.TypeOf(&Argument{})))
	assert.False(t, IsScorer(reflect.TypeOf(&Inference{})))
	assert.False(t, IsScorer(reflect.TypeOf(&BaseClaimEdge{})))
	assert.False(t, IsScorer(reflect.TypeOf(&PremiseEdge{})))
	assert.False(t, IsScorer(reflect.TypeOf(&Context{})))
	assert.False(t, IsScorer(reflect.TypeOf(&ContextEdge{})))
	assert.False(t, IsScorer(reflect.TypeOf(&UserScore{})))
}
