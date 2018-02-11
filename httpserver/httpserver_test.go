package httpserver

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

/*
SET UP FOR HASH HANDLER TESTS
*/
var ts *httptest.Server

func TestMain(m *testing.M) {
	//setup test server
	var h hashHandler
 	http.Handle("/hash", h)
  ts = httptest.NewServer(h)
  defer ts.Close()

  m.Run()
}


/*
TEST FOR MY SERVER
*/
func TestMyServerMakeServer(t *testing.T) {
  cases := []struct {
    in string;
    addr string;
    handler hashHandler;
  } {
    {"8080", "localhost:8080", hashHandler{}},
  }
    for _, c := range cases {
        makeServer(c.in)
    if s.Addr != c.addr && s.Handler != c.handler {
      t.Errorf("makeServer(%q) == {Addr : %q, Hanlder : %q}, want {Addr : %q, Handler : %q", c.in, s.Addr, s.Handler, c.addr, c.handler)
    }
  }
}

func TestMyServerRun(t *testing.T) {
  go Run("8082")
  _, err := http.Get("localhost:8082/")
  if err == nil {
    t.Errorf("Server not running")
  }
  shutdownMyServer()
}

func TestMyServerShutdown(t *testing.T) {
  go Run("8082")

  var wait sync.WaitGroup
  wait.Add(1)

  var res *http.Response
  v := url.Values{}
  v.Set("password", "angryMonkey")
  go func () { 
    defer wait.Done()
    res, _ = http.PostForm(ts.URL + "/hash", v)
  }()

  shutdownMyServer()

  wait.Wait()

  got, _ := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if string(got) != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
    t.Errorf("Shutdown not processing requests already sent")
  }

  _, err := http.Get("localhost:8082/")
  if err == nil {
    t.Errorf("Server accepting requests after shutdown")
  }
}


/*
TESTS FOR HASH HANDLER
*/
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