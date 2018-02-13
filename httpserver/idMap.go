package httpserver

import (
	"strconv"
	"sync"
)

/*
Stores objects by identifiers
Ensures that the key is an string integer for easy
Enables concurrent manipulation w/ read & write locks
*/
type IdMap struct {
	sync.RWMutex
	m map[string]string
}

//Next id to return, insures uniqueness by synced incrementing
var nextID int = 0

/*
Returns a unique id by incrementing nextId, returns through an int channel
*/
func GetCurrentId(i chan int) {
	id := nextID
	nextID++
	i <- id
}

/*
Blocks until write lock is obtained
Then writes the id, value pair to the map
Converts the integer into a string for easily getting from path string
*/
func (i *IdMap) WriteToMap(value string, id int) {
	strid := strconv.Itoa(id)
	i.Lock()
	i.m[strid] = value
	i.Unlock()
}

/*
Blocks until read lock is obtained
Then reads value for given id
*/
func (i *IdMap) ReadFromMap(id string) string {
	i.RLock()
	value := i.m[id]
	i.RUnlock()
	return value
}
