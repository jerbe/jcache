package utils

import "testing"

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/27 23:31
  @describe :
*/

func TestSignal(t *testing.T) {
	sig := NewSignal()
	ch := sig.Subscribe()
	sig.Unsubscribe(ch)
}
