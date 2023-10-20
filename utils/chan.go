package utils

import "sync"

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/27 10:07
  @describe :
*/

type Signal struct {
	subscribers map[*SignalSubscriber]struct{}
	rwMute      sync.RWMutex
	closed      bool
}

type SignalSubscriber struct {
	ch     chan interface{}
	C      <-chan interface{}
	close  chan struct{}
	signal *Signal
}

func (sb *SignalSubscriber) Close() error {
	sb.signal.Unsubscribe(sb)
	return nil
}

func (s *Signal) Publish(i interface{}) {
	if s.closed {
		panic("signal is closed")
	}
	s.rwMute.RLock()
	defer s.rwMute.RUnlock()
	for sb := range s.subscribers {
		go func(in *SignalSubscriber) {
			select {
			case <-in.close:
				return
			case in.ch <- i:
			}
		}(sb)
	}
}

func (s *Signal) Subscribe() *SignalSubscriber {
	if s.closed {
		panic("signal is closed")
	}
	ch := make(chan interface{})
	sb := &SignalSubscriber{
		ch:     ch,
		C:      ch,
		close:  make(chan struct{}),
		signal: s,
	}
	s.rwMute.Lock()
	defer s.rwMute.Unlock()
	s.subscribers[sb] = struct{}{}
	return sb
}

func (s *Signal) Unsubscribe(sb *SignalSubscriber) {
	s.rwMute.Lock()
	defer s.rwMute.Unlock()
	if _, ok := s.subscribers[sb]; ok {
		delete(s.subscribers, sb)
		close(sb.close)
	}
}

func (s *Signal) Close() {
	if s.closed {
		return
	}
	s.rwMute.Lock()
	defer s.rwMute.Unlock()

	for sb := range s.subscribers {
		delete(s.subscribers, sb)
	}
	s.closed = true
}

func NewSignal() *Signal {
	return &Signal{subscribers: make(map[*SignalSubscriber]struct{})}
}
