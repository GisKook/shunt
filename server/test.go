package main

import (
	"log"
	"time"
)

func main() {
	var aa string
	aa = time.Now().Format("2006-01-02 15:04:05")
	aa += "asdfasdf"
	log.Println(aa)
}
