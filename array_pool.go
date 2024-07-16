package timermgr

type ArrayPool[T any] struct {
	Arr  []T
	size int
	free []int
}

func NewPool[T any](cap int) *ArrayPool[T] {
	return &ArrayPool[T]{
		Arr:  make([]T, cap),
		free: make([]int, cap),
	}
}

func (ap *ArrayPool[T]) grow() {
	newArray := make([]T, len(ap.Arr)*2)
	copy(newArray, ap.Arr)
	ap.Arr = newArray
}

func (ap *ArrayPool[T]) Alloc() int {
	if ap.size < len(ap.Arr) {
		ap.size++
		return ap.size
	}

	if len(ap.free) > 0 {
		res := ap.free[len(ap.free)-1]
		ap.free = ap.free[:len(ap.free)-1]
		return res
	}

	ap.grow()
	return ap.Alloc()
}

func (ap *ArrayPool[T]) Free(id int) {
	if id < 0 {
		return
	}

	if id == ap.size-1 {
		ap.size--
		return
	}

	ap.free = append(ap.free, id)
}
