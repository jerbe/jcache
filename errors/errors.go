package errors

import (
	"github.com/jerbe/go-errors"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/1 11:14
  @describe :
*/

var (
	Nil              = errors.New("jcache: nil")
	ErrNoCacheClient = errors.New("no cache client init")
)

// ErrorValuer 带返回错误的值对象
type ErrorValuer interface {
	Err() error
}
