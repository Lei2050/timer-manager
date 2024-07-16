package timermgr

const (
	initCap            = 64
	maxExecTimePerTick = 50 //每次tick的最大运行时间(ms)，用于平滑
)

type TimerID = uint64

func encodeTimerID(idx, gen int) TimerID {
	return uint64(gen)<<32 | uint64(idx)
}

func decodeTimerID(uid TimerID) (idx, gen int) {
	return int(uid), int(uid >> 32)
}
