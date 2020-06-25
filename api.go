package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	_ "testing"
)

func apiRequest(verb string, url string, data map[string]string) (string, error) {
	client := &http.Client{}
	req, err := func() (*http.Request, error) {
		if verb == http.MethodGet {
			if req, err := http.NewRequest(http.MethodGet, url, nil); err == nil {
				q := req.URL.Query()
				for k, v := range data {
					q.Set(k, v)
				}
				req.URL.RawQuery = q.Encode()
				return req, nil
			} else {
				return nil, err
			}
		} else if verb == http.MethodPost {
			if marshalData, err := json.Marshal(data); err == nil {
				res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalData))
				return res, err
			} else {
				return nil, err
			}
		} else {
			return nil, errors.New("Bad Verb!")
		}
	}()
	if err == nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/vnd.github.v3+json, application/json")
		req.Header.Add("User-Agent", "https://example-app.com")
		if token, ok := session["access_token"]; ok == true {
			req.Header.Add("Authorization", "Bearer "+token)
		}
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if data, err := ioutil.ReadAll(res.Body); err == nil {
				return string(data), nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
