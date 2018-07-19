package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/crm/rest/server"
)

func MountRoutes(r chi.Router) {
	field := &server.FieldHandlers{Field: Field{}.New()}
	module := &server.ModuleHandlers{Module: Module{}.New()}
	r.Group(func(r chi.Router) {
		r.Use(field.Field.Authenticator())
		r.Route("/field", func(r chi.Router) {
			r.Get("/", field.List)
			r.Get("/{id}", field.Type)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(module.Module.Authenticator())
		r.Route("/module", func(r chi.Router) {
			r.Get("/", module.List)
			r.Post("/", module.Create)
			r.Get("/{id}", module.Read)
			r.Post("/{id}", module.Edit)
			r.Delete("/{id}", module.Delete)
			r.Get("/{module}/content", module.ContentList)
			r.Post("/{module}/content", module.ContentCreate)
			r.Get("/{module}/content/{id}", module.ContentRead)
			r.Post("/{module}/content/{id}", module.ContentEdit)
			r.Delete("/{module}/content/{id}", module.ContentDelete)
		})
	})
}