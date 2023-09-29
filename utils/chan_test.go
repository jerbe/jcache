package utils

import (
	"fmt"
	"testing"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/27 23:31
  @describe :
*/

func TestSignal(t *testing.T) {
	sig := NewSignal()
	ch := sig.Subscribe()

	go func() {
		//ticker := time.Tick(time.Second)
		var i = 0
		for {
			i++
			sig.Publish(i)
			//select {
			//case <-ticker:
			//	i++
			//	sig.Publish(i)
			//}
		}
	}()

	for {
		select {
		case s := <-ch:
			fmt.Println("sig:", s)
			if s == 0 {
				return
			}
		}
	}
}
