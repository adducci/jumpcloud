package main

import (
   "testing"
   "net/http"
   "net/http/httptest"
   "net/url"
   "io/ioutil"
   "log"
   "time"
   "sync"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	//setup test server
	var h hashHandler
 	http.Handle("/hash", h)
  	ts = httptest.NewServer(h)
  	defer ts.Close()

  	m.Run()
}

 func TestHashHandlerReturnsCorrectHash(t *testing.T) {
 	cases := []struct {
		in, want string
	} {
		{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
	}
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

 func TestHashHandlerNoForm(t *testing.T) {
    res, err := http.Post(ts.URL + "/hash", "html", nil)
    if err != nil {
       	log.Fatal(err)
    }
    if res.StatusCode != http.StatusBadRequest {
    	t.Errorf("Should not allow posts without a form")
    }
}

 func TestHashHandlerNoPasswordQuery(t *testing.T) {
 	v := url.Values{}
    v.Set("other", "angryMonkey")
    res, err := http.PostForm(ts.URL + "/hash", v)
    if err != nil {
       	log.Fatal(err)
    }
    if res.StatusCode != http.StatusBadRequest {
    	t.Errorf("Should not allow posts without a password query")
    }
}

func TestHashHandlerGet(t *testing.T) {
    res, err := http.Get(ts.URL + "/hash")
    if err != nil {
       	log.Fatal(err)
    }
    if res.StatusCode != http.StatusMethodNotAllowed {
    	t.Errorf("Get is not an allowed method")
    }
}

 func TestHashHandlerLeavesSocketOpen(t *testing.T) {
 	start := time.Now()

	v := url.Values{}
    v.Set("password", "angryMonkey")

    http.PostForm(ts.URL + "/hash", v)

    duration := time.Since(start).Seconds()
	if duration < 5 || duration > 6 {
		t.Errorf("Socket not lagging for 5 seonds")
	} 
}

func TestHashHandlerConcurrentRequests(t *testing.T) {
	start := time.Now()

	var wait sync.WaitGroup
	wait.Add(2)

	v := url.Values{}
    v.Set("password", "angryMonkey")

    for i := 0; i < 2; i++ {
        go func () { 
            defer wait.Done()
            http.PostForm(ts.URL + "/hash", v)
        }()
    }

    wait.Wait()

    duration := time.Since(start).Seconds()
	if duration > 6 {
		t.Errorf("Not handling requests concurrently")
	} 
}



