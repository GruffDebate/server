package gruff

import (
	"fmt"
)

// Inference is an edge from the target (a Claim or Argument) of an Argument
// to the Argument that is making the inference
type Inference struct {
	Edge
}

// ArangoObject interface

func (i Inference) CollectionName() string {
	return "inferences"
}

func (i Inference) ArangoKey() string {
	return i.Key
}

func (i Inference) ArangoID() string {
	return fmt.Sprintf("%s/%s", i.CollectionName(), i.ArangoKey())
}

func (i Inference) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

// Creator

func (i *Inference) Create(ctx *ServerContext) GruffError {
	col, err := ctx.Arango.CollectionFor(i)
	if err != nil {
		return err
	}

	i.PrepareForCreate(ctx)

	_, aerr := col.CreateDocument(ctx.Context, i)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}

func (i *Inference) Delete(ctx *ServerContext) GruffError {
	i.PrepareForDelete(ctx)
	patch := map[string]interface{}{
		"end": i.DeletedAt,
	}
	col, err := ctx.Arango.CollectionFor(i)
	if err != nil {
		return err
	}
	_, aerr := col.UpdateDocument(ctx.Context, i.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}

	return nil
}
