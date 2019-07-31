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

func (i *Inference) Create(ctx *ServerContext) Error {
	return CreateArangoObject(ctx, i)
}

func (i *Inference) Update(ctx *ServerContext, updates map[string]interface{}) Error {
	return NewServerError("This item cannot be modified")
}

func (i *Inference) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, i)
}
