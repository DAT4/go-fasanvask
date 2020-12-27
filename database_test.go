package main

import (
	"fmt"
	"testing"
)

func Test_authenticate(t *testing.T) {
	fmt.Println("Testing Auth")
	user := Resident{
		Name:     "Martin",
		Room:     "2b",
		Password: "pass",
	}

	if authenticate(user) {
		println("OK")
	} else {
		println("Fail")
	}
}

func Test_getUser(t *testing.T) {
	fmt.Println("Testing GetUser")
	id := "5fe8945482cfc5c53385d50f"
	user, err := getUser(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)
}