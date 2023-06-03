package main

import (
	"fmt"
	"testing"

	"time"
)

type Animal interface {
	Kind() string
}

type Dog struct {
	BornTime time.Time
}

func (d *Dog) Kind() string { 
	return fmt.Sprintf("dog: %v", d.BornTime)
}

type People struct {
	Pet Animal
}

func TestInterfaceNil(t *testing.T) {
	p := new(People)
	var d *Dog = nil
	p.Pet = d

	if p.Pet == nil {
		t.Log(1)
	} else {
		t.Log(2)
		t.Log(p.Pet.Kind())
	}
}
