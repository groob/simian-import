package main

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	tag := NewTag("foobar", "groob@groob.io", "groob-io")
	tag.addKind("Computer", "my-mac.local")
	tag.addKind("Computer", "your-mac.local")
	if err := tag.marshalKeys(); err != nil {
		t.Fatal(err)
	}
	fmt.Println(tag.Keys[0])
}
