package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	stopReq bool
	stopped chan struct{}

	MsgQueue messageQueue
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.MsgQueue = messageQueue{}
	l.stopped = make(chan struct{})

	go l.eventProcess()
}

func (l *Loop) eventProcess() {
	for {
		op := l.MsgQueue.Pull()
		if op == nil {
			continue
		}

		if update := op.Do(l.next); update {
			l.Receiver.Update(l.next)
			l.next, l.prev = l.prev, l.next
		}

		if l.stopReq {
			close(l.stopped)
			return
		}
	}
}

func (l *Loop) Post(op Operation) {
	if op != nil {
		l.MsgQueue.Push(op)
	}
}

func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(screen.Texture) {
		l.stopReq = true
	}))
	<-l.stopped
}

type messageQueue struct {
	Queue   []Operation
	mu      sync.Mutex
	blocked chan struct{}
}

func (mq *messageQueue) Push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.Queue = append(mq.Queue, op)

	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *messageQueue) Pull() Operation {
	mq.mu.Lock()
	for len(mq.Queue) == 0 {
		mq.blocked = make(chan struct{})
		blocked := mq.blocked
		mq.mu.Unlock()
		<-blocked
		mq.mu.Lock()
	}
	op := mq.Queue[0]
	mq.Queue[0] = nil
	mq.Queue = mq.Queue[1:]
	mq.mu.Unlock()
	return op
}
