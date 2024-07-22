package timermgr

import "fmt"

type ArrayPool[T any] struct {
	arr   []T //Arr[0]是哨兵（sentinel），不会分配出去
	alloc int //下一次分配哪个
	// free  []int
	free map[int]struct{}
}

func NewPool[T any](cap int) *ArrayPool[T] {
	if cap < 0 {
		panic("cap is less than zero")
	}
	if cap == 0 {
		cap = 1
	}
	cap++
	return &ArrayPool[T]{
		arr:   make([]T, cap),
		alloc: 1,
		free:  make(map[int]struct{}),
	}
}

func (ap *ArrayPool[T]) grow() {
	newArray := make([]T, len(ap.arr)*2) //TODO
	copy(newArray, ap.arr)
	ap.arr = newArray
}

// return >=1
func (ap *ArrayPool[T]) Alloc() int {
	if ap.alloc < len(ap.arr) {
		id := ap.alloc
		ap.alloc++
		return id
	}

	// if len(ap.free) > 0 {
	// 	res := ap.free[len(ap.free)-1]
	// 	ap.free = ap.free[:len(ap.free)-1]
	// 	return res
	// }
	if len(ap.free) > 0 {
		for k := range ap.free {
			delete(ap.free, k)
			return k
		}
	}

	ap.grow()

	return ap.Alloc()
}

func (ap *ArrayPool[T]) Free(id int) {
	if id <= 0 || id >= ap.alloc {
		panic(fmt.Errorf("free invalid id:%d, next alloc pos:%d", id, ap.alloc))
	}

	ap.arr[id] = ap.arr[0] //重置为零值，防止内存泄露

	if id == ap.alloc-1 {
		ap.alloc--
		return
	}

	// ap.free = append(ap.free, id)

	_, ok := ap.free[id]
	if ok {
		return
	}
	ap.free[id] = struct{}{}
}

func (ap *ArrayPool[T]) Get(id int) T {
	return ap.arr[id]
}

func (ap *ArrayPool[T]) GetRef(id int) *T {
	return &ap.arr[id]
}
