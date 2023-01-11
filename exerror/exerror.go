package exerror

import "fmt"

type ErgoError struct {
	Message string
}

func (ee *ErgoError) Error() string {
	return ee.Message
}

func (ee *ErgoError) String() string {
	return ee.Message
}

func Bomb(format string, args ...interface{}) {
	panic(ErgoError{Message: fmt.Sprintf(format, args...)})
}

func Dangerous(v interface{}) {
	if v == nil {
		return
	}

	switch t := v.(type) {
	case string:
		if t != "" {
			panic(ErgoError{Message: t})
		}
	case error:
		panic(ErgoError{Message: t.Error()})
	}
}

func Boka(value string, v interface{}) {
	if v == nil {
		return
	}
	Bomb(value)
}

// CheckAndExit check & exit
func CheckAndExit(err error) {
	if err != nil {
		panic(err)
	}
}
