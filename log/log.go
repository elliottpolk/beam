//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/log/log.go
//
package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	logger.Println(args...)
}

func NewError(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Error(err error, message string) {
	logger.Error(errors.Wrap(err, message))
}

func Errorf(err error, format string, args ...interface{}) {
	logger.Error(errors.Wrapf(err, format, args...))
}

func Errorln(err error, message string) {
	logger.Errorln(errors.Wrap(err, message))
}

func Fatal(args ...interface{}) {
	logger.Panic(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalln(args ...interface{}) {
	logger.Panicln(args...)
}
