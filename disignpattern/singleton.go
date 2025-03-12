package disignpattern

import (
	"fmt"
	"sync"
)

type Singleton struct {
	Foo string
}

var (
	instance *Singleton
	once     sync.Once
)

func GetSingleInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Create a new instance")
		instance = &Singleton{
			Foo: "bar",
		}
	})
	return instance
}

func SingleExample() {
	for i := 0; i < 10; i++ {
		s := GetSingleInstance()
		fmt.Println(s)
	}
}
