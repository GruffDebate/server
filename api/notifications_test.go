package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GruffDebate/server/gruff"
	"github.com/petmondo/petmondo-server/support"
	"github.com/stretchr/testify/assert"
)

func TestListNotifications(t *testing.T) {
	setup()
	defer teardown()

	u1 := gruff.User{Name: "User1", Username: "user1", Email: "email1@gruff.org"}
	u2 := gruff.User{Name: "User2", Username: "user2", Email: "email2@gruff.org"}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r := New(tokenForTestUser(u1))

	n1 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM)}
	n2 := gruff.Notification{UserID: u2.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM)}
	n3 := gruff.Notification{UserID: u2.ID, Type: gruff.NOTIFICATION_TYPE_NEW_ARGUMENT, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM), NewID: gruff.NewNUUID(), NewType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	n4 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM), Viewed: true}
	n5 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	n6 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_NEW_ARGUMENT, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), NewID: gruff.NewNUUID(), NewType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	n7 := gruff.Notification{UserID: u2.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	TESTDB.Create(&n1)
	TESTDB.Create(&n2)
	TESTDB.Create(&n3)
	TESTDB.Create(&n4)
	TESTDB.Create(&n5)
	TESTDB.Create(&n6)
	TESTDB.Create(&n7)

	expectedResults, _ := json.Marshal([]gruff.Notification{n6, n5, n1})

	r.GET("/api/notifications")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestMarkNotificationViewed(t *testing.T) {
	setup()
	defer teardown()

	u1 := gruff.User{Name: "User1", Username: "user1", Email: "email1@gruff.org"}
	u2 := gruff.User{Name: "User2", Username: "user2", Email: "email2@gruff.org"}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r := New(tokenForTestUser(u1))

	n1 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM)}
	n2 := gruff.Notification{UserID: u2.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM)}
	n3 := gruff.Notification{UserID: u2.ID, Type: gruff.NOTIFICATION_TYPE_NEW_ARGUMENT, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM), NewID: gruff.NewNUUID(), NewType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	n4 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_CLAIM), Viewed: true}
	n5 := gruff.Notification{UserID: u1.ID, Type: gruff.NOTIFICATION_TYPE_MOVED, ItemID: gruff.NewNUUID(), ItemType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT), OldID: gruff.NewNUUID(), OldType: support.IntPtr(gruff.OBJECT_TYPE_ARGUMENT)}
	TESTDB.Create(&n1)
	TESTDB.Create(&n2)
	TESTDB.Create(&n3)
	TESTDB.Create(&n4)
	TESTDB.Create(&n5)

	n5.Viewed = true
	expectedResults, _ := json.Marshal(n5)

	m := map[string]interface{}{}

	r.POST(fmt.Sprintf("/api/notifications/%d", n5.ID))
	r.SetBody(m)
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)

	TESTDB.First(&n5)
	assert.True(t, n5.Viewed)
}
