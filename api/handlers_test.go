package api

import (
	_ "errors"

	"github.com/GruffDebate/server/gruff"
	arango "github.com/arangodb/go-driver"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

var CTX *gruff.ServerContext
var TEST_CLIENT arango.Client
var TESTDB arango.Database
var ROLE string = "user"

var TESTTOKEN string
var READ_ONLY bool = false

func init() {
	TEST_CLIENT, TESTDB = gruff.InitTestDB()
}

func setup() {
	//TESTDB = INITDB.Begin()

	if CTX == nil {
		CTX = &gruff.ServerContext{}
	}

	CTX.Arango.DB = TESTDB
}

func teardown() {
	//TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	//TESTDB.LogMode(true)
}

func stopDBLog() {
	//TESTDB.LogMode(false)
}

func Router() *echo.Echo {
	return SetUpRouter(TestMiddlewareConfigurer{})
}

type TestMiddlewareConfigurer struct{}

func (mc TestMiddlewareConfigurer) ConfigureDefaultApiMiddleware(root *echo.Echo) *echo.Echo {
	root.Use(middleware.Recover())
	root.Use(middleware.CORS())
	root.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	root.Use(DBMiddleware(TESTDB))
	root.Use(DetermineType)
	root.Use(InitializePayload)

	return root
}

func (mc TestMiddlewareConfigurer) ConfigurePublicApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	public := api.Group("/api")
	public.Use(SettingHeaders(true))
	public.Use(Session)

	return public
}

func (mc TestMiddlewareConfigurer) ConfigurePrivateApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/api")
	private.Use(middleware.Gzip())
	private.Use(SettingHeaders(true))
	// private.Use(SetUpTestUser(ROLE))
	// private.Use(SetTestUserToken)
	private.Use(Session)

	return private
}

// func SetUpTestUser(role string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			password := "password"
// 			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			
// 			user := argon.User{}
// 			dbFind := TESTDB.Unscoped().First(&user, 999)
// 			if dbFind.RecordNotFound() {
// 				u := gruff.User{
					// 	Name:     name,
					// 	Username: username,
					// 	Email:    email,
					// 	Password: "123456",
					// }
					// password := u.Password
					// u.Password = ""
					// u.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

					// u.Create(CTX)
// 			} else {
// 				user.HashedPassword = hashedPassword
// 				TESTDB.Unscoped().Save(&user)
// 			}

// 			c.Set("User", user)

// 			return next(c)
// 		}
// 	}
// }

// func SetTestUserToken(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		user := c.Get("User").(argon.User)
// 		expireAt := argon.JWTTokenExpirationDate()
// 		jwt, _ := argon.IssueJWToken(user.ID, []string{"user"}, expireAt)
// 		c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
// 		return next(c)
// 	}
// }
