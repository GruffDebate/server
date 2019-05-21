package gruff

import (
	"testing"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArgumentValidateForCreate(t *testing.T) {
	a := Argument{}

	assert.Equal(t, "title: non zero value required;", a.ValidateForCreate().Error())

	a.Title = "A"
	assert.Equal(t, "title: A does not validate as length(3|1000);", a.ValidateForCreate().Error())

	a.Title = "This is a real argument"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.Description = "D"
	assert.Equal(t, "desc: D does not validate as length(3|4000);", a.ValidateForCreate().Error())

	a.Description = "This is a real description"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = ""
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = uuid.New().String()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.NoError(t, a.ValidateForCreate())

	a.TargetClaimID = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.Nil(t, a.ValidateForCreate())
}

func TestArgumentValidateForUpdate(t *testing.T) {
	a := Argument{}

	assert.Equal(t, "title: non zero value required;", a.ValidateForUpdate().Error())

	a.Title = "A"
	assert.Equal(t, "title: A does not validate as length(3|1000);", a.ValidateForUpdate().Error())

	a.Title = "This is a real argument"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate().Error())

	a.Description = "D"
	assert.Equal(t, "desc: D does not validate as length(3|4000);", a.ValidateForUpdate().Error())

	a.Description = "This is a real description"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate().Error())

	a.ClaimID = ""
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate().Error())

	a.ClaimID = uuid.New().String()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.NoError(t, a.ValidateForUpdate())

	a.TargetClaimID = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.Nil(t, a.ValidateForUpdate())
}
