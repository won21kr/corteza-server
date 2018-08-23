package crm

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/SentimensRG/sigctx"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/crm/rest"

	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"
)

func Init() error {
	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	// start/configure database connection
	factory.Database.Add("default", config.db.dsn)
	db, err := factory.Database.Get()
	if err != nil {
		return err
	}
	switch config.db.profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
		// @todo: profiling as an external service?
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: config.http.pretty,
		Trace:  config.http.tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	return nil
}

func Start() error {
	var ctx = sigctx.New()

	log.Println("Starting http server on address " + config.http.addr)
	listener, err := net.Listen("tcp", config.http.addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", config.http.addr))
	}

	// JWT Auth
	jwtAuth, err := auth.JWT()
	if err != nil {
		return errors.Wrap(err, "Error creating JWT Auth object")
	}

	r := chi.NewRouter()
	r.Use(jwtAuth.Verifier(), jwtAuth.Authenticator())

	// mount routes
	MountRoutes(r, config, rest.MountRoutes(jwtAuth))

	go http.Serve(listener, r)
	<-ctx.Done()

	return nil
}