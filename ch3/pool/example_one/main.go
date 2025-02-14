package main

import (
	"fmt"
	"sync"
)

func main() {
	pool := &sync.Pool{
		New: func() interface{} {
			s := "New Object"

			return &s
		},
	}

	obj1 := pool.Get().(*string)
	fmt.Println("Get:", *obj1)

	r := "Reused Object"

	pool.Put(&r)

	obj2 := pool.Get().(*string)
	fmt.Println("Get:", *obj2)

	obj3 := pool.Get().(*string)
	fmt.Println("Get:", *obj3)
}

// func main() {
// 	myPool := &sync.Pool{
// 		New: func() interface{} {
// 			fmt.Println("Creating new instance")
// 			return struct{}{}
// 		},
// 	}
//
// 	myPool.Get()
// 	instance := myPool.Get()
// 	myPool.Put(instance)
// 	myPool.Get()
// }
