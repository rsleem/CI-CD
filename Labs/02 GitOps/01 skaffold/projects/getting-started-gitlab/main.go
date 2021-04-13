package main

import (
    "fmt"
	"time"
    "log"
    "net/http"
    "os"
    "net"
)

func handler(w http.ResponseWriter, r *http.Request) {
    addrs,_ := net.InterfaceAddrs()
    host,_ := os.Hostname()
    fmt.Fprintf(w, "Hello James\n")
    fmt.Fprintf(w, "%s\n", host)
    fmt.Fprintf(w, "%s", addrs[1].(*net.IPNet).IP.To4())
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":80", nil))
}