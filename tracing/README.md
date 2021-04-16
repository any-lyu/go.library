
## Package Feature Overview

- Jeager tracing middleware for echo web framework
- Easy to tracing any component like MySQL, Redis, Mongo
- Distributed context propagation



### Example

```go
package main

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.innotechx.com/shm/go.common/database/sqlx"
	
	"gitlab.innotechx.com/shm/go.common/http/echo/middleware"
	"gitlab.innotechx.com/shm/go.common/tracing"
	"gitlab.innotechx.com/shm/go.common/http/apicall"
)

func main() {
	
	e := echo.New()
	
	serviceName := "example"
	endpoint := "http://example"
	
	// Init global tracer
	// SampleRate must in the range between 0.0 and 1.0
	// 1.0 for sample all requests, default is 0.001.
	
	if closer, err := tracing.InitGlobalTracer(serviceName, endpoint, tracing.DefaultSampleRate); err != nil {
		panic("init tracing failed: " + err.Error())
	} else {
		defer closer.Close()
	}

	// Install tracing middleware to echo web framework.
	e.Use(middleware.Tracing(serviceName))

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	
	// If you want to tracing your redis, start a span manually.
	sp := tracing.Redis().StartSpan(c.Request().Context())
	
	_, err = redisIns.HGetAll("device_user_map").Result()
	
	// Set any tag you like to the span, optional.
	sp.SetTag("example", "example")
	
	// Log something you like to the span, optional.
	sp.LogKV("server", "server")
	
	tracing.Redis().FinishSpan(sp, err)
	
	// finish span manually.
	return c.String(http.StatusOK, "Hello, World!")
}

func propagateTracing()  {
	
	// call any api with tracing, tracing enable by default.
	
    _ := apicall.Post(ctx, uri, nil, struct {}{}, &result)
    
}

func initMysqlDB() (sqlx.DB, error) {
	
	// Init your db here.
	
	// Enable MySQL tracing by default.
	
	return sqlx.WrapDB(db), nil
}
```

### Help
- tracing endpoints list in this url:http://km.innotechx.com/display/SHOP/Tracing

