package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/gorilla/mux"
)

type CacheApp interface {

}

type cacheApp struct {
	bazelCache BazelCache
}

func NewCacheApp(r *mux.Router, bazelCache BazelCache) CacheApp {
	c := &cacheApp{
		bazelCache: bazelCache,
	}
	
	r.HandleFunc("/cache/ac/{sha}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sha := vars["sha"]
		key := fmt.Sprintf("ac_%s", sha)
		
		switch r.Method {
		case "GET":
			c.get(key, w, r)
		case "PUT":
			c.put(key, w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusBadRequest)
		}
	})
	
	r.HandleFunc("/cache/cas/{sha}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sha := vars["sha"]
		key := fmt.Sprintf("cas_%s", sha)
		
		switch r.Method {
		case "GET":
			c.get(key, w, r)
		case "PUT":
			c.put(key, w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusBadRequest)
		}
	})

	return c
}

func (c *cacheApp) get(key string, w http.ResponseWriter, r *http.Request) {
	data, success := c.bazelCache.Get(key)
	
	if !success {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	
	w.Write([]byte(data))
}

func (c *cacheApp) put(key string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	
	if !c.bazelCache.Set(key, string(body)) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
