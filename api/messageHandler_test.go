package api

import (
	"reflect"
	"testing"
)
func TestX(t *testing.T){
	str1:="hello"
	str2:="hallo"
	if reflect.DeepEqual(str1,str2){
		t.Errorf("got %v,want %v",str1,str2)
	}
}