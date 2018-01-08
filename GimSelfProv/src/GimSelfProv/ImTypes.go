package main

type Inf struct {
	Uri string `json:"uri"`
}

type InfList struct{
	UriList []Inf `json:"uri-list"`
}
