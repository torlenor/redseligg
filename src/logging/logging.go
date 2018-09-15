package logging

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Init the logging framework
// has to be called only once
func Init() {
	logrus.SetFormatter(new(MyFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

// Get a logger with prefix name
func Get(name string) *logrus.Entry {
	return logrus.WithField("name", name)
}

type MyFormatter struct{}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	name, ok := entry.Data["name"]
	if !ok {
		name = "default"
	}
	fmt.Fprintf(b, "%s [%-5.5s] (%s): %s\n", entry.Time.Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String()), name, entry.Message)
	return b.Bytes(), nil
}
