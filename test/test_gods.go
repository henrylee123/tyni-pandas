package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func main() {
	m := linkedhashmap.New() // empty (keys are of type int)
	m.Put(1, "x")            // 2->b, 1->x (insertion-order)
	m.Put(2, "b")            // 2->b
	m.Put(1, "a")            // 2->b, 1->a (insertion-order)
	it := m.Iterator()
	for it.Next() {
		fmt.Println(it.Key(), it.Value())
	}
	_, _ = m.Get(2) // b, true
	_, _ = m.Get(3) // nil, false
	_ = m.Values()  // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()    // []interface {}{2, 1} (insertion-order)
	m.Remove(1)     // 2->b
	m.Clear()       // empty
	m.Empty()       // true
	m.Size()        // 0
}
