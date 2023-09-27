package utils

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/27 10:07
  @describe :
*/

type Signal struct {
	subscribers map[chan int]interface{}
	closed      bool
}

func (s *Signal) Publish(i int) {
	if s.closed {
		panic("signal is closed")
	}
	for ch := range s.subscribers {
		go func(in chan int) {
			in <- i
		}(ch)
	}
}

func (s *Signal) Subscribe() <-chan int {
	if s.closed {
		panic("signal is closed")
	}
	output := make(chan int)
	s.subscribers[output] = nil
	return output
}

func (s *Signal) Close() {
	if s.closed {
		return
	}
	for ch := range s.subscribers {
		close(ch)
		delete(s.subscribers, ch)
	}
}

func NewSignal() *Signal {
	return &Signal{subscribers: make(map[chan int]interface{})}
}
