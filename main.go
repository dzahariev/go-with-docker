package main

import (
	"dzahariev/go-with-prow/pkg01"
	"dzahariev/go-with-prow/pkg02"
	"fmt"
)

func main() {
	fmt.Println(hello())
	fmt.Println(pkg01.HelloPkg01())
	fmt.Println(pkg02.HelloPkg02())
}

func hello() string {
	return "Hello World!"
}
