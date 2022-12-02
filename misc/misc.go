package misc

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// timestamp for filenames unique to the second
const DATE_POG = "20060102150405"

// DeferError
// accounts for an at-close function that
// returns an error at its close
func DeferError(f func() error) {
	err := f()

	if nil != err {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		_, _ = fmt.Fprintf(os.Stderr,
			"[%s] error in DeferError from file: %s line %d\n"+
				" error: %s\n\t(may be harmless!)",
			time.Now().UTC().Format(time.RFC822),
			file, line, err.Error())
	}
}

// WriteSB Add a series of strings to a strings.Builder
func WriteSB(sb *strings.Builder, inputStrings ...string) {
	if nil == sb || nil == inputStrings {
		panic("null pointer instead of *strings.Builder or inputStrings in misc.WriteSB()")
	}
	if len(inputStrings) <= 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Got 0-length array of strings in misc.WriteSB()")
		return
	}
	for _, val := range inputStrings {
		_, err := sb.WriteString(val)
		if nil != err {
			_, _ = fmt.Fprintf(os.Stderr,
				"Got error in misc.WriteSB() while writing strings.\nError: %s",
				err.Error())
			for ix, str := range inputStrings {
				_, _ = fmt.Fprintf(os.Stderr, "%05d: [ %s ]", ix, str)
			}
			panic("panic in misc.WriteSB")
		}
	}
}

func IsStringSet(s *string) (isSet bool) {
	if nil != s && "" != *s {
		return true
	}
	return false
}
