package gii

import (
	"fmt"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		fmt.Printf("path: %s  ", ctx.Req.URL.Path)
		fmt.Printf("RawQuery: %s  ", ctx.Req.URL.RawQuery)

		ctx.Next()

		fmt.Printf("UseTime: %v\n\n", time.Since(start))
	}
}
