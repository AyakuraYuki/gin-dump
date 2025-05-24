# gin-dump

> Notice: This `gin-dump` is similar to `sambar/gin-dump`, but have a bit different to that one.

* Gin middleware/handler to dump header/body of request and response.
* Very helpful for debugging your applications.
* More beautiful and compact than `samber/gin-dump`

## Content-Type support

* application/json
* application/x-www-form-urlencoded

## Usage

Download and install it / import it in your code:

```shell
$ go get github.com/AyakuraYuki/gin-dump
```

```go
import "github.com/AyakuraYuki/gin-dump"
```

### Examples:

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/gin-dump"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// prints to stdout
	router.Use(gin_dump.Dump())

	// custom callback with default options
	router.Use(gin_dump.DumpFunc(func(dumpStr string) {
		fmt.Println(dumpStr)
	}))

	// dump with options and custom callback
	opts := []gin_dump.Option{
		gin_dump.WithShowReq(true),
		gin_dump.WithShowCookies(false),
		gin_dump.WithCallback(func(dumpStr string) {
			fmt.Println(dumpStr)
		}),
	}
	router.Use(gin_dump.DumpWithOptions(opts...))

	// ...
	router.Run()
}

```
