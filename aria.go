package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ericpfisher/aria/db"
)

var (
	domain = flag.String("domain", "short.example.com", "Custom domain")
	port   = flag.String("port", "8888", "Port to listen")
	db     = ariaDB.NewDB()
)

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", createHandler)

	log.Println("Listening on", *port)
	err := http.ListenAndServe(":"+*port, nil)

	if err != nil {
		panic(err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI, r.RemoteAddr)

	if r.RequestURI == "/" {
		form := `<html>
							<body>
								<form method="GET" action="shorten">
									URL:<br />
									<input type="text" name="url"><br />
									<input type="submit" value="Shorten!">
								</form>
							</body>
							</html>`
		fmt.Fprint(w, form)
		return
	}

	hash := strings.TrimPrefix(r.URL.RequestURI(), "/")
	lookupVal := db.GetLink(hash)

	if lookupVal == "" {
		http.NotFound(w, r)
		return
	}

	unescapedURL, unescapedError := url.QueryUnescape(lookupVal)

	if unescapedError != nil {
		log.Fatal(unescapedError)
	}

	if !strings.HasPrefix(unescapedURL, "http") {
		unescapedURL = fmt.Sprintf("http://%s", unescapedURL)
	}

	http.Redirect(w, r, unescapedURL, 301)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI, r.RemoteAddr)

	queryString := r.URL.Query()
	log.Printf("raw: %v\nescaped: %v", queryString.Get("url"), url.QueryEscape(queryString.Get("url")))
	toBeShortened := url.QueryEscape(queryString.Get("url"))

	if toBeShortened == "" {
		http.Error(w, "Bad Request", 400)
	}

	hash := db.AddLink(toBeShortened)
	fmt.Fprintf(w, "http://%s/%s", *domain, hash)
}
