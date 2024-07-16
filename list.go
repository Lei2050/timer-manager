package timermgr

type bucketEntry struct {
	nextTimerID int
}

func (be *bucketEntry) push(te *timerListEntry) {
	if be.nextTimerID != 0 {
		te.next = be.nextTimerID
	}
	be.nextTimerID = te.timerID
}

type timerListEntry struct {
	timer
	timerID int
	gen     int //当前回收次数，主要用于判断两个timerID是否是同一个
	next    int
}
