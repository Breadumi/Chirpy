package main

type reqS struct {
	Body string `json:"body"`
}

type eS struct {
	Error string `json:"error"`
}

type cleanedBody struct {
	Cleaned_Body string `json:"cleaned_body"`
}
