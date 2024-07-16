package timermgr

const (
	initCap            = 64
	maxExecTimePerTick = 50 //每次tick的最大运行时间(ms)，用于平滑
)

type TimerID = int
