package httpserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
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
	res, err := http.PostForm(ts.URL+"/hash", v)
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
		in      string
		addr    string
		handler hashHandler
	}{
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
	go Run("8083")

	var wait sync.WaitGroup
	wait.Add(1)

	var res *http.Response
	v := url.Values{}
	v.Set("password", "angryMonkey")
	go func() {
		defer wait.Done()
		res, _ = http.PostForm("http://localhost:8083/hash", v)
	}()

	shutdownMyServer()

	wait.Wait()

	got := getBodyOfResponse(res)
	if string(got) != "0" {
		t.Errorf("Shutdown not processing requests already sent")
	}

	_, err := http.Get("localhost:8083/")
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

	if id1 != id2-1 {
		t.Errorf("Consequetive requests returned ids %v, %v which are not incremental", got1, got2)
	}
}

func TestHashHandlerPasswordStoredAtCorrectId(t *testing.T) {
	cases := []struct {
		in, want string
	}{
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
	}{
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
	res, err := http.Post(ts.URL+"/hash", "html", nil)
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
	res, err := http.PostForm(ts.URL+"/hash", v)
	if err != nil {
		log.Fatal(err)
	}
	got := getBodyOfResponse(res)
	if string(got) != "400 needs password query" {
		t.Errorf("Status code for post to /hash without password query %v wanted 400 need password query", string(got))
	}
}

func TestHashHandlerPostHashId(t *testing.T) {
	res, err := http.Post(ts.URL+"/hash/2", "html", nil)
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

	computeHash("angryMonkey", 5)

	duration := time.Since(start).Seconds()
	if duration < 5 || duration > 6 {
		t.Errorf("Socket lagging for %v seconds wanted around 5", duration)
	}
}

func TestGetStatistics(t *testing.T) {
	res, _ := http.Get(ts.URL + "/stats")
	js := getBodyOfResponse(res)

	expected := fmt.Sprintf("{\"total\":%v,\"average\":%v}", st.Total, st.Average)

	if string(js) != expected {
		t.Errorf("/stats returned %v wanted %v", js, expected)
	}
}

/*
TEST FOR STATISTICS
*/

func TestNewStats(t *testing.T) {
	st := NewStats()

	if st.Total != 0 || st.Average >= 0 {
		t.Errorf("Expected new stats to start with 0, -1 got %v, %v", st.Total, st.Average)
	}
}

func TestAdjustAverage(t *testing.T) {
	got := adjustAverage(3, 3, 4)

	if got != 3.25 {
		t.Errorf("adjustAverage(3,3,4) got %v wanted 3.25", got)
	}
}

func TestUpdateStatistics(t *testing.T) {
	s := Stats{3, 3}

	s.UpdateStatistics(time.Millisecond * 4)

	if s.Total != 4 || s.Average != 3.25 {
		t.Errorf("Expected new total 4 got %v, expeted new average 3.25 got %v", s.Total, s.Average)
	}
}

func TestEncode(t *testing.T) {
	s := Stats{3, 3}

	j := s.Encode()

	expected := "{\"total\":3,\"average\":3}"

	if j != expected {
		t.Errorf("Encode returned %v wanted %v", s.Total, expected)
	}
}

/*
TEST FOR IDMAP
*/

func TestGetCurrentId(t *testing.T) {
	ic := make(chan int)
	go GetCurrentId(ic)
	go GetCurrentId(ic)
	id1, id2 := <-ic, <-ic

	diff := id1 - id2

	if diff != 1 && diff != -1 {
		t.Errorf("Expected GetCurrentId(chan int) to return incremental ids got %v, %v", id1, id2)
	}
}

func TestReadFromMap(t *testing.T) {
	i := IdMap{m: make(map[string]string)}

	i.m["0"] = "hash1"
	got := i.ReadFromMap("0")

	if got != "hash1" {
		t.Errorf("ReadFromMap returned %v wanted hash1", got)
	}
}

func TestConcurrentWriteToMap(t *testing.T) {
	i := IdMap{m: make(map[string]string)}

	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()
		i.WriteToMap("hash1", 0)
	}()

	go func() {
		defer wait.Done()
		i.WriteToMap("hash2", 1)
	}()

	wait.Wait()

	if i.m["0"] != "hash1" || i.m["1"] != "hash2" {
		t.Errorf("WriteToMap(hash1, 0) writes %v expected hash1", i.ReadFromMap("0"))
	}
}
