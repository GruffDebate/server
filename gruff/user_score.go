package gruff

import (
	"fmt"
)

// A UserScore is an edge that goes from a User
// to either a Claim (truth score) or an Argument (relevance score)
// UserScore corresponds directly to the "Personal Score" or "Belief Score"
// described in the Canonical Debate white paper: https://github.com/canonical-debate-lab/paper#33322-belief-scores
// The Score attribute is a float value from 0 to 1.0, corresponding to a % belief in the truth/relevance of the target
// (where 0 = 0% and 1.0 = 100%)
type UserScore struct {
	Edge
	Score float32 `json:"score"`
}

// ArangoObject interface

func (u UserScore) CollectionName() string {
	return "scores"
}

func (u UserScore) ArangoKey() string {
	return u.Key
}

func (u UserScore) ArangoID() string {
	return fmt.Sprintf("%s/%s", u.CollectionName(), u.ArangoKey())
}

func (u UserScore) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (u *UserScore) Create(ctx *ServerContext) Error {
	return CreateArangoObject(ctx, u)
}

func (u *UserScore) Update(ctx *ServerContext, updates Updates) Error {
	return NewServerError("This item cannot be modified")
}

func (u *UserScore) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, u)
}
