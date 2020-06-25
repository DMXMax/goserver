package main

import (
	"encoding/hex"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

func processAuthCall(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "processAuthCall")
	delete(session, "access_token")
	var barray [16]byte
	if _, err := rand.Read(barray[:]); err != nil {
		log.Fatal(err)
	} else {
		session["state"] = hex.EncodeToString(barray[:])
	}
	if urlrd, err := url.Parse(authorizeURL); err == nil {
		q := urlrd.Query()
		q.Set("response_type", "code")
		q.Set("client_id", githubClientID)
		q.Set("redirect_uri", baseURL+"/auth")
		q.Set("scope", "user_public_repo")
		state, _ := session["state"]
		q.Set("state", state)
		urlrd.RawQuery = q.Encode()
		http.Redirect(w, r, urlrd.String(), http.StatusSeeOther)
	}

}

func showIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err == nil {
		_ = tmpl.Execute(w, session)
	} else {
		log.Fatal(err)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Inside handleAuth")
	returnedState := r.URL.Query().Get("state")
	if returnedState == session["state"] {
		token, _ := apiRequest(http.MethodPost, tokenURL,
			map[string]string{"grant": "authorization_code",
				"client_id":     githubClientID,
				"client_secret": githubClientSecret,
				"redirect_uri":  baseURL,
				"code":          r.URL.Query().Get("code")})
		var access Access
		json.Unmarshal([]byte(token), &access)
		if access.Access_token == "" {
			delete(session, "access_token")
		} else {
			session["access_token"] = access.Access_token
		}

		http.Redirect(w, r, baseURL, http.StatusSeeOther)
	} else {
		//redirect
		http.Redirect(w, r, baseURL+"?error=invalid_state", http.StatusSeeOther)

	}
}
func handleRepoRequest(w http.ResponseWriter, r *http.Request) {
	log.Output(1, "Inside handleRepoRequest")
	repos, err := apiRequest(http.MethodGet, apiURLBase+"user/repos",
		map[string]string{"sort": "created", "direction": "desc"})
	if err == nil {
		var arr_repos []Repo
		json.Unmarshal([]byte(repos), &arr_repos)
		tmpl, err := template.ParseFiles("repos.html")
		if err == nil {
			_ = tmpl.Execute(w, arr_repos)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Output(0, "Cannot get repos: "+err.Error())
	}
}
