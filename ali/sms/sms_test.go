package sms

import (
	"testing"
)

func Test_Sms(t *testing.T) {
	sender := NewSender("LTAIFINnKi6pYhAL", "WO1K8DSUgewaJyfkpemtN2VdiYh5f8")
	sender.Send("18701260136", "百米客", "SMS_53610156", `{"code":"123","product":"456"}`)
}
