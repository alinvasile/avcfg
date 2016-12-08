package main

import (
    "flag"    
    "io/ioutil"
    "fmt"
    "log"
    "net/http"
    "html"
    "time"
    "strings"    
    "path/filepath"
    "github.com/antonholmquist/jason"
    "github.com/patrickmn/go-cache"    
)


var http_port_string = flag.String("http.port","8080", "HTTP Listen port")
var property_location_string = flag.String("json.location","./data", "Location of json files")

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main(){
	flag.Parse()

	var port=*http_port_string

	fmt.Printf("HTTP Listen port: %s\n", port)
	fmt.Printf("JSON location: %s\n", *property_location_string)

	// Create a cache with a default expiration time of 5 minutes, and which
    // purges expired items every 30 seconds
    propertyCache := cache.New(5*time.Minute, 30*time.Second)

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

		foo, found := propertyCache.Get(fullUrl)
        if found {
            //fmt.Fprintf(w, "Value (cached): %s\n", foo)
            fmt.Fprintf(w, "%s", foo)
            return
        }		

		components:=parseComponents(r)      
       
		ymlFile:=*property_location_string + "/" + components[1]+".json"     
		filename, _ := filepath.Abs(ymlFile)

        dat, err := ioutil.ReadFile(filename)
		if err != nil {
			errorHandler(w,r,http.StatusNotFound, filename + " property file does not exist")
			return
		}  

		v, err := jason.NewObjectFromBytes(dat)		
		
		if err != nil {
			errorHandler(w,r,http.StatusNotFound, filename + " property file cannot be parsed")
			return
		}   

        var children []string
        children=components[2:]        
        //fmt.Printf("Looking for property: %v\n", children)

		value,err:=v.GetString(children...)
		if err != nil {
			//log.Fatal(err)
			//panic(err)
			errorHandler(w,r,http.StatusNotFound, html.EscapeString(r.URL.Path) + " does not exist")
		    return
		}

		

		//fmt.Fprintf(w, "Requested json file: %s\n", components[1])
		//fmt.Fprintf(w, "Requested property: %v\n", children)
		//fmt.Fprintf(w, "Requested property: %q\n", fullUrl)

        //fmt.Fprintf(w, "Value (file): %s\n", value)
        fmt.Fprintf(w, "%s", value)

        propertyCache.Set(fullUrl, value, cache.DefaultExpiration)
		
	})

	log.Fatal(http.ListenAndServe(":" + port, nil))
	

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        fmt.Fprint(w, message)
    }
}

// http://learntogoogleit.com/post/56844473263/url-path-to-array-in-golang
func parseComponents(r *http.Request) []string {    
   
    //The URL that the user queried.
    path := r.URL.Path
    path = strings.TrimSpace(path)
    
    //Cut off the leading and trailing forward slashes, if they exist.
    //This cuts off the leading forward slash.
    if strings.HasPrefix(path, "/") {
        path = path[1:]
    }
    //This cuts off the trailing forward slash.
    if strings.HasSuffix(path, "/") {
        cut_off_last_char_len := len(path) - 1
        path = path[:cut_off_last_char_len]
    }
    //We need to isolate the individual components of the path.
    components:=strings.Split(path, "/")
    return components
}