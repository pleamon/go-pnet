package main

import (
	"log"
)

func main() {
	log.Println(8000 % 4096)
	length := 8000 + (4096 - 8000%4096)
	log.Println(length)
}
