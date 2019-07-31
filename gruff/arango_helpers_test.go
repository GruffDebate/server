package gruff

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/assert"
)

var TEST_CLIENT arango.Client
var TESTDB arango.Database
var CTX *ServerContext
var DEFAULT_USER User

func init() {
	CTX = &ServerContext{}
	TEST_CLIENT, TESTDB = InitTestDB()
	CTX.Arango.DB = TESTDB

	user := User{
		Name:            "Big Billy Goat Gruff",
		Username:        "BigBillyGoat",
		Email:           "bbg@gruff.org",
		Image:           "https://miro.medium.com/max/1400/1*h765MiOJBkf7fqPdrQDCPQ.jpeg",
		Curator:         false,
		Admin:           false,
		URL:             "https://github.com/canonical-debate-lab/paper",
		EmailVerifiedAt: support.TimePtr(time.Now()),
	}
	err := user.Create(CTX)
	if err != nil {
		fmt.Println("ERROR Creating test user:", err.Error())
	}

	DEFAULT_USER = user
	CTX.UserContext = user
}

func setupDB() {
	CTX.Arango.DB = TESTDB
	CTX.UserContext = DEFAULT_USER
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

	obj = &Claim{}
	assert.Equal(t, "FOR obj IN claims FILTER obj.end == null SORT obj.start DESC LIMIT 0, 20 RETURN obj", DefaultListQuery(obj, params))

	obj = &Argument{}
	params.Return = support.StringPtr("obj._id")
	assert.Equal(t, "FOR obj IN arguments FILTER obj.end == null SORT obj.start DESC LIMIT 0, 20 RETURN obj._id", DefaultListQuery(obj, params))
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
	CTX.RequestAt = nil
	c2.Create(CTX)
	CTX.RequestAt = nil
	c3.Create(CTX)
	CTX.RequestAt = nil
	c4.Create(CTX)
	CTX.RequestAt = nil
	c5.Create(CTX)
	CTX.RequestAt = nil
	c6.Create(CTX)
	CTX.RequestAt = nil

	query := DefaultListQuery(&Claim{}, DEFAULT_QUERY_PARAMETERS.Merge(ArangoQueryParameters{Limit: support.IntPtr(5)}))
	assert.Equal(t, "FOR obj IN claims FILTER obj.end == null SORT obj.start DESC LIMIT 0, 5 RETURN obj", query)
	objs, err := ListArangoObjects(CTX, reflect.TypeOf(Claim{}), query, map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(objs))
	assert.Equal(t, c6.ArangoID(), objs[0].(*Claim).ArangoID())
	assert.Equal(t, c5.ArangoID(), objs[1].(*Claim).ArangoID())
	assert.Equal(t, c4.ArangoID(), objs[2].(*Claim).ArangoID())
	assert.Equal(t, c3.ArangoID(), objs[3].(*Claim).ArangoID())
	assert.Equal(t, c2.ArangoID(), objs[4].(*Claim).ArangoID())
}

func TestGetArangoObject(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "Let's create a new claim for GetArangoObject"}
	err := c1.Create(CTX)
	assert.NoError(t, err)

	context := Context{ShortName: "GetArangoObject Context", Title: "Required Title", URL: "https://en.wikipedia.org/wiki/Context"}
	err = context.Create(CTX)
	assert.NoError(t, err)

	obj, err := GetArangoObject(CTX, reflect.TypeOf(c1), c1.ArangoKey())
	assert.NoError(t, err)
	claim := obj.(*Claim)
	assert.Equal(t, c1.ArangoID(), claim.ArangoID())

	obj, err = GetArangoObject(CTX, reflect.TypeOf(context), context.ArangoKey())
	assert.NoError(t, err)
	context1 := obj.(*Context)
	assert.Equal(t, context.ArangoID(), context1.ArangoID())

	obj, err = GetArangoObject(CTX, reflect.TypeOf(context), "blah blah blah")
	assert.Error(t, err)
	assert.Nil(t, obj)
	assert.Equal(t, "document not found", err.Error())
}
