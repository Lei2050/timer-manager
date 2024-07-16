package timermgr

type priorityQueue []int64

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i] < pq[j]
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	tmp := *pq
	n := len(tmp)
	tmp = tmp[0 : n+1]
	timer := x.(int64)
	tmp[n] = timer
	*pq = tmp
}

func (pq *priorityQueue) Pop() any {
	tmp := *pq
	n := len(tmp)
	timer := tmp[n-1]
	*pq = tmp[0 : n-1]
	return timer
}
