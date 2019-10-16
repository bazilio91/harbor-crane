package main

import (
	"github.com/bazilio91/harbor-crane/crane"
	"github.com/sirupsen/logrus"
)

func main() {

	c := crane.NewCrane()

	err := c.Sync()
	if err != nil {
		logrus.WithError(err).Error("fail")
		panic(err)
	}
}
