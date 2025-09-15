package main

import (
	"fmt"
)

var db = map[string]string{}

func shorten(url string) string {
	key := fmt.Sprintf("short%d", len(db)+1)
	db[key] = url
	return key
}

func resolve(key string) string {
	return db[key]
}

func main() {
	k := shorten("https://example.com")
	fmt.Println("Short key:", k)
	fmt.Println("Original URL:", resolve(k))
}
