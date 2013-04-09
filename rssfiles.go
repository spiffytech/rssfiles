package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "time"

    "github.com/stretchrcom/goweb/goweb"
    "github.com/hoisie/mustache"
)

var baseurl = "https://open.ge.tt/1"

type File struct {
    Fileid int
    Filename string
    Getturl string
    Created int64
    Timestamp time.Time
    Title string
    Size int64
    Guid string
    FileID int `json:"fileid"`
}
func (f File) renderTitle(s Share) string {
    title := f.Filename

    if s.Title != "" {
        title = s.Title + "-" + title
    }
    return title
}

type Share struct {
    Sharename string
    Title string
    Created int64
    Files []File
    Guid string
}
type AuthStuff struct {
    Accesstoken string
    Refreshtoken string
}
var authStuff AuthStuff


func main() {
    authStuff = gettLogin()
    shares := enumerateShares()
    fmt.Println(renderRSS(shares))
}

func gettLogin() AuthStuff {
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

    resp, err := http.Post(baseurl + "/users/login", "application/json", reader)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var authStuff AuthStuff
    json.Unmarshal(bytes, &authStuff)

    return authStuff;
}


func enumerateShares() []Share {
    u, err := url.Parse(baseurl)
    if err != nil {
        panic(err)
    }
    u.Path += "/shares"
    q := u.Query()
    q.Set("accesstoken", authStuff.Accesstoken)
    u.RawQuery = q.Encode()

    resp, err := http.Get(u.String())
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(body))

    var shares []Share
    json.Unmarshal(body, &shares)
    fmt.Println(shares)

    for shareIndex, share := range shares {
        shares[shareIndex].Guid = share.Sharename

        for fileIndex, file := range share.Files {
            shares[shareIndex].Files[fileIndex].Guid = share.Sharename + "_" + strconv.Itoa(file.FileID)

            shares[shareIndex].Files[fileIndex].Getturl += "/blob?download"

            shares[shareIndex].Files[fileIndex].Timestamp = time.Unix(file.Created, 0)
        }
    }

    return shares
}

func renderRSS(shares []Share) string {
    tpl, err := mustache.ParseFile("rss.mustache")
    if err != nil {
        panic(err)
    }

    context := make(map[string][]Share)
    context["shares"] = shares

    return tpl.Render(context)
    //return tpl.Render(shares)

}

func runserver() {

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
