package gruff

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContext(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Key = "testuser"
	CTX.UserContext = u

	context := Context{
		ShortName:   "test",
		Title:       "A test Context",
		Description: "This Context was created just for testing",
		URL:         "https://canonicaldebate.com",
	}

	saved := Context{}
	saved.Key = context.Key
	err := saved.Load(CTX)
	assert.Error(t, err)
	assert.Empty(t, saved.Key)

	err = context.Create(CTX)
	assert.NoError(t, err)
	saved = Context{}
	saved.Key = context.Key
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, context.ShortName, saved.ShortName)
	assert.Equal(t, context.Title, saved.Title)
	assert.Equal(t, context.Description, saved.Description)
	assert.Equal(t, context.URL, saved.URL)

	// TODO
	/*
		err = context.Create(CTX)
		assert.Error(t, err)
		assert.Equal(t, "A context with the same Short Name already exists", err.Error())
	*/

	context = Context{}
	err = context.Create(CTX)
	assert.Error(t, err)
	assert.Equal(t, "name: non zero value required;title: non zero value required;url: non zero value required", err.Error())

	context.Title = "Something more than 3 characters"
	context.Description = "AB"
	err = context.Create(CTX)
	assert.Error(t, err)
	assert.Equal(t, "name: non zero value required;desc: AB does not validate as length(3|4000);url: non zero value required", err.Error())

	context.ShortName = "Dwarf"
	context.URL = "https://lotr.com"
	context.Description = ""
	err = context.Create(CTX)
	assert.NoError(t, err)
}

func TestSearchContexts(t *testing.T) {
	setupDB()
	defer teardownDB()

	ctx1 := Context{
		ShortName:   "test",
		Title:       "A Context whose title doesn't match the short name",
		Description: "This Context was created just for testing",
		URL:         "https://canonicaldebate.com",
	}
	ctx2 := Context{
		ShortName:   "dwarf",
		Title:       "Dwarf",
		Description: "A humanoid creature that loves beards, beer and battle",
		URL:         "https://lotr.com/Dwarf",
	}
	ctx3 := Context{
		ShortName:   "hobbit",
		Title:       "Hobbit",
		Description: "A humanoid creature that loves food, friends and food",
		URL:         "https://lotr.com/Hobbit",
	}
	ctx4 := Context{
		ShortName:   "elf",
		Title:       "Elf",
		Description: "A humanoid creature that is much too tall",
		URL:         "https://lotr.com/Elf",
	}
	err := ctx1.Create(CTX)
	assert.NoError(t, err)
	err = ctx2.Create(CTX)
	assert.NoError(t, err)
	err = ctx3.Create(CTX)
	assert.NoError(t, err)
	err = ctx4.Create(CTX)
	assert.NoError(t, err)

	ctxs, err := SearchContexts(CTX, "")
	assert.Equal(t, "Query term required", err.Error())

	ctxs, err = SearchContexts(CTX, "%")
	assert.Equal(t, 0, len(ctxs))

	ctxs, err = SearchContexts(CTX, "f")
	assert.Equal(t, 2, len(ctxs))
	sort.Slice(ctxs, func(i, j int) bool { return ctxs[i].Title < ctxs[j].Title })
	assert.Equal(t, "Dwarf", ctxs[0].Title)
	assert.Equal(t, "Elf", ctxs[1].Title)

	ctxs, err = SearchContexts(CTX, "elf")
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, "Elf", ctxs[0].Title)

	ctxs, err = SearchContexts(CTX, "elF")
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, "Elf", ctxs[0].Title)

	ctxs, err = SearchContexts(CTX, "test")
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, "A Context whose title doesn't match the short name", ctxs[0].Title)
}
