package tinyurlwrapper

import (
	"appengine"
	"appengine/urlfetch"


	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
)

func init() {
	http.HandleFunc("/", handleShort)
}

func handleShort(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, false, "short", "", "", false)
		}
	}()

	//cxt := appengine.NewContext(r)

	args := r.URL.Query()
	q := args[PARAM][0]

	ch := make(chan string)
	go getTinyUrl(r, q, ch)
	res := <-ch

	status(w, true, "short", q, res, false)
}

func getTinyUrl(r *http.Request, orignalUrl string, ch chan string) {
	tingUrl := ""
	if orignalUrl != "" {
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
				} else {
					cxt.Errorf("getTinyUrl read: %v", err)
					tingUrl = orignalUrl
				}
			} else {
				cxt.Errorf("getTinyUrl doing: %v", err)
				tingUrl = orignalUrl
			}
		} else {
			cxt.Errorf("getTinyUrl: %v", err)
			tingUrl = orignalUrl
		}
	}
	ch <- tingUrl
}

func status(w http.ResponseWriter, ok bool, funcName string, q string, res string, stored bool) {
	s := fmt.Sprintf(`{"status":%s, "function":"%s", "q":"%s", "result":"%s", "stored":%s }`,
	strconv.FormatBool(ok),
	funcName,
	 q,
	res,
	strconv.FormatBool(stored))

	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}
