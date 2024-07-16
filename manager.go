package timermgr

import (
	"container/heap"
	"time"
)

type TimerManager struct {
	pq          priorityQueue
	time2Bucket map[int64]int
	timerPool   []timerListEntry

	bucketEntryPool *ArrayPool[bucketEntry]
	listEntryPool   *ArrayPool[timerListEntry]

	pendingExec []TimerID
}

func New() *TimerManager {
	tm := &TimerManager{
		pq:          make(priorityQueue, initCap),
		time2Bucket: make(map[int64]int, initCap),
		timerPool:   make([]timerListEntry, initCap),

		bucketEntryPool: NewPool[bucketEntry](initCap),
		listEntryPool:   NewPool[timerListEntry](initCap),
	}
	tm.listEntryPool.Alloc() //idx=0的元素不用，一般地id=0表示无效
	return tm
}

func (tm *TimerManager) getOrAllocBucket(endTime int64) (bucketId int) {
	bucketId, exist := tm.time2Bucket[endTime]
	if !exist {
		bucketId = tm.bucketEntryPool.Alloc()
		tm.time2Bucket[endTime] = bucketId
		heap.Push(&tm.pq, endTime)
	}
	return bucketId
}

func (tm *TimerManager) addTimer(d time.Duration, isRepeat bool, f TimerHandler, a ...any) TimerID {
	if d <= 0 {
		panic("timer duration <= 0")
	}

	endTime := time.Now().UnixMilli() + d.Milliseconds()
	bucketId := tm.getOrAllocBucket(endTime)

	be := &tm.bucketEntryPool.Arr[bucketId]

	timerID := tm.listEntryPool.Alloc()
	te := &tm.listEntryPool.Arr[timerID]
	*te = timerListEntry{
		timer: timer{
			end:      endTime,
			interval: d.Milliseconds(),
			repeat:   isRepeat,
			callback: f,
			args:     a,
		},
		timerID: timerID,
	}

	be.push(te) //这里te会逃逸吗？

	return timerID
}

func (tm *TimerManager) AddTimer(d time.Duration, f TimerHandler, a ...any) TimerID {
	return tm.addTimer(d, false, f, a...)
}

func (tm *TimerManager) AddRepeatTimer(d time.Duration, f TimerHandler, a ...any) TimerID {
	return tm.addTimer(d, true, f, a...)
}

func (tm *TimerManager) execTimer(timerID TimerID) time.Duration {
	te := &tm.listEntryPool.Arr[timerID]
	if te.isCancel {
		return 0
	}

	start := time.Now()
	te.do()
	return time.Since(start)
}

func (tm *TimerManager) repeatTimer(timerID TimerID) {
	te := &tm.listEntryPool.Arr[timerID]
	bucketId := tm.getOrAllocBucket(te.end + te.interval)
	be := &tm.bucketEntryPool.Arr[bucketId]
	be.push(te)
}

func (tm *TimerManager) execPendingTimer() bool {
	var cumulCost time.Duration
	for _, timerID := range tm.pendingExec {
		cumulCost += tm.execTimer(timerID)
		if tm.listEntryPool.Arr[timerID].repeat {
			tm.repeatTimer(timerID)
		} else {
			tm.listEntryPool.Free(timerID)
		}
		if maxExecTimePerTick > 0 && cumulCost.Milliseconds() > maxExecTimePerTick {
			return false
		}
	}
	return true
}

func (tm *TimerManager) Tick(now time.Time) {
	if tm.pq.Len() <= 0 {
		return
	}

	if !tm.execPendingTimer() { //之前未运行完的任务还无法执行完
		return
	}

	for tm.pq.Len() > 0 {
		headTime := tm.pq[0]
		nowt := now.UnixMilli()
		if headTime > nowt {
			break
		}

		onTime := heap.Pop(&tm.pq).(int64)
		bucketId := tm.time2Bucket[onTime]
		be := &tm.bucketEntryPool.Arr[bucketId]
		curTimerID := be.nextTimerID
		for curTimerID > 0 {
			tm.pendingExec = append(tm.pendingExec, curTimerID)
		}
		tm.bucketEntryPool.Free(bucketId)

		if !tm.execPendingTimer() {
			return
		}
	}
}

func (tm *TimerManager) CancelTimer(timerID TimerID) {
	te := &tm.listEntryPool.Arr[timerID]
	if te.isCancel {
		return
	}
	te.cancel()
}
