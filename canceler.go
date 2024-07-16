package timermgr

type Canceler struct {
	timerIDs []TimerID
	timerMgr *TimerManager
}

func (c *Canceler) register(timerID TimerID) {
	c.timerIDs = append(c.timerIDs, timerID)
}

func (c *Canceler) Cancel() {
	if len(c.timerIDs) <= 0 {
		return
	}
	for _, timerID := range c.timerIDs {
		c.timerMgr.CancelTimer(timerID)
	}
	c.timerIDs = c.timerIDs[:0]
}
