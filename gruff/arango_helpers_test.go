package gruff

import (
	"reflect"
	"testing"

	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/assert"
)

var TEST_CLIENT arango.Client
var TESTDB arango.Database
var CTX *ServerContext

func init() {
	CTX = &ServerContext{}
	TEST_CLIENT, TESTDB = InitTestDB()
}

func setupDB() {
	CTX.Arango.DB = TESTDB
}

func teardownDB() {
}

func startDBLog() {
}

func stopDBLog() {
}

func TestDefaultListQuery(t *testing.T) {
	params := DEFAULT_QUERY_PARAMETERS
	var obj ArangoObject

	obj = Claim{}
	assert.Equal(t, "FOR obj IN claims SORT obj.start DESC LIMIT 0, 20 RETURN obj", DefaultListQuery(obj, params))

	obj = Argument{}
	params.Return = support.StringPtr("obj._id")
	assert.Equal(t, "FOR obj IN arguments SORT obj.start DESC LIMIT 0, 20 RETURN obj._id", DefaultListQuery(obj, params))
}

func TestListArangoObjects(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "Let's create a new claim"}
	c2 := Claim{Title: "Ok, but we're going to need more than this"}
	c3 := Claim{Title: "How many is enough?"}
	c4 := Claim{Title: "Well, we need at least as many as the query limit"}
	c5 := Claim{Title: "If not more..."}
	c6 := Claim{Title: "Ok, yeah. Let's just do one more..."}
	c1.Create(CTX)
	c2.Create(CTX)
	c3.Create(CTX)
	c4.Create(CTX)
	c5.Create(CTX)
	c6.Create(CTX)

	query := DefaultListQuery(Claim{}, DEFAULT_QUERY_PARAMETERS.Merge(ArangoQueryParameters{Limit: support.IntPtr(5)}))
	assert.Equal(t, "FOR obj IN claims SORT obj.start DESC LIMIT 0, 5 RETURN obj", query)
	objs, err := ListArangoObjects(CTX, reflect.TypeOf(Claim{}), query, map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(objs))
	assert.Equal(t, c6.ArangoID(), objs[0].(*Claim).ArangoID())
	assert.Equal(t, c5.ArangoID(), objs[1].(*Claim).ArangoID())
	assert.Equal(t, c4.ArangoID(), objs[2].(*Claim).ArangoID())
	assert.Equal(t, c3.ArangoID(), objs[3].(*Claim).ArangoID())
	assert.Equal(t, c2.ArangoID(), objs[4].(*Claim).ArangoID())
}
