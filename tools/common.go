package tools

import (
	"gii/glog"
	"io"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		glog.Error(err)
	}
}
