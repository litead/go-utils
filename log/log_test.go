package log

import "testing"

func Test_Write(t *testing.T) {
	logger := Logger{ToStdErr: true, WithFile: true}

	logger.Start()
	logger.Write(Info, "this", "is a test log")
	logger.Close()
}
