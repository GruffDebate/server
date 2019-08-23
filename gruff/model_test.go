package gruff

import (
	"fmt"
	"testing"

	"github.com/GruffDebate/server/support"
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

func TestClearTransientFields(t *testing.T) {
	c := Claim{}
	err := ClearTransientFields(c)
	assert.Error(t, err)
	assert.Equal(t, "Cannot clear values on an immutable object", err.Error())

	var cPtr *Claim
	err = ClearTransientFields(cPtr)
	assert.Error(t, err)
	assert.Equal(t, "Cannot clear values on a nil item", err.Error())

	c.ID = fmt.Sprintf("SLDKFJDS:LFJDS:LJFD")
	c.Title = "Title of a thing"
	c.Negation = "Not title of a thing"
	c.Note = "Note of a thing"
	c.Truth = 1.0
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", c.ID)
	assert.Equal(t, "Title of a thing", c.Title)
	assert.Equal(t, "Not title of a thing", c.Negation)
	assert.Equal(t, float32(1.0), c.Truth)
	assert.Equal(t, []Claim(nil), c.PremiseClaims)
	assert.Equal(t, []Argument(nil), c.ProArgs)
	assert.Equal(t, []Argument(nil), c.ConArgs)
	assert.Equal(t, []Link(nil), c.Links)
	assert.Equal(t, []Context(nil), c.ContextElems)

	assert.NoError(t, ClearTransientFields(&c))
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", c.ID)
	assert.Equal(t, "Title of a thing", c.Title)
	assert.Equal(t, "Not title of a thing", c.Negation)
	assert.Equal(t, float32(1.0), c.Truth)
	assert.Equal(t, []Claim(nil), c.PremiseClaims)
	assert.Equal(t, []Argument(nil), c.ProArgs)
	assert.Equal(t, []Argument(nil), c.ConArgs)
	assert.Equal(t, []Link(nil), c.Links)
	assert.Equal(t, []Context(nil), c.ContextElems)

	c.PremiseClaims = []Claim{Claim{Title: "Just delete me"}}
	c.ProArgs = []Argument{Argument{Title: "Nothing else matters"}}
	c.ConArgs = []Argument{Argument{Title: "Let me go"}}
	c.Links = []Link{Link{Title: "A link is a pointer?"}}
	c.ContextElems = []Context{Context{ShortName: "What's the point?"}}
	assert.NoError(t, ClearTransientFields(&c))
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", c.ID)
	assert.Equal(t, "Title of a thing", c.Title)
	assert.Equal(t, "Not title of a thing", c.Negation)
	assert.Equal(t, float32(1.0), c.Truth)
	assert.Equal(t, []Claim(nil), c.PremiseClaims)
	assert.Equal(t, []Argument(nil), c.ProArgs)
	assert.Equal(t, []Argument(nil), c.ConArgs)
	assert.Equal(t, []Link(nil), c.Links)
	assert.Equal(t, []Context(nil), c.ContextElems)

	a := Argument{
		TargetClaimID: support.StringPtr(":LSDKFSDFJSLDKFJ"),
		ClaimID:       "LSDFJSD:LFJDSL:J",
		Title:         "More titles, like we're royalty or something",
		Pro:           true,
		Relevance:     0.32,
		Str:           0.11,
	}
	a.Key = ";alsdfas;ldkfjas;ldkfas;ldfkjsal;dfjksa"
	assert.NoError(t, ClearTransientFields(&a))
	assert.Equal(t, ":LSDKFSDFJSLDKFJ", *a.TargetClaimID)
	assert.Equal(t, ";alsdfas;ldkfjas;ldkfas;ldfkjsal;dfjksa", a.Key)
	assert.Equal(t, "More titles, like we're royalty or something", a.Title)
	assert.Equal(t, true, a.Pro)
	assert.Equal(t, float32(0.32), a.Relevance)
	assert.Equal(t, float32(0.11), a.Str)
	assert.Equal(t, (*Claim)(nil), a.TargetClaim)
	assert.Equal(t, (*Argument)(nil), a.TargetArgument)
	assert.Equal(t, (*Claim)(nil), a.Claim)
	assert.Equal(t, []Argument(nil), a.ProArgs)
	assert.Equal(t, []Argument(nil), a.ConArgs)

	a.TargetClaim = &Claim{Title: "Target"}
	a.TargetArgument = &Argument{Title: "Target"}
	a.Claim = &Claim{Title: "Not target"}
	a.ProArgs = []Argument{Argument{Title: "Nothing else matters"}}
	a.ConArgs = []Argument{Argument{Title: "You're just a sad copy of the Claim's sad argument"}}
	assert.NoError(t, ClearTransientFields(&a))
	assert.Equal(t, ":LSDKFSDFJSLDKFJ", *a.TargetClaimID)
	assert.Equal(t, ";alsdfas;ldkfjas;ldkfas;ldfkjsal;dfjksa", a.Key)
	assert.Equal(t, "More titles, like we're royalty or something", a.Title)
	assert.Equal(t, true, a.Pro)
	assert.Equal(t, float32(0.32), a.Relevance)
	assert.Equal(t, float32(0.11), a.Str)
	assert.Equal(t, (*Claim)(nil), a.TargetClaim)
	assert.Equal(t, (*Argument)(nil), a.TargetArgument)
	assert.Equal(t, (*Claim)(nil), a.Claim)
	assert.Equal(t, []Argument(nil), a.ProArgs)
	assert.Equal(t, []Argument(nil), a.ConArgs)
}

func TestClearTransientData(t *testing.T) {
	c := Claim{}
	m := Updates{}
	data, err := ClearTransientData(c, m)
	assert.NoError(t, err)

	var cPtr *Claim
	err = ClearTransientFields(cPtr)
	assert.Error(t, err)
	assert.Equal(t, "Cannot clear values on a nil item", err.Error())

	m["id"] = fmt.Sprintf("SLDKFJDS:LFJDS:LJFD")
	m["title"] = "Title of a thing"
	m["negation"] = "Not title of a thing"
	m["note"] = "Note of a thing"
	m["truth"] = 1.0

	data, err = ClearTransientData(c, m)
	assert.NoError(t, err)
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", m["id"])
	assert.Equal(t, "Title of a thing", m["title"])
	assert.Equal(t, "Not title of a thing", m["negation"])
	assert.Equal(t, 1.0, m["truth"])
	assert.Equal(t, nil, m["premises"])
	assert.Equal(t, nil, m["proargs"])
	assert.Equal(t, nil, m["conargs"])
	assert.Equal(t, nil, m["links"])
	assert.Equal(t, nil, m["contexts"])
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", data["id"])
	assert.Equal(t, "Title of a thing", data["title"])
	assert.Equal(t, "Not title of a thing", data["negation"])
	assert.Equal(t, 1.0, data["truth"])
	assert.Equal(t, nil, data["premises"])
	assert.Equal(t, nil, data["proargs"])
	assert.Equal(t, nil, data["conargs"])
	assert.Equal(t, nil, data["links"])
	assert.Equal(t, nil, data["contexts"])

	m["premises"] = []Claim{Claim{Title: "Just delete me"}}
	m["proargs"] = []Argument{Argument{Title: "Nothing else matters"}}
	m["conargs"] = []Argument{Argument{Title: "Let me go"}}
	m["links"] = []Link{Link{Title: "A link is a pointer?"}}
	m["contexts"] = []Context{Context{ShortName: "What's the point?"}}
	data, err = ClearTransientData(c, m)
	assert.NoError(t, err)
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", m["id"])
	assert.Equal(t, "Title of a thing", m["title"])
	assert.Equal(t, "Not title of a thing", m["negation"])
	assert.Equal(t, 1.0, m["truth"])
	assert.Equal(t, []Claim{Claim{Title: "Just delete me"}}, m["premises"])
	assert.Equal(t, []Argument{Argument{Title: "Nothing else matters"}}, m["proargs"])
	assert.Equal(t, []Argument{Argument{Title: "Let me go"}}, m["conargs"])
	assert.Equal(t, []Link{Link{Title: "A link is a pointer?"}}, m["links"])
	assert.Equal(t, []Context{Context{ShortName: "What's the point?"}}, m["contexts"])
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", data["id"])
	assert.Equal(t, "Title of a thing", data["title"])
	assert.Equal(t, "Not title of a thing", data["negation"])
	assert.Equal(t, 1.0, data["truth"])
	assert.Equal(t, nil, data["premises"])
	assert.Equal(t, nil, data["proargs"])
	assert.Equal(t, nil, data["conargs"])
	assert.Equal(t, nil, data["links"])
	assert.Equal(t, nil, data["contexts"])

	a := Argument{
		TargetClaimID: support.StringPtr(":LSDKFSDFJSLDKFJ"),
		ClaimID:       "LSDFJSD:LFJDSL:J",
		Title:         "More titles, like we're royalty or something",
		Pro:           true,
		Relevance:     0.32,
		Str:           0.11,
	}
	m = data
	m["proargs"] = []Argument{Argument{Title: "Nothing else matters"}}
	m["conargs"] = []Argument{Argument{Title: "Let me go"}}
	m["claim"] = Claim{Title: "Just delete me"}
	m["targetClaim"] = Claim{Title: "Just delete to me"}
	m["targetArg"] = Argument{Title: "Let me to go"}
	data, err = ClearTransientData(a, m)
	assert.NoError(t, err)
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", m["id"])
	assert.Equal(t, "Title of a thing", m["title"])
	assert.Equal(t, "Not title of a thing", m["negation"])
	assert.Equal(t, 1.0, m["truth"])
	assert.Equal(t, Claim{Title: "Just delete me"}, m["claim"])
	assert.Equal(t, Claim{Title: "Just delete to me"}, m["targetClaim"])
	assert.Equal(t, Argument{Title: "Let me to go"}, m["targetArg"])
	assert.Equal(t, []Argument{Argument{Title: "Nothing else matters"}}, m["proargs"])
	assert.Equal(t, []Argument{Argument{Title: "Let me go"}}, m["conargs"])
	assert.Equal(t, "SLDKFJDS:LFJDS:LJFD", data["id"])
	assert.Equal(t, "Title of a thing", data["title"])
	assert.Equal(t, "Not title of a thing", data["negation"])
	assert.Equal(t, 1.0, data["truth"])
	assert.Equal(t, nil, data["claim"])
	assert.Equal(t, nil, data["targetClaim"])
	assert.Equal(t, nil, data["targetArg"])
	assert.Equal(t, nil, data["proargs"])
	assert.Equal(t, nil, data["conargs"])
}
