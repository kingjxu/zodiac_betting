package rolllog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Level uint32

// Log levels.
const (
	LEVEL_TRACE Level = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_PANIC
)

var levelNames = []string{
	"TRACE",
	"DEBUG",
	"INFO ",
	"WARN ",
	"ERROR",
	"FATAL",
	"PANIC",
}

type Logger struct {
	sync.RWMutex
	w     *RollWriter
	l     *baseLogger
	level Level
}

var defaultLogger *Logger

func New(path string, baseName string) (*Logger, error) {
	l := &Logger{}
	err := l.Init(path, baseName)

	return l, err
}

func (l *Logger) Init(path string, baseName string) error {
	l.level = LEVEL_INFO

	if path == "" || baseName == "" {
		l.l = NewbaseLogger(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		return nil
	}

	var err error
	l.w, err = NewRollWriter(path, baseName)
	if err != nil {
		fmt.Printf("NewRollWriter()!err:%+v", err)
		return err
	}

	l.l = NewbaseLogger(l.w, "", log.LstdFlags|log.Lshortfile)

	return err
}

func (l *Logger) SetFlags(flag int) {
	l.l.SetFlags(flag)
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetMaxFileSize(size uint32) {
	l.w.SetMaxFileSize(size)
}

func (l *Logger) SetMaxFileCnt(count uint32) {
	l.w.SetMaxFileCnt(count)
}

func (l *Logger) output(level Level, format string, v ...interface{}) {
	if level >= l.level {
		l.l.Output(3, levelNames[level], fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.output(LEVEL_TRACE, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.output(LEVEL_DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.output(LEVEL_INFO, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(LEVEL_WARN, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.output(LEVEL_ERROR, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.output(LEVEL_FATAL, format, v...)
}

func (l *Logger) Panic(format string, v ...interface{}) {
	l.output(LEVEL_PANIC, format, v...)
}

func Trace(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_TRACE, format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_FATAL, format, v...)
	os.Exit(1)
}

func Panic(format string, v ...interface{}) {
	defaultLogger.output(LEVEL_PANIC, format, v...)
	panic(fmt.Sprintf("["+levelNames[LEVEL_PANIC]+"]"+format, v...))
}

func Init(path string, baseName string) error {
	var err error
	defaultLogger, err = New(path, baseName)
	if err != nil {
		fmt.Printf("failed to new defaultLogger!err:%+v", err)
		return err
	}

	return nil
}

func SetFlags(flag int) {
	defaultLogger.l.SetFlags(flag)
}

func SetLevel(level Level) {
	defaultLogger.level = level
}

func SetMaxFileSize(size uint32) {
	if defaultLogger.w != nil {
		defaultLogger.w.SetMaxFileSize(size)
	}
}

func SetMaxFileCnt(count uint32) {
	if defaultLogger.w != nil {
		defaultLogger.w.SetMaxFileCnt(count)
	}
}

func getExeBaseName() (string, error) {
	exeAbsPathName, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}

	exeFullName := filepath.Base(exeAbsPathName)
	exeBaseName := strings.TrimSuffix(exeFullName, filepath.Ext(exeFullName))

	return exeBaseName, nil
}

func init() {
	exeAbsPathName, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}

	exePath := filepath.Dir(exeAbsPathName)
	basePath := exePath
	if filepath.Base(exePath) == "bin" {
		basePath = filepath.Dir(exePath)
	}

	baseName, err := getExeBaseName()
	if err != nil {
		log.Fatal("Failed to get baseName,err:%+v", err)
	}

	Init(basePath+"/log/", baseName)
}
