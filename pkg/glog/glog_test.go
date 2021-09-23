package glog

import "testing"

func TestGlog(t *testing.T) {
	l := New(&Options{
		TimeFormat: "2006-01-02 15:04:05",
	})
	l.Info("ABC1")
	dfl.Info("ABC")

	SetDefaultLog(l)

	var i int
	dfl.WithFields("a", "v", i, "s", "s", "dd").Info("V")
}
