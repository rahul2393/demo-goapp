// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample helloworld is a basic App Engine flexible app.
package demogoappcloud

import (
	"fmt"
	//	"log"
	"google.golang.org/appengine/datastore"
	_ "google.golang.org/appengine/remote_api"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
)

//func main() {
//	// Set this in app.yaml when running in production.
//	projectID := os.Getenv("GCLOUD_DATASET_ID")
//
//	var err error
//	datastoreClient, err = datastore.NewClient(ctx, projectID)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	appengine.Main()
//}

func init() {
	http.HandleFunc("/_ah/remote_api", handle)
	http.HandleFunc("/_ah/health", healthCheckHandler)

}
func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!\n")

	ctx := appengine.NewContext(r)

	// Get a list of the most recent visits.
	visits, err := queryVisits(ctx, 10)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Record this visit.
	if err := recordVisit(ctx, time.Now(), r.RemoteAddr); err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Previous visits:")
	for _, v := range visits {
		fmt.Fprintf(w, "[%s] %s\n", v.Timestamp, v.UserIP)
	}
	fmt.Fprintln(w, "\nSuccessfully stored an entry of the current request.")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type visit struct {
	Timestamp time.Time
	UserIP    string
}

func recordVisit(ctx context.Context, now time.Time, userIP string) error {
	v := &visit{
		Timestamp: now,
		UserIP:    userIP,
	}

	k := datastore.NewIncompleteKey(ctx, "Visit", nil)

	_, err := datastore.Put(ctx, k, v)
	return err
}

func queryVisits(ctx context.Context, limit int64) ([]*visit, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Visit").
		Order("-Timestamp").
		Limit(10)

	visits := make([]*visit, 0)
	_, err := q.GetAll(ctx, &visits)
	return visits, err
}
