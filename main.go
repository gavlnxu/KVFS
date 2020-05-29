package main

import (
	"fmt"
	"github.com/gavlnxu/KVFS/core"
)

type T struct {
	A int
	B string
}

func main() {
	// create a folder under folder test
	f, err := core.Open("test", false)
	if err != nil {
		panic(f)
	}
	// allow storage of any JSON
	f.Set("Hello", "World")
	f.Set("PI", 3.1415926)
	f.Set("test", &T{1, "OK"})
	fmt.Println("Hello:", f.Get("Hello"))
	fmt.Println("dummy will nil:", f.Get("dummy"))

	f.Increment(func(k string, v interface{}) {
		fmt.Println(k, v)
	})

	f.Del("Hello")
	// empty the directory, will remove directory
	f.Cls()
}
