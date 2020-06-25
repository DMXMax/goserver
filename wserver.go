package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	_ "testing"

	"golang.org/x/net/publicsuffix"
)

var githubClientID = "7aadf605ead5a0fda1a5"
var githubClientSecret = ""
var authorizeURL = "https://github.com/login/oauth/authorize"
var tokenURL = "https://github.com/login/oauth/access_token"
var baseURL = "http://34.105.54.95:8000"
var apiURLBase = "https://api.github.com/"
var session map[string]string
var m_jar http.CookieJar

func main() {
	session = make(map[string]string)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	m_jar = jar
	if err != nil {
		log.Fatal(err)
	}

	if secret, ok := os.LookupEnv("GITHUB_SECRET"); ok == true {
		githubClientSecret = secret
	} else {
		log.Fatal("GITHUB_SECRET not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		qry := r.URL.Query()
		action, ok := qry["action"]
		if ok == true {
			if action[0] == "login" {
				if _, err := url.Parse(authorizeURL); err != nil {
					log.Fatal(err)
				} else {
					processAuthCall(w, r)
				}
			} else {
				if _, ok := session["access_token"]; ok == false {
					http.Redirect(w, r, baseURL, http.StatusSeeOther)
				} else {
					switch action[0] {
					case "logout":
						{
							delete(session, "access_token")
							http.Redirect(w, r, baseURL, http.StatusSeeOther)
						}
					case "repos":
						{
							handleRepoRequest(w, r)
						}
					default:
						{
							http.Redirect(w, r, baseURL, http.StatusSeeOther)
						}
					}
				}
			}
		} else {
			showIndexPage(w, r)
		}
	})
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		handleAuth(w, r)
	})
	log.Output(0, "Starting Web Server")
	if err := http.ListenAndServeTLS(":8000", "/home/glen_clarkson_gmail_com/cert.pem", "/home/glen_clarkson_gmail_com/key.pem", nil); err != nil {
		log.Fatal(err)
	}
}
