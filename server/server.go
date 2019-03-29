package main
import (
    "fmt"
    "html"
    "net/http"
)

func doRedirect(w http.ResponseWriter, r *http.Request) {
    //println("Doing doRedirect from ", r.URL.Path)
    //println("query: ", r.URL.RawQuery)
    for k, v := range r.Header { 
        fmt.Printf("key[%s] value[%s]\n", k, v)
    }
    urlTo := r.URL
    urlTo.Scheme = "https"
    urlTo.Host = "solace-pubsub-imb.bosh-lite.com"
    
    //var redirectTo  = "https://solace-pubsub-imb.bosh-lite.com" + r.URL.Path
    var redirectTo  = urlTo.String()
    println("redirectTo: ", redirectTo)
    //http.Redirect(w, r, "http://localhost:8082/endpoint", 301)
    http.Redirect(w, r, redirectTo, 301)
}

func endpoint(w http.ResponseWriter, r *http.Request) {
    println("Doing endpoint from ", r.URL.Path)
    for k, v := range r.Header { 
        fmt.Printf("key[%s] value[%s]\n", k, v)
    }
    fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
    // http.Redirect(w, r, "https://solace-pubsub-imb.bosh-lite.com/SEMP/v2/config/msgVpns", 301)
}

func main() {
    http.HandleFunc("/", doRedirect)
    //http.HandleFunc("/SEMP/v2/config/msgVpns/default/restDeliveryPoints/testRdp", doRedirect)
    //http.HandleFunc("/endpoint", endpoint)
    if err := http.ListenAndServe(":8081", nil); err != nil {
        panic(err)
    }
}
