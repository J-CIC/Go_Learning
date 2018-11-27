package add

import "testing"

func TestAdd(t *testing.T) {
    list := []int{1,2,3,4,5}
    if Add(list) != 16 {
        t.Log("add for ",list," fails")
        t.Fail()
    }
}