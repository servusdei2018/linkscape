// Copyright 2021 The Linkscape Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command linkscape serves the URL shortener.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/servusdei2018/linkscape/internal/db"
	"github.com/servusdei2018/linkscape/internal/router"
)

var (
	bind string
	length int
	timeout time.Duration
	lifetime time.Duration
	frequency time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 15*time.Second, "connection read/write timeout")
	flag.DurationVar(&lifetime, "lifetime", 24*time.Hour, "shortened url duration")
	flag.DurationVar(&frequency, "frequency", 1*time.Hour, "frequency between purge of expired urls")
	flag.IntVar(&length, "length", 4, "length of shortened urls")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "bind address")

	flag.Parse()
}

func main() {
	// Configure database.
	dbCfg := db.Config{
		// Remember shortened URLs for one day.
		Lifetime: 24*time.Hour,
		// Purge expired URLs every hour.
		Frequency: 1*time.Hour,
	}

	// Configure router.
	routerCfg := router.Config{
		Length: length,
	}

	r := router.New(&routerCfg, &dbCfg)

	// Configure HTTP server.
	srv := &http.Server{
		Handler: r,
		Addr: bind,
		WriteTimeout: timeout,
		ReadTimeout: timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
