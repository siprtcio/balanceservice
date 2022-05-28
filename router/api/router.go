package api

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	. "github.com/siprtcio/balanceservice/configs"
	"github.com/siprtcio/balanceservice/modules/cache"
)

// Handler
func health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy!")
}

//-----
// API Routers
//-----
func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()
	// Context
	e.Use(NewContext())
	// Customization
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("api")
	e.Logger.SetLevel(GetLogLvl())

	// Gzip
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Cache
	e.Use(cache.Cache())

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Validator = NewValidator()

	e.Use(mw.RequestIDWithConfig(mw.RequestIDConfig{
		Generator: func() string {
			return genUUID()
		},
	}))
	e.GET("v1/Accounts/:vendor_auth_id/Balances/:tenant_auth_id", handler(GetBalanceByVendorAuthIdAuthId))
	e.GET("v1/Accounts/:auth_id/Balances", handler(GetBalanceByAuthId))
	e.GET("v1/Accounts/:auth_id/Balance", handler(GetBalanceByAuthId))
	e.PATCH("v1/Accounts/:vendor_auth_id/Balances/:tenant_auth_id", handler(UpdateBalanceByVendorAuthIdAuthId))
	return e
}

//-----
// API Routers
//-----
func PrivRouters() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Context自定义
	e.Use(NewContext())

	// Customization
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("api")
	e.Logger.SetLevel(GetLogLvl())

	// Gzip
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Cache
	e.Use(cache.Cache())

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Validator = NewValidator()

	e.Use(mw.RequestIDWithConfig(mw.RequestIDConfig{
		Generator: func() string {
			return genUUID()
		},
	}))
	e.GET("v1/health", health)
	e.GET("v1/Accounts/:vendor_auth_id/Balances/:tenant_auth_id", handler(GetBalanceByVendorAuthIdAuthId))
	e.GET("v1/Accounts/:auth_id/Balance", handler(GetBalanceByAuthId))
	e.GET("v1/Accounts/:auth_id/Balances", handler(GetBalanceByAuthId))
	e.PATCH("v1/Accounts/:vendor_auth_id/Balances/:tenant_auth_id", handler(UpdateBalanceByVendorAuthIdAuthId))
	e.PATCH("v1/Accounts/:auth_id/Balances", handler(UpdateBalanceByAuthId))
	return e
}

type (
	HandlerFunc func(*Context) error
)

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
func genUUID() string {
	// Create a Version 4 UUID.
	requestUUID, err := uuid.NewV4()
	if err != nil {
		//add error handler
	}
	return requestUUID.String()
}
