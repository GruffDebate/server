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

// Creator

func (c *ContextEdge) Create(ctx *ServerContext) GruffError {
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

// Deleter

func (c *ContextEdge) Delete(ctx *ServerContext) GruffError {
	c.PrepareForDelete(ctx)
	patch := map[string]interface{}{
		"end": c.DeletedAt,
	}
	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err
	}
	_, aerr := col.UpdateDocument(ctx.Context, c.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}

	return nil
}

// Business methods

func FindContextEdge(ctx *ServerContext, contextArangoKey, claimArangoKey string) (ContextEdge, GruffError) {
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
