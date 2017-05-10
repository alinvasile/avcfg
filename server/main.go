package main

import (
    "flag"
    "io/ioutil"
    "fmt"
    "net/http"
    "html"
    "time"
    "strings"
    "strconv"
	"path/filepath"
    "github.com/antonholmquist/jason"
    "github.com/patrickmn/go-cache"

)

var type_string = flag.String("type","http", "Listen type (http or tcp)")
var tcp_port_string = flag.String("tcp.port","8087", "TCP Listen port")
var http_port_string = flag.String("http.port","8080", "HTTP Listen port")
var property_location_string = flag.String("json.location","./data", "Location of json files")

var cache_ttl_string = flag.String("cache.ttl","300", "Property cache TTL in seconds")
var cache_purge_interval_string = flag.String("cache.purge.interval","30", "Property cache purge interval in seconds")



func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getPropertyIntValue(value string) (int64){
	var p,err=strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Printf("Invalid integer value: %s\n", value)
		panic(err)
	}	
	return p
}

func readPropertyThroughCache(propertyCache *cache.Cache, path string) (string, error){
	foo, found := propertyCache.Get(path)
	if found {
		return fmt.Sprintf("%s",foo), nil
	}

	components := parseComponents(path)

	ymlFile:=*property_location_string + "/" + components[1]+".json"
	filename, _ := filepath.Abs(ymlFile)

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("property file %s does not exist", filename)
	}

	v, err := jason.NewObjectFromBytes(dat)

	if err != nil {
		return "", fmt.Errorf("property file %s cannot be parsed", filename)
	}

	var children []string
	children=components[2:]

	value,err:=v.GetString(children...)
	if err != nil {
		return "", fmt.Errorf("Property %s does not exist", html.EscapeString(path))
	}

	propertyCache.Set(path, value, cache.DefaultExpiration)

	return value, nil
}





func main(){
	flag.Parse()

	var serverType=*type_string

	var port=*http_port_string
	var tcpPort=*tcp_port_string

	ttl:=getPropertyIntValue(*cache_ttl_string)
	purge_interval:=getPropertyIntValue(*cache_purge_interval_string)


	fmt.Printf("HTTP Listen port: %s\n", port)
	fmt.Printf("TCP Listen port: %s\n", tcpPort)
	fmt.Printf("JSON location: %s\n", *property_location_string)

	fmt.Printf("Cache TTL in seconds: %d\n", ttl)
	fmt.Printf("Cache purge interval in seconds: %d\n", purge_interval)

	// Create a cache with a default expiration time of 5 minutes, and which
    	// purges expired items every 30 seconds
    	propertyCache := cache.New(time.Duration(ttl) * time.Second, time.Duration(purge_interval) * time.Second)

	if(serverType == "http"){
		startHttpServer(propertyCache, port)
	} else {
		startTcpServer(propertyCache, tcpPort)
	}




}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        fmt.Fprint(w, message)
    }
}

// http://learntogoogleit.com/post/56844473263/url-path-to-array-in-golang
func parseComponents(path string) []string {

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