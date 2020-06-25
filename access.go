package main

type Access struct {
	Access_token string
	Token_type   string
	Scope        string
}

type Repo struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	GHUrl    string `json:"html_url"`
}
