package main

import (
	"fmt"
	_ "testing"
	"net/http"
)

func main(){
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Welcome to the Dark Side!")
	})

	fs:=http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":5000", nil)
	}

