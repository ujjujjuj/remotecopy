package client

import (
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
)

func Copy() (string, error) {
	return clipboard.ReadAll()
}

func Paste(str string) {
	robotgo.TypeStr(str)
}

func Test() {
	robotgo.KeyPress("Enter")
}
