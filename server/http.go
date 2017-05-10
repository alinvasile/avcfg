package main

import (
	"github.com/patrickmn/go-cache"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"log"
)

func startHttpServer(propertyCache *cache.Cache, port string){
	http.HandleFunc("/service/json", func(w http.ResponseWriter, r *http.Request) {
		files, _ := ioutil.ReadDir(*property_location_string)
		for _, f := range files {
			filename:=f.Name()
			extension:= filepath.Ext(filename)
			var name = filename[0:len(filename)-len(extension)]
			fmt.Fprintf(w, "%s\n",name)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/json/", func(w http.ResponseWriter, r *http.Request) {

		fullUrl:=r.URL.Path

		value, err := readPropertyThroughCache(propertyCache, fullUrl)

		if err != nil {
			errorHandler(w,r,http.StatusNotFound, fmt.Sprintf("Error reading: %s", err.Error()))
			return
		}

		fmt.Fprintf(w, "%s", value)

	})

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
