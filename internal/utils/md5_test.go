package utils

import (
	"testing"
)

func TestStringToMd5(t *testing.T) {
	s := "test"
	expected := "098f6bcd4621d373cade4e832627b4f6"
	actual := StringToMd5(s)
	if actual != expected {
		t.Errorf("StringToMd5(%s) = %s; want %s", s, actual, expected)
	}

}

func TestSetToMd5(t *testing.T) {
	s := []string{"test1", "test2"}
	expected := "beff3fcba56f29677c5d52b843df365e"
	actual := SetToMd5(s)
	if actual != expected {
		t.Errorf("SetToMd5(%s) = %s; want %s", s, actual, expected)
	}

}