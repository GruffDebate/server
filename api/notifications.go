package api

import (
	_ "net/http"

	_ "github.com/GruffDebate/server/gruff"
	_ "github.com/labstack/echo"
)

/*
func ListNotifications(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	userID := ctx.UserContext.ID

	notifications := []gruff.Notification{}
	db = db.Where("user_id = ?", userID)
	db = db.Where("viewed = false")
	db = db.Order("created_at DESC")
	if err := db.Find(&notifications).Error; err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, notifications)
}

func MarkNotificationViewed(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	userID := ctx.UserContext.ID
	notificationID := c.Param("id")
	if notificationID == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	notification := gruff.Notification{}
	if err := db.First(&notification, notificationID).Error; err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if notification.UserID != userID {
		return AddError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
	}

	notification.Viewed = true
	if err := db.Save(&notification).Error; err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, notification)
}
*/
