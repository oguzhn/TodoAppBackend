package main

type Todo struct {
	Id    int    `json:"-"`
	Title string `json:"title"`
}
