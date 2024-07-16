package timermgr

type Canceler struct {
	timerUids []uint64
	timerMgr  *TimerManager
}

func (c *Canceler) register(timerID TimerID, gen int) {
	c.timerUids = append(c.timerUids, uint64(gen)<<32|uint64(timerID))
}

func (c *Canceler) Cancel() {
	if len(c.timerUids) <= 0 {
		return
	}
	for _, uid := range c.timerUids {
		timerID, gen := int(uid>>32), TimerID(uid)
		te := &c.timerMgr.listEntryPool.Arr[timerID]
		if gen != te.gen {
			//说明该timer之前被回收并复用了，当前Canceler持有的是个失效的
			continue
		}
		c.timerMgr.CancelTimer(timerID)
	}
	c.timerUids = c.timerUids[:0]
}
