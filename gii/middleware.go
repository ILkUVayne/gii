package gii

import (
	"fmt"
)

func Base() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("base Middleware before call...\n")
		ctx.Next()
		fmt.Printf("base Middleware after call...\n")
	}
}

func V1() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v1 Middleware before call...\n")
		ctx.Next()
		fmt.Printf("v1 Middleware after call...\n")
	}
}

func V2() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v2 Middleware before call...\n")
		ctx.Next()
		fmt.Printf("v2 Middleware after call...\n")
	}
}

func V3() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v3 Middleware before call...\n")
		ctx.Next()
		fmt.Printf("v3 Middleware after call...\n")
	}
}
