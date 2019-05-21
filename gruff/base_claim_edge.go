package gruff

import (
	"fmt"
)

// BaseClaim is an edge pointing from an Argument to the Claim on which it is based
// (the true/false part of the Argument)
type BaseClaimEdge struct {
	Edge
	From string `json:"_from,omitempty"`
	To   string `json:"_to,omitempty"`
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

// Creator

func (bc *BaseClaimEdge) Create(ctx *ServerContext) GruffError {
	col, err := ctx.Arango.CollectionFor(bc)
	if err != nil {
		return err
	}

	bc.PrepareForCreate()

	_, aerr := col.CreateDocument(ctx.Context, bc)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}

func (bc *BaseClaimEdge) Delete(ctx *ServerContext) GruffError {
	bc.PrepareForDelete()
	patch := map[string]interface{}{
		"end": bc.DeletedAt,
	}
	col, err := ctx.Arango.CollectionFor(bc)
	if err != nil {
		return err
	}
	_, aerr := col.UpdateDocument(ctx.Context, bc.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}

	return nil
}
