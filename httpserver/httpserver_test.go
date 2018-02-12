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
   "strconv"
   "fmt"
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
HELPER FUNCTIONS FOR REPEATED CODE
*/
func postToHash(password string) *http.Response {
  v := url.Values{}
  v.Set("password", password)
  res, err := http.PostForm(ts.URL + "/hash", v)
  if err != nil {
    log.Fatal(err)
  }
  return res
}


func getBodyOfResponse(res *http.Response) []byte {
  got, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  return got
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
    if s.Addr != c.addr {
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
  go func () { 
    defer wait.Done()
    res = postToHash("angryMonkey")
  }()

  shutdownMyServer()

  wait.Wait()

  got := getBodyOfResponse(res)
  if string(got) != "0" {
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
 func TestHashHandlerIdIncrements(t *testing.T) {
  res1 := postToHash("first")
  got1 := getBodyOfResponse(res1)
  id1, _ := strconv.ParseInt(string(got1), 10, 16)

  res2 := postToHash("second")
  got2 := getBodyOfResponse(res2)
  id2, _ := strconv.ParseInt(string(got2), 10, 16)

	if id1 != id2 - 1 {
		t.Errorf("Consequetive requests returned ids %v, %v which are not incremental", got1, got2)
	}
}

func TestHashHandlerPasswordStoredAtCorrectId(t *testing.T) {
  cases := []struct {
    in, want string
  } {
    {"thisPassword", "qoNIxVxpNORI0PURYPpzz34mCogGX7LcHopAADCdq/E7ywCJ8kou1dhw/HM2p0qfuQv9FDIa6VFl1RaOxwExSw=="},
  }
  for _, c := range cases {
    res := postToHash(c.in)
    got := getBodyOfResponse(res)
    time.Sleep(time.Second * 6)

    hashes.RLock()
    hash := hashes.m[string(got)]
    hashes.RUnlock()
    
    if hash != c.want {
      t.Errorf("/hash?password=%q == %q, want %q", c.in, hashes.m[string(got)], c.want)
    }
  }
}

func TestHashHandlerPasswordGetId(t *testing.T) {
  cases := []struct {
    in, want string
  } {
    {"thisPassword", "qoNIxVxpNORI0PURYPpzz34mCogGX7LcHopAADCdq/E7ywCJ8kou1dhw/HM2p0qfuQv9FDIa6VFl1RaOxwExSw=="},
  }
  for _, c := range cases {
    resPost := postToHash(c.in)
    got := getBodyOfResponse(resPost)
    time.Sleep(time.Second * 6)


    resGet, _ := http.Get(ts.URL + "/hash/" + string(got))
    hash := getBodyOfResponse(resGet)
    
    if string(hash) != c.want {
      t.Errorf("/hash?password=%q == %q, want %q", c.in, hashes.m[string(got)], c.want)
    }
  }
}


 func TestHashHandlerNoForm(t *testing.T) {
    res, err := http.Post(ts.URL + "/hash", "html", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.StatusCode != http.StatusBadRequest {
      t.Errorf("Status code for post to /hash without form %v wanted 400 needs form", res.Status)
    }
}

func TestHashHandlerNoPasswordQuery(t *testing.T) {
 	v := url.Values{}
  v.Set("other", "angryMonkey")
  res, err := http.PostForm(ts.URL + "/hash", v)
  if err != nil {
   	log.Fatal(err)
  }
  got := getBodyOfResponse(res)
  if string(got) != "400 needs password query" {
  	t.Errorf("Status code for post to /hash without password query %v wanted 400 need password query", string(got))
  }
}

func TestHashHandlerPostHashId(t *testing.T) {
    res, err := http.Post(ts.URL + "/hash/2", "html", nil)
    if err != nil {
        log.Fatal(err)
    }
    got := getBodyOfResponse(res)
    if string(got) != "405 method not allowed" {
      t.Errorf("Status code for post to /hash/{:id} without form %v wanted 405 method not allowed", string(got))
    }
}

func TestHashHandlerGetNoId(t *testing.T) {
  res, err := http.Get(ts.URL + "/hash")
  if err != nil {
   	log.Fatal(err)
  }
  got := getBodyOfResponse(res)
  if string(got) != "405 method not allowed" {
   	t.Errorf("Status code for get on /hash %v wanted 405 method not allowed", string(got))
  }
}

func TestHashHandlerGetNonIntegerId(t *testing.T) {
  res, err := http.Get(ts.URL + "/hash/hey")
  if err != nil {
    log.Fatal(err)
  }
  got := getBodyOfResponse(res)
  if string(got) != "400 invalid path" {
    t.Errorf("Status code for non integer id %v wanted 400 invalid path", string(got))
  }
}

func TestHashHandlerPasswordGetInvalidId(t *testing.T) {
  res, _ := http.Get(ts.URL + "/hash/500")
  got := getBodyOfResponse(res)  
  if string(got) != "404 hash not found" {
    t.Errorf("Status code for invalid id number %v wanted 404 hash not found", string(got))
  }
}

func TestHashHandlerPasswordNotHashedFor5Seconds(t *testing.T) {
	start := time.Now()

  computeHash("angryMonkey", "5")

  duration := time.Since(start).Seconds()
	if duration < 5 || duration > 6 {
		t.Errorf("Socket lagging for %v seconds wanted around 5", duration)
	} 
}

/*
TEST FOR STATISTICS
*/
func TestUpdateStatistics(t *testing.T) {
  tot := stats.Total
  a := stats.Average

  updateStatistics(time.Millisecond)

  new_avg := (tot * a + 1)/(tot + 1) 

  if stats.Total != tot + 1 || stats.Average != new_avg {
    t.Errorf("Expected new total %v got %v, expeted new average %v got %v", tot + 1, stats.Total, new_avg, stats.Average)
  }
}

func TestGetStatistics(t *testing.T) {
  res, _ := http.Get(ts.URL + "/stats")
  js := getBodyOfResponse(res)

  expected := fmt.Sprintf("{\"total\":%v,\"average\":%v}", stats.Total, stats.Average)
  fmt.Println(expected)
  
  if string(js) != expected {
    t.Errorf("/stats returned %v wanted %v", js, expected)
  }
}
