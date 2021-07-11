# ginlambda
A tiny adapter that enables Gin for Lambda

## Quick Start

```
go get -u github.com/dgravesa/ginlambda
```

Set up your Gin routes as you normally would, and instead of calling *lambda.Start()* just pass your Gin engine instance to `ginlambda.Start()`.

```go
package main

import (
    "net/http"

    "github.com/dgravesa/ginlambda"
    "github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
    r = gin.Default()
    r.GET("/greeting", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })
}

func main() {
    // equivalent to lambda.Start(handler)
    ginlambda.Start(r)
}
```

## Also Check Out

[AWS Labs Framework Adaptors](https://github.com/awslabs/aws-lambda-go-api-proxy) - This is an AWS-supported package that does essentially the same job. Plus, it provides adaptors for a bunch of other frameworks as well.

There are two benefits to using `ginlambda` instead. First and foremost, `ginlambda` has a much smaller dependency tree,
as it's only meant to support Gin. The AWS Labs module has indirect dependencies from a plethora of frameworks,
so for example you're depending on Iris and its dependencies even if you aren't using Iris.

Second, `ginlambda` can save you two lines of code by using `ginlambda.Start()`.
This translates to roughly 200 bytes, saving you precious disk space.
