package timermgr

import (
	"fmt"
	"testing"
)

type TestArrayPoolStruct struct {
	Val int
}

func TestArrayPool(t *testing.T) {
	sttAp := NewPool[TestArrayPoolStruct](0)
	id := sttAp.Alloc()
	fmt.Println(id, sttAp.Get(id))
	id1 := sttAp.Alloc()
	fmt.Println(id1, sttAp.Get(id1))
	id2 := sttAp.Alloc()
	stt := sttAp.GetRef(id2)
	fmt.Println(id2, stt)
	stt.Val = 1000
	fmt.Println(id2, stt)
	fmt.Println(id2, sttAp.Get(id2))
	fmt.Println(id2, sttAp.GetRef(id2))
	id3 := sttAp.Alloc()
	fmt.Println(id3, sttAp.Get(id3))

	fmt.Println("sttAp:", sttAp)
	sttAp.GetRef(id).Val = 1
	sttAp.GetRef(id1).Val = 2
	sttAp.GetRef(id3).Val = 4
	fmt.Println("sttAp:", sttAp)
	sttAp.Free(id)
	fmt.Println("sttAp:", sttAp)
	sttAp.Free(id3)
	fmt.Println("sttAp:", sttAp)
	sttAp.Free(id1)
	sttAp.Free(id2)
	fmt.Println("sttAp:", sttAp)

	sttAp.Alloc()
	sttAp.Alloc()
	sttAp.Alloc()
	sttAp.Alloc()
	sttAp.Alloc()
	fmt.Println("sttAp:", sttAp)
	sttAp.Alloc()
	fmt.Println("sttAp:", sttAp)
	sttAp.Alloc()
	fmt.Println("sttAp:", sttAp)
	sttAp.Alloc()
	fmt.Println("sttAp:", sttAp)

	sttAp.Free(1)
	sttAp.Free(1)
	sttAp.Free(1)
	sttAp.Free(1)
	sttAp.Free(1)
	fmt.Println("sttAp:", sttAp)
	sttAp.Free(9)
}
