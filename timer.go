package timermgr

type TimerHandler func(a ...any)

type timer struct {
	end      int64 //结束时间，millisecond
	interval int64 //间隔，millisecond
	repeat   bool  //是否重复
	isCancel bool  //是否已取消
	callback TimerHandler
	args     []any
}

func (t *timer) do() {
	if t.isCancel {
		return
	}
	t.callback(t.args...)
}

func (t *timer) cancel() {
	t.repeat = false
	t.isCancel = true
	t.args = nil
}
