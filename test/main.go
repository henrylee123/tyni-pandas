package main

import (
	"fmt"
	"github.com/noaway/dateparse"
)

func main() {
	d, _ := dateparse.ParseLocal("2012/01/01")
	fmt.Println(d)
}
