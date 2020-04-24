package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestPidFileWriteOk(t *testing.T) {
	cases := []struct {
		dirName  string
		fileName string
		pid      int
		pidMod   string
		want1    tPidMark
		want2    error
	}{
		{os.TempDir(), "testpidfilewrite", 12345, "", tPidMark{os.TempDir(), "testpidfilewrite", 12345}, nil},
		{os.TempDir(), ".testpidfilewrite", 12345, "", tPidMark{os.TempDir(), ".testpidfilewrite", 12345}, nil},
		{os.TempDir(), "testpidfilewrite", 12345, "mod", tPidMark{os.TempDir(), "testpidfilewrite", -1}, errors.New("The PID file already exist, but is contents are incorrect")},
	}
	for _, c := range cases {
		defer os.Remove(c.dirName + "/" + c.fileName)
		pidFile, _ := os.Create(c.dirName + "/" + c.fileName)
		defer pidFile.Close()
		_, _ = pidFile.WriteString(strings.Join(strings.Split(strconv.Itoa(c.pid), ""), c.pidMod))
		got1, got2 := initPidMark(c.dirName, c.fileName, c.pid)
		if got1 != c.want1 { //|| (got2 != c.want2) {
			t.Errorf("TestPidFileWriteOk(%q,%q,%v) == %v,%v, want %v,%v", c.dirName, c.fileName, c.pid, got1, got2, c.want1, c.want2)
		}

		/*got := pidFileWrite(c.dirName, c.fileName, c.pid)
		if got != c.want {
			t.Errorf("TestPidFileWriteOk(%q,%q,%v) == %v, want %v", c.dirName, c.fileName, c.pid, got, c.want)
		}
		st, err := ioutil.ReadFile(c.dirName + "/" + c.fileName)
		if (err == nil) && (string(st) != "12345") {
			t.Errorf("The recorded and read values do not match: read == %v, want %q", c.pid, st)
		}*/

	}
}
