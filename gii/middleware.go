package gii

import (
	"fmt"
)

func Base() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("base Middleware call...\n")
	}
}

func V1() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v1 Middleware call...\n")
	}
}

func V2() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v2 Middleware call...\n")
	}
}

func V3() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("v3 Middleware call...\n")
	}
}
