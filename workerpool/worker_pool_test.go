package workerpool

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	a := 0
	e := &worker{
		run: func() { a += 1 },
	}
	e.Go(nil)
	time.Sleep(time.Second)
	assert.Equal(t, 1, a)
}

func TestWorkerWithRelease(t *testing.T) {
	a := 0
	e := &worker{
		run: func() { a += 1 },
	}
	e.Go(func() { a += 1 })
	time.Sleep(time.Second)
	assert.Equal(t, 2, a)
}

func TestWorkerPanic(t *testing.T) {
	a := 0
	e := &worker{
		run: func() { panic("panic test.") },
	}
	e.Go(func() { a = 5 })
	time.Sleep(time.Second)
	assert.Equal(t, 5, a)
}

func TestPoolSchedule(t *testing.T) {
	a := 0
	start := time.Now()
	p := NewWorkerPool(3)

	err0 := p.Run(nil, nil)
	assert.Equal(t, ErrParamInvalid, err0)

	ctx := context.Background()
	err1 := p.Run(ctx, func() { time.Sleep(time.Second); a += 1 })
	err2 := p.Run(ctx, func() { time.Sleep(time.Second); a += 1 })
	err3 := p.Run(ctx, func() { time.Sleep(time.Second); a += 1 })
	err4 := p.Run(ctx, func() { time.Sleep(time.Second); a += 1 })
	err5 := p.Run(ctx, func() { time.Sleep(time.Second); a += 1 })

	// 资源池只有3个worker，只有前3个任务运行完（1秒过后），后两个worker才能分配到资源
	assert.True(t, time.Since(start) > time.Second)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Nil(t, err4)
	assert.Nil(t, err5)

	// wait 10ms确保前3个worker已经运行完，后两个worker还在运行中
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 3, a)
	assert.Equal(t, 1, len(p.tickets))
	assert.Equal(t, 0, len(p.running))

	// 再过1秒，来检查后两个worker是否也运行完，以及资源是否已经全部释放
	tick := time.NewTicker(time.Second + time.Millisecond*100)
	<-tick.C
	assert.Equal(t, 5, a)
	assert.Equal(t, 3, len(p.tickets))

}

func TestRunWithKey(t *testing.T) {
	key1 := TaskKey("Task-1")
	key2 := TaskKey("Task-1")
	key3 := TaskKey("Task-1")
	key4 := TaskKey("Task-1")
	key5 := TaskKey("Task-1")
	assert.True(t, key1 == key2 && key2 == key3 && key3 == key4 && key4 == key5)

	p := NewWorkerPool(10)
	err1 := p.RunWithKey(nil, key1, func() { time.Sleep(time.Second) })
	err2 := p.RunWithKey(nil, key2, func() { time.Sleep(time.Second) })
	err3 := p.RunWithKey(nil, key3, func() { time.Sleep(time.Second) })
	err4 := p.RunWithKey(nil, key4, func() { time.Sleep(time.Second) })
	err5 := p.RunWithKey(nil, key5, func() { time.Sleep(time.Second) })

	assert.Nil(t, err1)
	assert.Equal(t, ErrAlreadyRunning, err2)
	assert.Equal(t, ErrAlreadyRunning, err3)
	assert.Equal(t, ErrAlreadyRunning, err4)
	assert.Equal(t, ErrAlreadyRunning, err5)

	key6 := TaskKey("Task-6")
	key7 := TaskKey("Task-7")
	err6 := p.RunWithKey(nil, key6, func() { time.Sleep(time.Second) })
	err7 := p.RunWithKey(nil, key7, func() { time.Sleep(time.Second) })
	assert.Nil(t, err6)
	assert.Nil(t, err7)
}

func TestRunWithKeyConcurrent(t *testing.T) {
	key := TaskKey("Task-1")
	p := NewWorkerPool(10)
	var errors = make([]error, 5)

	go func() { errors[0] = p.RunWithKey(nil, key, func() { time.Sleep(time.Second) }) }()
	go func() { errors[1] = p.RunWithKey(nil, key, func() { time.Sleep(time.Second) }) }()
	go func() { errors[2] = p.RunWithKey(nil, key, func() { time.Sleep(time.Second) }) }()
	go func() { errors[3] = p.RunWithKey(nil, key, func() { time.Sleep(time.Second) }) }()
	go func() { errors[4] = p.RunWithKey(nil, key, func() { time.Sleep(time.Second) }) }()

	time.Sleep(time.Millisecond * 10)
	var successNum, alreadyRunningErrNum int

	for _, err := range errors {
		if err == nil {
			successNum += 1
		} else if err == ErrAlreadyRunning {
			alreadyRunningErrNum += 1
		}
	}
	assert.Equal(t, 1, successNum)
	assert.Equal(t, 4, alreadyRunningErrNum)
	assert.True(t, p.isRunning(key))
	assert.Equal(t, 9, len(p.tickets))
}
