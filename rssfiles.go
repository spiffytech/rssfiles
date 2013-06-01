package main

import (
    "bytes"
    "encoding/csv"
    "encoding/json"
    "fmt"
    //"io"
    "io"
    "io/ioutil"
    "path/filepath"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"

    mux "github.com/gorilla/mux"
    "github.com/hoisie/mustache"
)

var baseurl = "https://open.ge.tt/1"
var mimetypes map[string]string

type File struct {
    FileID int
    FileIDRaw string `json:"fileid"`
    Filename string
    Getturl string
    GetturlRaw string `json:"getturl"`
    Created int64
    Timestamp time.Time
    TimestampStr string
    Title string
    Size int64
    Guid string
    Mimetype string
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
    //authStuff = gettLogin()
    //shares := enumerateShares()
    //fmt.Println(renderRSS(shares))

    mimetypes = make(map[string]string)

    f, err := os.Open("mimetypes")
    if err != nil {
        panic(err)
    }
    csvReader := csv.NewReader(f)
    csvReader.TrimLeadingSpace = true
    csvReader.Comma = '\t'
    csvReader.Comment = '#'
    for {
        fields, err := csvReader.Read()
        if err != nil {
            if err == io.EOF {
                break
            } else {
                panic(err)
            }
        }

        for _, extension := range strings.Fields(fields[1]) {
            mimetypes["." + extension] = fields[0]
        }
    }

    runServer()
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
    fmt.Println(string(body))
    json.Unmarshal(body, &shares)
    fmt.Println(shares)

    for shareIndex, share := range shares {
        shares[shareIndex].Guid = share.Sharename

        for fileIndex, file := range share.Files {
            shares[shareIndex].Files[fileIndex].Guid = share.Sharename + "_" + strconv.Itoa(file.FileID)

            stuff := fmt.Sprintf("http://ge.tt/api/1/files/%s/%s/blob?download", share.Sharename, file.FileIDRaw)

            // http://ge.tt/api/1/files/79HHOBi/v/4/blob?download
            // http://ge.tt/api/1/files/79HHOBi/4/blob?download
            shares[shareIndex].Files[fileIndex].Getturl = stuff

            shares[shareIndex].Files[fileIndex].Timestamp = time.Unix(file.Created, 0)
            shares[shareIndex].Files[fileIndex].TimestampStr = shares[shareIndex].Files[fileIndex].Timestamp.Format(time.RFC1123Z)

            ext := filepath.Ext(file.Filename)
            shares[shareIndex].Files[fileIndex].Mimetype = mimetypes[ext]
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


func runServer() {
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        authStuff = gettLogin()
        shares := enumerateShares()
        rss := renderRSS(shares)
        fmt.Fprint(w, rss)
    })
    http.Handle("/", r)

    fmt.Println("Listening for connections")
    http.ListenAndServe(":3001", nil)
}
