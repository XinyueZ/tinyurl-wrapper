package tinyurlwrapper

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

//Auto-update handler.
func handleAutoUpdate(w http.ResponseWriter, r *http.Request) {
	ch := make(chan []Tinyurl)
	go getAll(w, r, ch)
	all := <-ch
	if len(all) > 0 {
		update(w, r, all)
	}
}

//Get all saved Tinyurls.
func getAll(w http.ResponseWriter, r *http.Request, ch chan []Tinyurl) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, EMPTY, EMPTY, false)
			cxt := appengine.NewContext(r)
			cxt.Errorf("getAll: %v", err)
			close(ch)
		}
	}()

	cxt := appengine.NewContext(r)
	q := datastore.NewQuery("Tinyurl").Filter("EditTime >", 0)
	turls := make([]Tinyurl, 0)
	if _, err := q.GetAll(cxt, &turls); err == nil {
		if len(turls) > 0 {
			ch <- turls
		} else {
			ch <- nil
		}
	} else {
		panic(err)
	}
}

//Update existings.
func update(w http.ResponseWriter, r *http.Request, tinyurls []Tinyurl) {
	total := len(tinyurls)
	ch := make(chan int, total)

	for i := 0; i < total; i++ {
		go func(index int, pturl *Tinyurl) {
			cxt := appengine.NewContext(r)
			q := datastore.NewQuery("Tinyurl").Filter("OrignalUrl=", pturl.OrignalUrl)
			tinyurls := make([]Tinyurl, 0)
			keys, _ := q.GetAll(cxt, &tinyurls)
			build(w, r, keys[0], pturl)
			ch <- i
		}(i, &(tinyurls[i]))
	}

	for i := 0; i < total; i++ {
		<-ch
	}
}
