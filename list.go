package timermgr

type bucketEntry struct {
	timerListEntryIdx int
}

func (be *bucketEntry) push(te *timerListEntry) {
	if be.timerListEntryIdx != 0 {
		te.next = be.timerListEntryIdx
	}
	be.timerListEntryIdx = te.idx
}

type timerListEntry struct {
	timer
	idx  int
	gen  int //当前回收次数，主要用于判断两个timerID是否是同一个
	next int
}
