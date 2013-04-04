package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    //"net/url"
    "os"

    "github.com/stretchrcom/goweb/goweb"
)

func main() {
    type login struct {
        Apikey string `json:"apikey"`
        Email string `json:"email"`
        Password string `json:"password"`
    }
    j := login{Apikey: os.Getenv("apikey"), Email: os.Getenv("email"), Password: os.Getenv("password")}

    b, err := json.Marshal(j)
    if err != nil {
        panic(err)
    }
    reader := bytes.NewReader(b)

    resp, err := http.Post("https://open.ge.tt/1/users/login", "application/json", reader)
    if err != nil {
        panic(err)
    }

    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s", bytes)


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

    //goweb.ListenAndServe("0.0.0.0:3000")
}
