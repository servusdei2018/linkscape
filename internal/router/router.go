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

package router

import (
	"fmt"
	"net/http"

	"github.com/servusdei2018/linkscape/internal/db"
	"github.com/servusdei2018/linkscape/internal/site/static"

	"github.com/gorilla/mux"
)

var (
	database *db.DB
	gen      *generator
	router   *mux.Router
)

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/shorten", shortener).Methods("GET")
	router.HandleFunc("/{key:[a-zA-Z]{4}}", redirector).Methods("GET")
}

// New configures the database and returns the router.
func New(cfg *Config, dbcfg *db.Config) *mux.Router {
	gen = newGenerator(cfg.Length)
	database = db.New(dbcfg)
	return router
}

// index serves the application's index page.
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index reached")
	w.WriteHeader(http.StatusOK)
	w.Write(static.Index)
}

// shortener shortens a URL and returns the shortened link.
func shortener(w http.ResponseWriter, r *http.Request) {
	// Parse the url to shorten.
	vars := r.URL.Query()
	url, present := vars["url"]
	if !present {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	// Generate a unique identifier.
	uid := gen.Next()
	// Store the URL by uid.
	database.Set(uid, &url[0])
	// Return the shortened url.
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, uid)))
}

// redirector redirects traffic.
func redirector(w http.ResponseWriter, r *http.Request) {
	// Parse the key.
	vars := mux.Vars(r)
	key, present := vars["key"]
	if !present {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	// Retrieve the URL from the database.
	val, found := database.Get(key)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	// Redirect to the URL.
	http.Redirect(w, r, *val, http.StatusTemporaryRedirect)
}
