package main_test

import "testing"

func TestAdd(t *testing.T){
	if 1 + 1 != 2 {
		t.Fail()
	}
}