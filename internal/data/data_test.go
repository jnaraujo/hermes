package data

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	d := New("GET", "key")
	if d.Decode() != "GET\nkey" {
		t.Error("Decode() failed")
	}
}

func TestEncodeFromStr(t *testing.T) {
	expected := New("SET", "user=john")

	data, err := FromStr("SET\nuser=john")
	if err != nil {
		t.Error("Decode() failed", err)
	}

	if !reflect.DeepEqual(expected, data) {
		t.Error("Data not equal")
	}
}
