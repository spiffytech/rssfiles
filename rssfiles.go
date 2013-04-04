package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    //"os"

    "github.com/stretchrcom/goweb/goweb"
)

func main() {
    goweb.MapFunc("/", func(c *goweb.Context) {
        resp, err := http.Get("http://almamater.xkcd.com/")
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        //teeReader := io.TeeReader(resp.Body, os.Stdout)
        teeReader := io.TeeReader(resp.Body, c.ResponseWriter)
        ioutil.ReadAll(teeReader)
    })

    goweb.ListenAndServe("0.0.0.0:3000")

    //bytes := make([]byte, 256)
    //for {
    //    n, err := resp.Body.Read(bytes)
    //    if err != nil && n != 0{
    //        panic(err)
    //    }
    //    fmt.Printf("%s", bytes)
    //    if n == 0 {
    //        fmt.Println("leaving loop...")
    //        break
    //    }
    //}
    fmt.Println("here")
}
