package main

import (
   "testing"
   "net/http"
   "net/http/httptest"
   "net/url"
   "io/ioutil"
   "log"
)

 func TestHashHandlerGoodRequest(t *testing.T) {
 	cases := []struct {
		in, want string
	} {
		{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
	}

 	var h hashHandler
 	http.Handle("/hash", h)
  	ts := httptest.NewServer(h)
  	defer ts.Close()
  
    for _, c := range cases {
    	v := url.Values{}
    	v.Set("password", c.in)
        res, err := http.PostForm(ts.URL + "/hash", v)
        if err != nil {
        	log.Fatal(err)
        }
        got, err := ioutil.ReadAll(res.Body)
        res.Body.Close()
        if err != nil {
        	log.Fatal(err)
        }
		if string(got) != c.want {
			t.Errorf("/hash?password=%q == %q, want %q", c.in, got, c.want)
		}
	}
}

