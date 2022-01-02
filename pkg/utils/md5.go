package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func Md5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func BuildToken() string {
	worker, _ := NewWorker(2)
	return strconv.FormatInt(worker.GetId(), 10)
}
