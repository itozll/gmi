package glog

import "testing"

func TestGlog(t *testing.T) {
	l := New(&Options{
		TimeFormat: "2006-01-02 15:04:05",
	})
	l.Info("ABC")
	dfl.Info("ABC")
}
