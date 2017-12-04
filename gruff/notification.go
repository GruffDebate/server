package gruff

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

const OBJECT_TYPE_CLAIM int = 1
const OBJECT_TYPE_ARGUMENT int = 2

const NOTIFICATION_TYPE_MOVED int = 1
const NOTIFICATION_TYPE_PARENT_MOVED int = 2
const NOTIFICATION_TYPE_NEW_ARGUMENT int = 3

type Notification struct {
	Model
	UserID   uint64        `json:"userId" sql:"not null"`
	Type     int           `json:"type" sql:"not null"`
	ItemID   *NullableUUID `json:"itemId,omitempty" sql:"type:uuid"`
	ItemType *int          `json:"itemType"`
	Item     interface{}   `json:"item,omitempty" gorm:"-"`
	OldID    *NullableUUID `json:"oldId,omitempty" sql:"type:uuid"`
	OldType  *int          `json:"oldType"`
	NewID    *NullableUUID `json:"newId,omitempty" sql:"type:uuid"`
	NewType  *int          `json:"newType"`
	Viewed   bool          `json:"viewed" sql:"not null"`
}

func (n Notification) ValidateForCreate() GruffError {
	return ValidateStruct(n)
}

func (n Notification) ValidateForUpdate() GruffError {
	return n.ValidateForCreate()
}

func (n Notification) ValidateField(f string) GruffError {
	return ValidateStructField(n, f)
}

func NotifyArgumentMoved(ctx *ServerContext, userId uint64, argId uuid.UUID, oldTargetId uuid.UUID, oldTargetType int) GruffError {
	n := Notification{
		UserID:   userId,
		Type:     NOTIFICATION_TYPE_MOVED,
		ItemID:   NUUID(argId),
		ItemType: IntPtr(OBJECT_TYPE_ARGUMENT),
		OldID:    NUUID(oldTargetId),
		OldType:  IntPtr(oldTargetType),
	}
	if err := ctx.Database.Create(&n).Error; err != nil {
		return NewServerError(err.Error())
	}
	return nil
}

func NotifyParentArgumentMoved(ctx *ServerContext, userId uint64, parentArgId uuid.UUID, oldTargetId uuid.UUID, oldTargetType int) GruffError {
	n := Notification{
		UserID:   userId,
		Type:     NOTIFICATION_TYPE_PARENT_MOVED,
		ItemID:   NUUID(parentArgId),
		ItemType: IntPtr(OBJECT_TYPE_ARGUMENT),
		OldID:    NUUID(oldTargetId),
		OldType:  IntPtr(oldTargetType),
	}
	if err := ctx.Database.Create(&n).Error; err != nil {
		return NewServerError(err.Error())
	}
	return nil
}

func NotifyNewArgument(ctx ServerContext, userId uint64, item interface{}, newArg Argument) GruffError {
	n := Notification{
		UserID:  userId,
		Type:    NOTIFICATION_TYPE_NEW_ARGUMENT,
		NewID:   NUUID(newArg.ID),
		NewType: IntPtr(OBJECT_TYPE_ARGUMENT),
	}
	if claim, ok := item.(Claim); ok {
		n.ItemID = NUUID(claim.ID)
		n.ItemType = IntPtr(OBJECT_TYPE_CLAIM)
	} else if arg, ok := item.(Argument); ok {
		n.ItemID = NUUID(arg.ID)
		n.ItemType = IntPtr(OBJECT_TYPE_ARGUMENT)
	}
	if err := ctx.Database.Create(&n).Error; err != nil {
		return NewServerError(err.Error())
	}
	return nil
}

// Scopes

func FindArgumentMovedNotifications(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", NOTIFICATION_TYPE_MOVED).
		Where("item_type = ?", OBJECT_TYPE_ARGUMENT)
}

func FindParentArgumentMovedNotifications(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", NOTIFICATION_TYPE_PARENT_MOVED).
		Where("item_type = ?", OBJECT_TYPE_ARGUMENT)
}
