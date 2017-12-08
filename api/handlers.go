package api

import (
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/GruffDebate/server/gruff"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func List(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	db = BasicJoins(ctx, c, db)
	db = BasicPaging(ctx, c, db)

	items := reflect.New(reflect.SliceOf(ctx.Type)).Interface()
	err := db.Find(items).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	items = itemsOrEmptySlice(ctx.Type, items)

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = items
		return c.JSON(http.StatusOK, ctx.Payload)
	} else {
		return c.JSON(http.StatusOK, items)
	}
}

func Create(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	item := reflect.New(ctx.Type).Interface()
	if err := c.Bind(item); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	valerr := BasicValidationForCreate(ctx, c, item)
	if valerr != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(valerr.Error()))
	}

	if gruff.IsIdentifier(ctx.Type) {
		gruff.SetCreatedByID(item, ctx.UserContext.ID)
	}

	dberr := db.Create(item).Error
	if dberr != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(dberr.Error()))
	}

	return c.JSON(http.StatusCreated, item)
}

func Get(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item := reflect.New(ctx.Type).Interface()

	db = BasicJoins(ctx, c, db)
	//db = BasicFetch(ctx, c, db, id)

	err := db.Where("id = ?", id).First(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, item)
}

func Update(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item := reflect.New(ctx.Type).Interface()
	err := db.Where("id = ?", id).First(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	err = BasicValidationForUpdate(ctx, c, item)
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := c.Bind(item); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	dberr := db.Set("gorm:save_associations", false).Save(item).Error
	if dberr != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(dberr.Error()))
	}

	return c.JSON(http.StatusAccepted, item)
}

func Delete(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item := reflect.New(ctx.Type).Interface()
	err := db.Where("id = ?", id).First(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	err = db.Delete(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, item)
}

func Destroy(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item := reflect.New(ctx.Type).Interface()
	err := db.Where("id = ?", id).First(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	err = db.Unscoped().Delete(item).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, item)
}

func AddAssociation(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	parentID := c.Param("parentId")
	id := c.Param("id")

	parentItem := reflect.New(ctx.ParentType).Interface()
	if err := db.Where("id = ?", parentID).First(parentItem).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	item := reflect.New(ctx.Type).Interface()
	if err := db.Where("id = ?", id).First(item).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	associationName := AssociationFieldNameFromPath(c)
	if err := db.Model(parentItem).Association(associationName).Append(item).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusCreated, item)
}

func ReplaceAssociation(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	parentID := c.Param("parentId")

	model := gruff.ReplaceMany{}
	if err := c.Bind(&model); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	parentItem := reflect.New(ctx.ParentType).Interface()
	if err := db.Where("id = ?", parentID).First(parentItem).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	items := reflect.New(reflect.SliceOf(ctx.Type)).Interface()
	err := db.Where("id in (?)", model.IDS).Find(items).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	associationName := AssociationFieldNameFromPath(c)
	if err := db.Model(parentItem).Association(associationName).Replace(items).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, items)
}

func RemoveAssociation(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	parentID := c.Param("parentId")
	id := c.Param("id")

	parentItem := reflect.New(ctx.ParentType).Interface()
	if err := db.Where("id = ?", parentID).First(parentItem).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	item := reflect.New(ctx.Type).Interface()
	if err := db.Where("id = ?", id).First(item).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	associationName := AssociationFieldNameFromPath(c)
	if err := db.Model(parentItem).Association(associationName).Delete(item).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusOK, item)
}

func BasicJoins(ctx *gruff.ServerContext, c echo.Context, db *gorm.DB) *gorm.DB {
	db = joinsFor(db, ctx)
	return db
}

func BasicFetch(ctx *gruff.ServerContext, c echo.Context, db *gorm.DB, uid int) *gorm.DB {
	if uid > 0 {
		id := uint(uid)
		path := c.Path()
		db = fetchFor(db, path, id)
		return db
	}
	return db
}

func fetchFor(db *gorm.DB, path string, userId uint) *gorm.DB {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		switch part {
		case "users":
		}
	}
	return db
}

func joinsFor(db *gorm.DB, ctx *gruff.ServerContext) *gorm.DB {
	t := ctx.Type
	elemT := t
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}
	for i := 0; i < elemT.NumField(); i++ {
		f := elemT.Field(i)
		tag := elemT.Field(i).Tag
		fetch := tag.Get("fetch")
		if fetch == "eager" {
			db = db.Preload(f.Name)
		}
	}
	return db
}

func BasicPaging(ctx *gruff.ServerContext, c echo.Context, db *gorm.DB, opts ...bool) *gorm.DB {
	queryTC := true
	if len(opts) > 0 {
		queryTC = opts[0]
	}

	st := c.QueryParam("start")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit > 0 && queryTC {
		QueryTotalCount(ctx, c)
	}

	if st != "" {
		startIdx, _ := strconv.Atoi(st)
		if startIdx > 0 {
			db = db.Offset(startIdx)
		}
	}

	if limit > 0 {
		db = limitQueryByConfig(ctx, db, "", limit)
	}

	return db
}

func QueryTotalCount(ctx *gruff.ServerContext, c echo.Context) {
	item := reflect.New(ctx.Type).Interface()
	var n int

	ctx.Database.Model(item).
		Select("COUNT(*)").
		Row().
		Scan(&n)

	ctx.Payload["ct"] = n
}

func limitQueryByConfig(ctx *gruff.ServerContext, db *gorm.DB, key string, requestLimit int) *gorm.DB {
	dbLimit := requestLimit
	limitStr := os.Getenv(key)
	limit, err := strconv.Atoi(limitStr)
	if err == nil {
		if dbLimit <= 0 || (limit > 0 && limit < dbLimit) {
			dbLimit = limit
		}
	}
	if dbLimit > 0 {
		db = db.Limit(dbLimit)
	}
	return db
}

func itemsOrEmptySlice(t reflect.Type, items interface{}) interface{} {
	if reflect.ValueOf(items).IsNil() {
		items = reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	}
	return items
}

func BasicValidationForCreate(ctx *gruff.ServerContext, c echo.Context, item interface{}) gruff.GruffError {
	if gruff.IsValidator(reflect.TypeOf(item)) {
		validator := item.(gruff.Validator)
		return validator.ValidateForCreate()
	}

	return nil
}

func BasicValidationForUpdate(ctx *gruff.ServerContext, c echo.Context, item interface{}) error {
	if gruff.IsValidator(ctx.Type) {
		return gruff.ValidateStructFields(item)
	}

	return nil
}
