package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/GruffDebate/server/gruff"
	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

var ARANGODB_CLIENT arango.Client
var ARANGODB_POOL arango.Database

const (
	HeaderReferrerPolicy = "Referrer-Policy"
)

type securityMiddlewareOption func(*echo.Response)

func ReferrerPolicy(p string) securityMiddlewareOption {
	return func(r *echo.Response) {
		r.Header().Set(HeaderReferrerPolicy, p)
	}
}

func Secure(headers ...securityMiddlewareOption) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			for _, m := range headers {
				m(res)
			}
			return next(c)
		}
	}
}

func DBMiddleware(db arango.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			arangoCtx := gruff.ArangoContext{Context: context.Background(), DB: db}
			c.Set("Arango", arangoCtx)

			return next(c)
		}
	}
}

func InitializePayload(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("Payload", make(map[string]interface{}))
		c.Set("AppName", fmt.Sprintf("%s-%s", os.Getenv("GRUFF_NAME"), os.Getenv("GRUFF_ENV")))
		c.Set("RequestID", uuid.NewV4().String())
		c.Set("Method", c.Request().Method)
		c.Set("Endpoint", fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path))
		c.Set("Path", c.Request().URL.String())

		return next(c)
	}
}

func SettingHeaders(test bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !test {
				xGruff := c.Request().Header.Get("X-Gruff")
				if xGruff != "Gruff" {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
			}

			return next(c)
		}
	}
}

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := ServerContext(c)
		user := gruff.User{}

		authorization := c.Request().Header.Get("Authorization")
		secretKey := os.Getenv("JWT_KEY_SIGNIN")
		if authorization != "" {
			tokenSlice := strings.Split(authorization, " ")
			if len(tokenSlice) == 2 && tokenSlice[0] == "Bearer" {
				token, err := gruff.VerifyJWTToken(tokenSlice[1], secretKey)
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}

				claims := token.Claims.(jwt.MapClaims)
				if token.Valid && claims["iss"] == gruff.JWT_ISS {
					uidParse := claims["user"].(string)
					// uid := string(uidParse)
					user.Key = uidParse
					if err := user.Load(ctx); err != nil {
						return echo.NewHTTPError(http.StatusUnauthorized)
					}
				} else {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
			} else {
				user.Key = ""
			}
		} else {
			user.Key = ""
		}

		c.Set("UserContext", user)

		return next(c)
	}
}

func DetermineType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tType reflect.Type
		var parentType reflect.Type

		parts := PathParts(c.Path())
		var pathType string
		for i := 0; i < len(parts); i++ {
			pathType = parts[i]
			t := StringToType(pathType)
			if t != nil {
				if tType != nil {
					parentType = tType
				}
				tType = t
			}
		}

		c.Set("ParentType", parentType)
		c.Set("Type", tType)

		return next(c)
	}
}

func AssociationFieldNameFromPath(c echo.Context) string {
	var tType reflect.Type
	if c.Get("Type") != nil {
		tType = c.Get("Type").(reflect.Type)
	}
	path := c.Path()
	parts := strings.Split(path, "/")
	associationPath := ""
	for _, part := range parts {
		if StringToType(part) == tType {
			associationPath = part
		}
	}
	associationName := support.SnakeToCamel(associationPath)
	return associationName
}

func PathParts(path string) []string {
	parts := strings.Split(strings.Trim(path, " /"), "/")
	return parts
}

func StringToType(typeName string) (t reflect.Type) {
	switch typeName {
	case "users":
		var m gruff.User
		t = reflect.TypeOf(m)
	case "claims", "premises":
		var m gruff.Claim
		t = reflect.TypeOf(m)
	case "claim_opinions":
		var m gruff.ClaimOpinion
		t = reflect.TypeOf(m)
	case "arguments":
		var m gruff.Argument
		t = reflect.TypeOf(m)
	case "argument_opinions":
		var m gruff.ArgumentOpinion
		t = reflect.TypeOf(m)
	case "contexts":
		var m gruff.Context
		t = reflect.TypeOf(m)
	case "links":
		var m gruff.Link
		t = reflect.TypeOf(m)
	}
	return
}

func ServerContext(c echo.Context) *gruff.ServerContext {
	var tType reflect.Type
	var ParentType reflect.Type
	var user gruff.User
	var arango gruff.ArangoContext

	if c.Get("UserContext") != nil {
		user = c.Get("UserContext").(gruff.User)
	}

	if c.Get("Type") != nil {
		tType = c.Get("Type").(reflect.Type)
	}

	if c.Get("ParentType") != nil {
		ParentType = c.Get("ParentType").(reflect.Type)
	}

	if c.Get("Arango") != nil {
		arango = c.Get("Arango").(gruff.ArangoContext)
	}

	return &gruff.ServerContext{
		Context:     context.Background(),
		RequestID:   c.Get("RequestID").(string),
		Arango:      arango,
		UserContext: user,
		Test:        false,
		Type:        tType,
		ParentType:  ParentType,
		Payload:     make(map[string]interface{}),
		// TODO: what about Request, AppName, Method, etc.?
	}
}
