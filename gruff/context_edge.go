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
	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err
	}

	c.PrepareForCreate(ctx)

	_, aerr := col.CreateDocument(ctx.Context, c)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}

func (c *ContextEdge) Update(ctx *ServerContext, updates map[string]interface{}) Error {
	return NewServerError("This item cannot be modified")
}

func (c *ContextEdge) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, c)
}

// Business methods

func FindContextEdge(ctx *ServerContext, contextArangoKey, claimArangoKey string) (ContextEdge, Error) {
	db := ctx.Arango.DB

	context := Context{}
	context.Key = contextArangoKey
	claim := Claim{}
	claim.Key = claimArangoKey

	edge := ContextEdge{}
	bindVars := map[string]interface{}{
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
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return edge, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edge, NewServerError(err.Error())
		}
	}

	return edge, nil
}
