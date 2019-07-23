package gruff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetKey(t *testing.T) {
	c := Claim{}
	err := SetKey(&c, "key")
	assert.NoError(t, err)
	assert.Equal(t, "key", c.Key)

	a := Argument{}
	err = SetKey(&a, "key")
	assert.NoError(t, err)
	assert.Equal(t, "key", a.Key)

	u := User{}
	err = SetKey(&u, "key")
	assert.NoError(t, err)
	assert.Equal(t, "key", u.Key)

	ctx := ServerContext{}
	err = SetKey(&ctx, "key")
	assert.Error(t, err)
	assert.Equal(t, "Item does not have a Key field", err.Error())

	var nilGuy *Claim
	err = SetKey(nilGuy, "key")
	assert.Error(t, err)
	assert.Equal(t, "Cannot set value on nil item", err.Error())
}
