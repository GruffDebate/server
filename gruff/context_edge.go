package gruff

import (
	"fmt"
)

// A ContextEdge is an edge that goes from a Context
// to one or more Claims that are made in that Context
type ContextEdge struct {
	Edge
}

// ArangoObject interface

func (c ContextEdge) CollectionName() string {
	return "context_edges"
}

func (c ContextEdge) ArangoKey() string {
	return c.Key
}

func (c ContextEdge) ArangoID() string {
	return fmt.Sprintf("%s/%s", c.CollectionName(), c.ArangoKey())
}

func (c ContextEdge) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (c *ContextEdge) Create(ctx *ServerContext) Error {
	return CreateArangoObject(ctx, c)
}

func (c *ContextEdge) Update(ctx *ServerContext, updates Updates) Error {
	return NewServerError("This item cannot be modified")
}

func (c *ContextEdge) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, c)
}

// Business methods

func FindContextEdge(ctx *ServerContext, contextArangoKey, claimArangoKey string) (ContextEdge, Error) {
	context := Context{}
	context.Key = contextArangoKey
	claim := Claim{}
	claim.Key = claimArangoKey

	edge := ContextEdge{}
	bindVars := BindVars{
		"claim":   claim.ArangoID(),
		"context": context.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                      FILTER obj._to == @claim
                                         AND obj._from == @context
                                         AND obj.end == null
                                       LIMIT 1
                                      RETURN obj`,
		ContextEdge{}.CollectionName(),
	)
	err := FindArangoObject(ctx, query, bindVars, &edge)
	return edge, err
}
