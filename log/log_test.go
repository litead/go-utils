package log

import "testing"

func Test_Write(t *testing.T) {
	logger := Logger{ToStdErr: true, WithFile: true}

	logger.Start()
	logger.Infoln("this", "is a test log")
	logger.Close()
}
