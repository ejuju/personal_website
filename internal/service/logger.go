package service

import (
	"errors"
	"io"
	"os"
)

type Logger struct {
	Writer io.Writer
}

func NewLogger(fpath string) (l Logger, err error) {
	if fpath == "" {
		return l, errors.New("empty config file path")
	}

	switch fpath {
	default:
		l.Writer, err = openLogfile(fpath)
		if err != nil {
			return l, err
		}
	case "stderr":
		l.Writer = os.Stderr
	case "stdout":
		l.Writer = os.Stdout
	}

	return l, nil
}

func (l Logger) Log(s string) {
	_, err := io.WriteString(l.Writer, s+"\n")
	if err != nil {
		panic(err)
	}
}

func openLogfile(fpath string) (f *os.File, err error) {
	f, err = os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fstat.Mode().IsRegular() {
		return nil, errors.New("logfile must be regular file")
	}
	return f, nil
}
