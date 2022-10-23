package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	q := QueryService{db: db}
	apis := e.Group("/api")
	apis.Use(middleware.AddTrailingSlash())
	apis.POST("/query", func(c echo.Context) error {
		req := QueryReq{}
		err := c.Bind(&req)
		if err != nil {
			return c.String(400, "bad request: struct error")
		}
		query, err := q.Query(c.Request().Context(), req)

		if err != nil {
			if errors.Is(err, BadRequest) {
				return c.String(400, err.Error())
			}
			return err
		}
		return c.JSON(200, query)
	})
	apis.GET("/btc/blocks/:hash", func(c echo.Context) error {
		hash := c.Param("hash")
		rt, err := q.GetBlockByHash(c.Request().Context(), hash)
		if err != nil {
			if errors.Is(err, BadRequest) {
				return c.String(400, err.Error())
			} else if errors.Is(err, NotFound) {
				return c.String(404, err.Error())
			}
			return err
		}
		return c.JSON(200, rt)
	})
	static := Static{base: os.Getenv("STATIC")}
	e.GET("/*", static.Serve)
	e.Logger.Fatal(e.Start(os.Getenv("LISTEN")))
}

type Static struct {
	base string
}

func (s Static) Serve(c echo.Context) error {
	p := c.Param("*")
	p, err := url.PathUnescape(p)
	if err != nil {
		return c.String(400, fmt.Sprintf("failed to unescape path variable: %s", err.Error()))
	}

	// fs.FS.Open() already assumes that file names are relative to FS root path and considers name with prefix `/` as invalid
	name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))
	path := filepath.Join(s.base, name)

	fi, err := os.Stat(path)
	if err != nil {
		name = "index.html"
		path = filepath.Join(s.base, name)
		fi, err = os.Stat(path)
		if err != nil {
			return c.String(404, "not found")
		}
	}
	p = c.Request().URL.Path
	if fi.IsDir() && len(p) > 0 && p[len(p)-1] != '/' {
		return c.Redirect(http.StatusMovedPermanently, (p + "/"))
	}
	return c.File(path)
}
