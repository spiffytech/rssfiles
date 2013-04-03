package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    resp, err := http.Get("http://google.com")
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    bytes, err := ioutil.ReadAll(resp.Body)
    fmt.Println(string(bytes))
    
}
