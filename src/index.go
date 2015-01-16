package tinyurlwrapper

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const EMPTY = ""

//Domain class for database.
type Tinyurl struct {
	OrignalUrl string
	Tinyurl    string
	EditTime   int64
}

func init() {
	http.HandleFunc("/", handleMain)
}

//Main handle
func handleMain(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, EMPTY, EMPTY, false)
			cxt := appengine.NewContext(r)
			cxt.Errorf("handleShort: %v", err)
		}
	}()

	turl := new(Tinyurl)

	//Get original url.
	args := r.URL.Query()
	turl.OrignalUrl = args[PARAM][0]

	//Transform to tinyurl.
	ch := make(chan string)
	go getTinyUrl(w, r, turl.OrignalUrl, ch)
	turl.Tinyurl = <-ch

	//To find a existing one
	xh := make(chan *Tinyurl)
	go find(w, r, turl.OrignalUrl, xh)
	savedTinyurl := <-xh

	if savedTinyurl == nil {
		//Save in DB.
		if editTime, err := strconv.ParseInt(time.Now().Local().Format("20060102150405"), 10, 64); err == nil {
			turl.EditTime = editTime
			sh := make(chan bool)
			go save(w, r, turl, sh)
			if <-sh {
				status(w, true, turl.OrignalUrl, turl.Tinyurl, false)
			}
		} else {
			panic(err)
		}
	} else {
		status(w, true, savedTinyurl.OrignalUrl, savedTinyurl.Tinyurl, true)
	}
}

//Transform an orignalUrl to Tinyurl.
func getTinyUrl(w http.ResponseWriter, r *http.Request, orignalUrl string, ch chan string) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, EMPTY, EMPTY, false)
			cxt := appengine.NewContext(r)
			cxt.Errorf("getTinyUrl: %v", err)
			close(ch)
		}
	}()
	tingUrl := EMPTY
	if orignalUrl != EMPTY {
		cxt := appengine.NewContext(r)
		if req, err := http.NewRequest(API_METHOD, TINY+orignalUrl, nil); err == nil {
			httpClient := urlfetch.Client(cxt)
			res, err := httpClient.Do(req)
			if res != nil {
				defer res.Body.Close()
			}
			if err == nil {
				if bytes, err := ioutil.ReadAll(res.Body); err == nil {
					tingUrl = string(bytes)
					ch <- tingUrl
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		ch <- EMPTY
	}
}

//Save a Tinyurl in database.
func save(w http.ResponseWriter, r *http.Request, tinyurl *Tinyurl, ch chan bool) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, EMPTY, EMPTY, false)
			cxt := appengine.NewContext(r)
			cxt.Errorf("find: %v", err)
			close(ch)
		}
	}()

	//Save in db.
	cxt := appengine.NewContext(r)
	if _, err := datastore.Put(cxt, datastore.NewIncompleteKey(cxt, "Tinyurl", nil), tinyurl); err == nil {
		ch <- true
	} else {
		panic(err)
	}
}

//To find an existing url that has been transformed by tinyurl before.
//A validate Tinyurl returns back through ch, otherwise a nil.
func find(w http.ResponseWriter, r *http.Request, url string, ch chan *Tinyurl) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, EMPTY, EMPTY, false)
			cxt := appengine.NewContext(r)
			cxt.Errorf("find: %v", err)
			close(ch)
		}
	}()

	cxt := appengine.NewContext(r)
	q := datastore.NewQuery("Tinyurl").Filter("OrignalUrl =", url)
	turls := make([]Tinyurl, 0)
	if _, err := q.GetAll(cxt, &turls); err == nil {
		if len(turls) > 0 {
			ch <- &turls[0]
		} else {
			ch <- nil
		}
	} else {
		panic(err)
	}
}

//Response json to browser.
func status(w http.ResponseWriter, ok bool, q string, res string, stored bool) {
	s := fmt.Sprintf(`{"status":%s,   "q":"%s", "result":"%s", "stored":%s }`,
		strconv.FormatBool(ok),
		q,
		res,
		strconv.FormatBool(stored))

	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}
