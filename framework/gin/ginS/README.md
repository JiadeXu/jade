# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"a-projects/geekbang/framework/gin"
	"a-projects/geekbang/framework/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
