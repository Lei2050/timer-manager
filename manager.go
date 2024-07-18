package timermgr

import (
	"container/heap"
	"fmt"
	"time"
)

type TimerManager struct {
	pq          priorityQueue
	time2Bucket map[int64]int
	timerPool   []timerListEntry

	bucketEntryPool *ArrayPool[bucketEntry]
	listEntryPool   *ArrayPool[timerListEntry]

	pendingExec []int
}

func New() *TimerManager {
	tm := &TimerManager{
		pq:          make(priorityQueue, 0, initCap),
		time2Bucket: make(map[int64]int, initCap),
		timerPool:   make([]timerListEntry, initCap),

		bucketEntryPool: NewPool[bucketEntry](initCap),
		listEntryPool:   NewPool[timerListEntry](initCap),
	}
	return tm
}

func (tm *TimerManager) getOrAllocBucket(endTime int64) (bucketId int) {
	bucketId, exist := tm.time2Bucket[endTime]
	if !exist {
		bucketId = tm.bucketEntryPool.Alloc()
		tm.bucketEntryPool.Arr[bucketId] = bucketEntry{}
		tm.time2Bucket[endTime] = bucketId
		heap.Push(&tm.pq, endTime)
	}
	return bucketId
}

func (tm *TimerManager) addTimer(d time.Duration, isRepeat bool, f TimerHandler, canceler *Canceler, a ...any) TimerID {
	if d <= 0 {
		panic("timer duration <= 0")
	}

	// nowMilliSec := time.Now().UnixMilli()
	endTime := time.Now().UnixMilli() + d.Milliseconds()
	bucketId := tm.getOrAllocBucket(endTime)

	be := &tm.bucketEntryPool.Arr[bucketId]

	idx := tm.listEntryPool.Alloc()
	te := &tm.listEntryPool.Arr[idx]
	gen := te.gen + 1
	*te = timerListEntry{
		timer: timer{
			end:      endTime,
			interval: d.Milliseconds(),
			repeat:   isRepeat,
			isCancel: false,
			callback: f,
			args:     a,
		},
		idx: idx,
		gen: gen,
	}

	// fmt.Printf("now:%d, endTime:%d, bucketId:%d, idx:%d, gen:%d, args:%+v\n", nowMilliSec, endTime, bucketId, idx, gen, a)

	be.push(te) //这里te会逃逸吗？不会

	timerID := encodeTimerID(te.idx, te.gen)
	if canceler != nil {
		canceler.register(timerID)
	}

	return timerID
}

func (tm *TimerManager) AddTimer(d time.Duration, f TimerHandler, canceler *Canceler, a ...any) TimerID {
	return tm.addTimer(d, false, f, canceler, a...)
}

func (tm *TimerManager) AddRepeatTimer(d time.Duration, f TimerHandler, canceler *Canceler, a ...any) TimerID {
	return tm.addTimer(d, true, f, canceler, a...)
}

func (tm *TimerManager) execTimer(idx int) time.Duration {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("execTimer err:%+v\n", err)
		}
	}()

	te := &tm.listEntryPool.Arr[idx]
	if te.isCancel {
		return 0
	}

	start := time.Now()
	te.do()
	return time.Since(start)
}

func (tm *TimerManager) repeatTimer(idx int) {
	te := &tm.listEntryPool.Arr[idx]
	te.end = te.end + te.interval
	bucketId := tm.getOrAllocBucket(te.end)
	be := &tm.bucketEntryPool.Arr[bucketId]
	be.push(te)
	// fmt.Printf("end:%d, interval:%d, bucketId:%d, idx:%d\n", te.end, te.interval, bucketId, idx)
}

func (tm *TimerManager) execPendingTimer() bool {
	var cumulCost time.Duration
	for i, idx := range tm.pendingExec {
		cumulCost += tm.execTimer(idx)
		if tm.listEntryPool.Arr[idx].repeat {
			tm.repeatTimer(idx)
		} else {
			tm.listEntryPool.Free(idx)
		}
		if maxExecTimePerTick > 0 && cumulCost.Milliseconds() > maxExecTimePerTick {
			tm.pendingExec = tm.pendingExec[i+1:]
			return false
		}
	}
	tm.pendingExec = tm.pendingExec[:0]
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
		timerEntryIdx := be.timerListEntryIdx
		for timerEntryIdx > 0 {
			tm.pendingExec = append(tm.pendingExec, timerEntryIdx)
			next := tm.listEntryPool.Arr[timerEntryIdx].next
			tm.listEntryPool.Arr[timerEntryIdx].next = 0 //断开链表
			timerEntryIdx = next
		}
		// fmt.Printf("free bucket:%d\n", bucketId)
		delete(tm.time2Bucket, onTime)
		tm.bucketEntryPool.Free(bucketId)

		if !tm.execPendingTimer() {
			return
		}
	}
}

func (tm *TimerManager) CancelTimer(timerID TimerID) {
	idx, gen := decodeTimerID(timerID)
	te := &tm.listEntryPool.Arr[idx]
	if gen != te.gen {
		//说明该timer之前被回收并复用了，当前Canceler持有的是个失效的
		return
	}
	if te.isCancel {
		return
	}
	te.cancel()
}

func (tm *TimerManager) NewCanceler() *Canceler {
	return &Canceler{timerMgr: tm}
}
