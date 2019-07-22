package gruff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	setupDB()
	defer teardownDB()

	user := User{
		Name:     "Imma User",
		Username: "ImmaUser",
		Email:    "immauser@gruff.org",
		Password: "monkey",
		Image:    "https://i.ytimg.com/vi/hYuViV9NgzA/hqdefault.jpg",
		Curator:  true,
		Admin:    true,
		URL:      "https://thetruth2020.org/",
	}

	saved := User{}
	saved.Username = user.Username
	err := saved.Load(CTX)
	assert.Empty(t, saved.Key)

	err = user.Create(CTX)
	assert.NoError(t, err)
	saved = User{}
	saved.Username = user.Username
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, "", saved.Password)
	assert.NotEmpty(t, saved.HashedPassword)
	assert.Equal(t, user.Name, saved.Name)
	assert.Equal(t, user.Username, saved.Username)
	assert.Equal(t, user.Email, saved.Email)
	assert.Equal(t, user.Image, saved.Image)
	assert.True(t, saved.Curator)
	assert.True(t, saved.Admin)
	assert.Equal(t, user.URL, saved.URL)
}

// TODO: test update
// TODO: test change password
// TODO: test validations
