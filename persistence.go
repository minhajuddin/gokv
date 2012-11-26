package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

//loads the key value data from the persistence file
//when the server is started
func loadKv() error {
	data, err := ioutil.ReadFile(kvFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &kv)
}

//read writes are done after locking the current go routine
func getValue(key string) (value interface{}, ok bool) {
	mutex.Lock()
	defer mutex.Unlock()
	value, ok = kv[key]
	return
}

func setValue(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	kv[key] = value
}

//persists the data to the persistence file when the
//server shuts down
func persistKv() error {
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(kvFile, bytes, 0600)
}

//returns keys starting with this prefix
func listKeys(prefix string) []string {
	keys := make([]string, 0, len(kv))
	mutex.Lock()
	for k := range kv {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys)
		}
	}
	mutex.Unlock()
	return keys
}

//deletes a key from the kv store
func deleteKey(key string) {
	mutex.Lock()
	delete(kv, key)
	mutex.Unlock()
}
