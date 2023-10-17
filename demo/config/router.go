package config

import (
	"gii/demo/controller"
	"gii/gii"
	"net/http"
)

func Router() *gii.Engine {
	Router := gii.Default()

	r1 := Router.Group("/v1").Use(gii.V1())
	{
		r1.Get("/ping", controller.Handle)
		r1.Get("/hello", controller.ReJson)
		r1.Post("/hello", controller.ReJson)
		r1.Get("/rexml", controller.ReXML)
	}

	r2 := Router.Group("/v2").Use(gii.V2())
	{
		r2.Get("/ping", controller.Handle)
		r2.Get("/hello", controller.ReJson)
		r2.Post("/hello", controller.ReJson)
		r2.Get("/rexml", controller.ReXML)
	}

	r3 := r2.Group("/v3").Use(gii.V3())
	{
		r3.Get("/ping", controller.Handle)
		r3.Get("/hello", controller.ReJson)
		r3.Post("/hello", controller.ReJson)
		r3.Get("/rexml", controller.ReXML)
		r3.Get("/panic", controller.PanicC)
	}
	//
	Router.Get("/ping", controller.Handle)
	Router.Get("/hello", controller.ReJson)
	Router.Post("/hello", controller.ReJson)
	Router.Get("/rexml", controller.ReXML)
	Router.Get("/panic", controller.PanicC)

	Router.Get("id/:name", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "id/:name")
	})
	Router.Get("id/name", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "id/name")
	})
	Router.Get("id/name/sd", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "id/name/sd")
	})
	Router.Get("id/:name/asdas", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "id/:name/asdas")
	})
	// /favicon.ico
	Router.Get(":id", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", ":id")
	})
	Router.Delete("book/:id", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "book/:id")
	})

	Router.Patch("book1", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "book1")
	})

	Router.Put("book2", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "book2")
	})

	Router.Options("book3", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "book3")
	})

	Router.Head("book4", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "book4")
	})

	Router.Any("bookAny", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "path: %s", "bookAny")
	})
	return Router
}
