package gruff

import (
	"fmt"
)

// BaseClaim is an edge pointing from an Argument to the Claim on which it is based
// (the true/false part of the Argument)
type BaseClaimEdge struct {
	Edge
}

// ArangoObject interface

func (bc BaseClaimEdge) CollectionName() string {
	return "base_claims"
}

func (bc BaseClaimEdge) ArangoKey() string {
	return bc.Key
}

func (bc BaseClaimEdge) ArangoID() string {
	return fmt.Sprintf("%s/%s", bc.CollectionName(), bc.ArangoKey())
}

func (bc BaseClaimEdge) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (bc *BaseClaimEdge) Create(ctx *ServerContext) Error {
	return CreateArangoObject(ctx, bc)
}

func (bc *BaseClaimEdge) Update(ctx *ServerContext, updates map[string]interface{}) Error {
	return NewServerError("This item cannot be modified")
}

func (bc *BaseClaimEdge) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, bc)
}
