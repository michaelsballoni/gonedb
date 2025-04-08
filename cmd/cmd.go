package main

import (
	"fmt"

	"github.com/michaelsballoni/gonedb/pkg"
)

func main() {
	n := gonedb.Node{}
	l := gonedb.Link{}

	fmt.Printf("Null Node: %d\n", n.Id)
	fmt.Printf("Null Link: %d\n", l.Id)
}
