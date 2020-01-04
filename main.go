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
	"github.com/rebel-l/branma_be/config"
	"github.com/rebel-l/branma_be/endpoint/doc"
	"github.com/rebel-l/branma_be/endpoint/ping"
	"github.com/rebel-l/branma_be/endpoint/repository"
	"github.com/rebel-l/smis"

	"github.com/sirupsen/logrus"
)

const (
	version = "0.1.0"

	timeOutWrite = 15 * time.Second
	timeOutRead  = 15 * time.Second

	defaultConfigFile = "./etc/config.json"
)

var (
	cfg           *config.Config
	configFile    *string
	databaseReset *bool
	db            *sqlx.DB
	log           logrus.FieldLogger
	svc           *smis.Service
)

func initCustomFlags() {
	/**
	  1. Add your custom service flags below, for more details see https://golang.org/pkg/flag/
	*/

	// DB
	cfg.GetDB().StoragePath = flag.String("s", cfg.GetDB().GetStoragePath(), "path to storage of database file")
	cfg.GetDB().SchemaScriptsPath = flag.String(
		"schema",
		cfg.GetDB().GetSchemaScriptPath(),
		"path to schema scripts database is created from",
	)
	databaseReset = flag.Bool("reset", false, "resets the database, NOTE: all data will be lost!")

	// GIT
	cfg.GetGit().BaseURL = flag.String(
		"git-url",
		cfg.GetGit().GetBaseURL(),
		"url of your git repository, e.g. https://github.com/rebel-l",
	)

	cfg.GetGit().ReleaseBranchPrefix = flag.String(
		"git-prefix",
		cfg.GetGit().GetReleaseBranchPrefix(),
		"prefix for release branches on your git repository",
	)

	// JIRA
	cfg.GetJira().BaseURL = flag.String(
		"jira-url",
		cfg.GetJira().GetBaseURL(),
		"url of your JIRA project, e.g. https://jira.atlassian.com/",
	)

	cfg.GetJira().Username = flag.String(
		"jira-user",
		cfg.GetJira().GetUsername(),
		"username of your login to JIRA",
	)

	cfg.GetJira().Password = flag.String(
		"jira-pw",
		cfg.GetJira().GetPassword(),
		"password of your login to JIRA",
	)
}

func initCustom() error {
	/**
	  2. add your custom service initialisation below, e.g. database connection, caches etc.
	*/
	var err error

	if *databaseReset {
		err = bootstrap.DatabaseReset(cfg.GetDB())
		if err != nil {
			return err
		}
	}

	db, err = bootstrap.Database(cfg.GetDB(), version)
	if err != nil {
		return err
	}

	return nil
}

func initCustomRoutes() error {
	/**
	  3. Register your custom routes below
	*/

	// repository
	if err := repository.Init(svc, db); err != nil {
		return err
	}

	return nil
}

func closeCustom() {
	/**
	  4. Close your connections
		nolint:godox TODO: include in go-project
	*/
	log.Info("Closing connections ...")

	if err := db.Close(); err != nil {
		log.Errorf("failed to close connections: %v", err)
	}
}

func main() {
	log = logrus.New()
	log.Info("Starting service: branma_be")

	cfg = config.New()

	initFlags()
	initService()

	initConfig()

	if err := initCustom(); err != nil {
		log.Fatalf("Failed to initialise custom settings: %s", err)
	}
	defer closeCustom()

	if err := initRoutes(); err != nil {
		log.Fatalf("Failed to initialise routes: %s", err)
	}

	log.Infof("Service listens to port %d", cfg.GetService().GetPort())
	if err := svc.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func initService() {
	router := mux.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", cfg.GetService().GetPort()),
		WriteTimeout: timeOutWrite,
		ReadTimeout:  timeOutRead,
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
	cfg.GetService().Port = flag.Int("p", cfg.GetService().GetPort(), "the port the service listens to")
	configFile = flag.String("c", defaultConfigFile, "config file in JSON format") // nolint:godox TODO: add to go-project
}

func initConfig() {
	// 1. load config from file (if exists)
	cfgFF := config.New()
	if err := cfgFF.Load(*configFile); err != nil {
		log.Warnf("load config from file failed, continue with default config/parameters: %v", err)
		return
	}

	// 2. merge with config from flags
	cfgFF.Merge(cfg)

	// 3. reset conf
	cfg = cfgFF
}
