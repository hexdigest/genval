package main

import "testing"

func TestMain(t *testing.T) {

	mainLogic("examples/simple", "validators.go", true)
	mainLogic("examples/complicated", "validators.go", true)
	mainLogic("examples/complicated_without_check", "validators.go", false)
	mainLogic("examples/overriding", "validators.go", true)
}
