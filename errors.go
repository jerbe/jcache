package jcache

import (
	"context"
	"errors"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/1 11:14
  @describe :
*/

var (
	ErrEmpty         = errors.New("empty")
	ErrNoCacheClient = errors.New("no cache client init")
	ErrNoRecord      = errors.New("no record")
	ErrCanceled      = context.Canceled
)
