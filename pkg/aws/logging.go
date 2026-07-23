/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"github.com/sirupsen/logrus"
)

type LoggerWrapper struct {
	logrusLogger *logrus.Logger
}

func NewLoggerWrapper(logrusLog *logrus.Logger) *LoggerWrapper {
	if logrusLog == nil {
		return nil
	}
	return &LoggerWrapper{
		logrusLogger: logrusLog,
	}
}

func (lw *LoggerWrapper) GetLevel() int {
	return int(lw.logrusLogger.GetLevel())
}

func (lw *LoggerWrapper) Debug(args ...interface{}) {
	lw.logrusLogger.Debug(args...)
}

func (lw *LoggerWrapper) Info(args ...interface{}) {
	lw.logrusLogger.Info(args...)
}

func (lw *LoggerWrapper) Warn(args ...interface{}) {
	lw.logrusLogger.Warn(args...)
}

func (lw *LoggerWrapper) Error(args ...interface{}) {
	lw.logrusLogger.Error(args...)
}

func (lw *LoggerWrapper) Fatal(args ...interface{}) {
	lw.logrusLogger.Fatal(args...)
}
