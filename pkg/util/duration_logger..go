package util

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type durationObj struct {
	sync.RWMutex
	objs map[string]time.Time
}

func (d *durationObj) Start(name string) {
	d.Lock()
	defer d.Unlock()

	_, ok := d.objs[name]
	if ok {
		panic(fmt.Sprintf("duration logger with same name = %s already exist!", name))
	}
	d.objs[name] = time.Now()
}

func (d *durationObj) End(name string) {
	d.Lock()
	defer d.Unlock()

	val, ok := d.objs[name]
	if !ok {
		panic(fmt.Sprintf("duration logger with name = %s does not exist!", name))
	}
	slog.Info("duration logger", "name", name, "duration", time.Since(val))
	delete(d.objs, name)
}

var singletonDurationObj DurationLoggerInterface = &durationObj{objs: make(map[string]time.Time)}

type DurationLoggerInterface interface {
	Start(name string)
	End(name string)
}

func DurationLogger() DurationLoggerInterface {
	return singletonDurationObj
}
