// Copyright (c) 2017-2018 Alexander Eichhorn
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package xmp

import (
	"log"
	"strconv"
	"strings"
)

type LogLevelType int
type LogType struct {
	bool
}

type Logger interface {
	Error(v ...interface{})
	Errorf(s string, v ...interface{})
	Warn(v ...interface{})
	Warnf(s string, v ...interface{})
	Info(v ...interface{})
	Infof(s string, v ...interface{})
	Debug(v ...interface{})
	Debugf(s string, v ...interface{})
}

const (
	LogLevelInvalid LogLevelType = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

func (l LogLevelType) Prefix() string {
	switch l {
	case LogLevelDebug:
		return "Debug:"
	case LogLevelInfo:
		return "Info:"
	case LogLevelWarning:
		return "Warn:"
	case LogLevelError:
		return "Error:"
	default:
		return strconv.Itoa(int(l))
	}
}

var (
	logLevel LogLevelType = LogLevelWarning
	Log      Logger       = &LogType{}
)

// package level logging control
func LogLevel() LogLevelType {
	return logLevel
}

func SetLogLevel(l LogLevelType) {
	logLevel = l
}

func SetLogger(v Logger) {
	Log = v
}

// internal
func defaultLogger(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	log.Printf(format, args...)
}

func output(lvl LogLevelType, v ...interface{}) {
	m := append(make([]interface{}, 0, len(v)+1), lvl.Prefix())
	m = append(m, v...)
	defaultLogger("%s", m...)
}

func outputf(lvl LogLevelType, s string, v ...interface{}) {
	s = strings.Join([]string{"%s ", s}, "")
	m := append(make([]interface{}, 0, len(v)+1), lvl.Prefix())
	m = append(m, v...)
	defaultLogger(s, m...)
}

func (x LogType) Error(v ...interface{}) {
	if logLevel > LogLevelError {
		return
	}
	output(LogLevelError, v...)
}

func (x LogType) Errorf(s string, v ...interface{}) {
	if logLevel > LogLevelError {
		return
	}
	outputf(LogLevelError, s, v...)
}

func (x LogType) Warn(v ...interface{}) {
	if logLevel > LogLevelWarning {
		return
	}
	output(LogLevelWarning, v...)
}

func (x LogType) Warnf(s string, v ...interface{}) {
	if logLevel > LogLevelWarning {
		return
	}
	outputf(LogLevelWarning, s, v...)
}

func (x LogType) Info(v ...interface{}) {
	if logLevel > LogLevelInfo {
		return
	}
	output(LogLevelInfo, v...)
}

func (x LogType) Infof(s string, v ...interface{}) {
	if logLevel > LogLevelInfo {
		return
	}
	outputf(LogLevelInfo, s, v...)
}

func (x LogType) Debug(v ...interface{}) {
	if logLevel > LogLevelDebug {
		return
	}
	output(LogLevelDebug, v...)
}

func (x LogType) Debugf(s string, v ...interface{}) {
	if logLevel > LogLevelDebug {
		return
	}
	outputf(LogLevelDebug, s, v...)
}

// package level forwarders to the real logger implementation
// func Error(v ...interface{})            { Log.Error(v...) }
// func Errorf(s string, v ...interface{}) { Log.Errorf(s, v...) }
// func Warn(v ...interface{})             { Log.Warn(v...) }
// func Warnf(s string, v ...interface{})  { Log.Warnf(s, v...) }
// func Info(v ...interface{})             { Log.Info(v...) }
// func Infof(s string, v ...interface{})  { Log.Infof(s, v...) }
// func Debug(v ...interface{})            { Log.Debug(v...) }
// func Debugf(s string, v ...interface{}) { Log.Debugf(s, v...) }
