package perror

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Error interface {
	Code() int
	Err() error
	Error() string
	Msg() string
	GetStack() string
	Unwrap() error
}

type fullerror struct {
	code int
	err  error
	msg  string
	*stack
}

// golang 官方error转换成内部error返回
func GoErr(cause error) fullerror {
	return fullerror{
		code:  -1,
		err:   cause,
		msg:   "内部错误",
		stack: callers(),
	}
}

// 将Go error进行包装后转换成我们自己的Error
func New(code int, msg string, cause error) fullerror {
	return fullerror{
		code:  code,
		err:   fmt.Errorf(msg+":%w", cause),
		msg:   msg,
		stack: callers(),
	}
}

// 返回一个我们自己的error
func NewNil(code int, msg string) fullerror {
	return fullerror{
		code:  code,
		err:   errors.New(msg),
		msg:   msg,
		stack: callers(),
	}
}

func (f fullerror) Err() error {
	return f.err
}

func (f fullerror) Msg() string {
	return f.msg
}

func (f fullerror) Code() int {
	return f.code
}

func (f fullerror) Error() string {
	s := f.TopStackTrace()
	return fmt.Sprintf("%s-%s:%d", f.err.Error(), s, s)
}

func (f fullerror) GetStack() string {
	if f.stack == nil {
		return ""
	}
	var err []string
	s := *f.stack
	for i := 0; i < len(s); i++ {
		f := Frame(s[i])
		err = append(err, fmt.Sprintf("\t\tfun:%s\tline:%d \n", f.name(), f.line()))
	}
	return strings.Join(err, "")
}

func (f fullerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, fmt.Sprintf("%v", f.err))
	case 'q':
		fmt.Fprintf(s, "%q", fmt.Sprintf("%v", f.err))
	}
}

func (w fullerror) Cause() error {
	return w.err
}
func (w fullerror) Unwrap() error {
	return w.err
}
