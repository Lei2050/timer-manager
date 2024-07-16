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
	next    int
}
