package gruff

import (
	"github.com/GruffDebate/server/support"
	"github.com/jinzhu/gorm"
)

const OBJECT_TYPE_CLAIM int = 1
const OBJECT_TYPE_ARGUMENT int = 2

const NOTIFICATION_TYPE_MOVED int = 1
const NOTIFICATION_TYPE_PARENT_MOVED int = 2
const NOTIFICATION_TYPE_NEW_ARGUMENT int = 3

type Notification struct {
	Model
	UserID   uint64      `json:"userId" sql:"not null"`
	Type     int         `json:"type" sql:"not null"`
	ItemID   *string     `json:"itemId,omitempty" sql:"type:uuid"`
	ItemType *int        `json:"itemType"`
	Item     interface{} `json:"item,omitempty" gorm:"-"`
	OldID    *string     `json:"oldId,omitempty" sql:"type:uuid"`
	OldType  *int        `json:"oldType"`
	NewID    *string     `json:"newId,omitempty" sql:"type:uuid"`
	NewType  *int        `json:"newType"`
	Viewed   bool        `json:"viewed" sql:"not null"`
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

func NotifyArgumentMoved(ctx *ServerContext, userId uint64, argId string, oldTargetId string, oldTargetType int) GruffError {
	/*
		n := Notification{
			UserID:   userId,
			Type:     NOTIFICATION_TYPE_MOVED,
			ItemID:   &argId,
			ItemType: support.IntPtr(OBJECT_TYPE_ARGUMENT),
			OldID:    &oldTargetId,
			OldType:  support.IntPtr(oldTargetType),
		}
	*/
	/*
		if err := ctx.Database.Create(&n).Error; err != nil {
			return NewServerError(err.Error())
		}
	*/
	return nil
}

func NotifyParentArgumentMoved(ctx *ServerContext, userId uint64, parentArgId string, oldTargetId string, oldTargetType int) GruffError {
	/*
		n := Notification{
			UserID:   userId,
			Type:     NOTIFICATION_TYPE_PARENT_MOVED,
			ItemID:   &parentArgId,
			ItemType: support.IntPtr(OBJECT_TYPE_ARGUMENT),
			OldID:    &oldTargetId,
			OldType:  support.IntPtr(oldTargetType),
		}
	*/
	/*
		if err := ctx.Database.Create(&n).Error; err != nil {
			return NewServerError(err.Error())
		}
	*/
	return nil
}

func NotifyNewArgument(ctx ServerContext, userId uint64, item interface{}, newArg Argument) GruffError {
	n := Notification{
		UserID:  userId,
		Type:    NOTIFICATION_TYPE_NEW_ARGUMENT,
		NewID:   &newArg.ID,
		NewType: support.IntPtr(OBJECT_TYPE_ARGUMENT),
	}
	if claim, ok := item.(Claim); ok {
		n.ItemID = &claim.ID
		n.ItemType = support.IntPtr(OBJECT_TYPE_CLAIM)
	} else if arg, ok := item.(Argument); ok {
		n.ItemID = &arg.ID
		n.ItemType = support.IntPtr(OBJECT_TYPE_ARGUMENT)
	}
	/*
		if err := ctx.Database.Create(&n).Error; err != nil {
			return NewServerError(err.Error())
		}
	*/
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
