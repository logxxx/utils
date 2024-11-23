package notify

import (
	"context"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/logxxx/utils/runutil"
)

type INotify interface {
	Pub(event *fsnotify.Event) error
	PubGlobal(event *fsnotify.Event) error
	Sub(fn func(e *fsnotify.Event))
	Close() error
}

type NotifyEngine struct {
	ctx      context.Context
	notifies []*Notify
	closeFn  func()
}

type Notify struct {
	engine     *NotifyEngine
	dirs       []string
	notifyChan chan *fsnotify.Event
	ctx        context.Context
	closeFn    func()
	name       string
	isClosed   bool
}

func NewNotifyEngine() *NotifyEngine {
	ctx, cancel := context.WithCancel(context.Background())
	e := &NotifyEngine{ctx: ctx, closeFn: cancel}
	return e
}

func newNotify(ctx context.Context, engine *NotifyEngine, name string) *Notify {
	newCtx, cancel := context.WithCancel(ctx)
	return &Notify{engine: engine, notifyChan: make(chan *fsnotify.Event, 1024), ctx: newCtx, closeFn: cancel, name: name}
}

// name NOT need unique, just a mark
func (e *NotifyEngine) NewNotify(name string) INotify {
	sn := newNotify(e.ctx, e, name)
	e.notifies = append(e.notifies, sn)
	return sn
}

func (e *NotifyEngine) Close() error {
	if e.closeFn != nil {
		e.closeFn()
	}
	for _, n := range e.notifies {
		n.Close()
	}
	return nil
}

func (e *NotifyEngine) pub(event *fsnotify.Event) {
	for _, elemSn := range e.notifies {
		elemSn.Pub(event)
	}
}

func (sn *Notify) Sub(fn func(e *fsnotify.Event)) {
	runutil.GoRunSafe(func() {
		for {
			if sn.isClosed {
				return
			}
			select {
			case <-sn.ctx.Done():
				return
			case e := <-sn.notifyChan:
				fn(e)
			}
		}
	})
}

func (sn *Notify) PubGlobal(event *fsnotify.Event) error {
	sn.engine.pub(event)
	return nil
}

func (sn *Notify) Pub(event *fsnotify.Event) error {
	select {
	case sn.notifyChan <- event:
	default:
		return errors.New("pub queue is full")
	}
	return nil
}

func (sn *Notify) Close() error {
	if sn.isClosed {
		return nil
	}
	sn.isClosed = true
	if sn.closeFn != nil {
		sn.closeFn()
	}
	close(sn.notifyChan)
	return nil
}
