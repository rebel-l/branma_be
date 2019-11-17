/*
This is the backend for the branch manager called 'branma'. It analyses your feature branches and connects it with
your JIRA tickets.

Copyright (C) 2019 Lars Gaubisch

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/bootstrap"
	"github.com/rebel-l/branma_be/endpoint/doc"
	"github.com/rebel-l/branma_be/endpoint/ping"
	"github.com/rebel-l/smis"

	"github.com/sirupsen/logrus"
)

const (
	version                    = "0.1.0"
	defaultPort                = 3000
	defaultPathToDatabase      = "./storage"
	defaultPathToSchemaScripts = "./scripts/schema"
)

var (
	log         logrus.FieldLogger
	port        *int
	svc         *smis.Service
	storagePath *string
	db          *sqlx.DB
)

func initCustomFlags() {
	/**
	  1. Add your custom service flags below, for more details see https://golang.org/pkg/flag/
	*/
	storagePath = flag.String("s", defaultPathToDatabase, "path to storage of database file")
}

func initCustom() error {
	/**
	  2. add your custom service initialisation below, e.g. database connection, caches etc.
	*/
	var err error

	db, err = bootstrap.Database(*storagePath, defaultPathToSchemaScripts, version) // TODO: use flag instead of constant
	if err != nil {
		return err
	}

	return nil
}

func initCustomRoutes() error {
	/**
	  3. Register your custom routes below
	*/

	return nil
}

func closeCustom() {
	/**
	  4. Close your connections
		TODO: include in go-project
	*/
	log.Info("Closing connections ...")

	if err := db.Close(); err != nil {
		log.Errorf("failed to close connections: %v", err)
	}
}

func main() {
	log = logrus.New()
	log.Info("Starting service: branma_be")

	initFlags()
	initService()

	if err := initCustom(); err != nil {
		log.Fatalf("Failed to initialise custom settings: %s", err)
	}
	defer closeCustom()

	if err := initRoutes(); err != nil {
		log.Fatalf("Failed to initialise routes: %s", err)
	}

	log.Infof("Service listens to port %d", *port)
	if err := svc.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func initService() {
	router := mux.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	var err error
	svc, err = smis.NewService(srv, router, log)
	if err != nil {
		log.Fatalf("failed to initialize service: %s", err)
	}
}

func initRoutes() error {
	if err := initDefaultRoutes(); err != nil {
		return fmt.Errorf("default routes failed: %s", err)
	}

	if err := initCustomRoutes(); err != nil {
		return fmt.Errorf("custom routes failed: %s", err)
	}

	return nil
}

func initDefaultRoutes() error {
	if err := ping.Init(svc); err != nil {
		return err
	}

	if err := doc.Init(svc); err != nil {
		return err
	}

	return nil
}

func initFlags() {
	initDefaultFlags()
	initCustomFlags()
	flag.Parse()
}

func initDefaultFlags() {
	port = flag.Int("p", defaultPort, "the port the service listens to")
}
