package pkg

import (
	"fmt"
	fs "github.com/facebookgo/stack"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var defaultFormatter = &logrus.TextFormatter{DisableColors: true}

type PathMap map[logrus.Level]string

type WriterMap map[logrus.Level]io.Writer

type LoggerHook struct {
	paths     PathMap
	writers   WriterMap
	levels    []logrus.Level
	lock      *sync.Mutex
	formatter logrus.Formatter

	defaultPath      string
	defaultWriter    io.Writer
	hasDefaultPath   bool
	hasDefaultWriter bool
}

func (hook *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func addStack(entry *logrus.Entry) {
	var skipFrames int
	if len(entry.Data) == 0 {
		// When WithField(s) is not used, we have 8 logrus frames to skip.
		skipFrames = 8
	} else {
		// When WithField(s) is used, we have 6 logrus frames to skip.
		skipFrames = 6
	}

	var frames fs.Stack

	// Get the complete stack track past skipFrames count.
	_frames := fs.Callers(skipFrames)

	// Remove logrus's own frames that seem to appear after the code is through
	// certain hoops. e.g. http handler in a separate package.
	// This is a workaround.
	for _, frame := range _frames {
		if !strings.Contains(frame.File, "github.com/sirupsen/logrus") {
			frames = append(frames, frame)
		}
	}

	if len(frames) > 0 {
		// If we have a frame, we set it to "caller" field for assigned levels.
		entry.Data["caller"] = frames[0]

		// Set the available frames to "stack" field.
		if entry.Level < logrus.WarnLevel {
			entry.Data["stack"] = frames
		}
	}
}

func (hook *LoggerHook) fileWrite(entry *logrus.Entry) error {
	var (
		fd   *os.File
		path string
		msg  []byte
		err  error
		ok   bool
	)

	hook.lock.Lock()
	defer hook.lock.Unlock()

	if path, ok = hook.paths[entry.Level]; !ok {
		if hook.hasDefaultPath {
			path = hook.defaultPath
		} else {
			return nil
		}
	}
	addStack(entry)

	dir := filepath.Dir(path)
	os.MkdirAll(dir, os.ModePerm)

	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open logfile:", path, err)
		return err
	}
	defer fd.Close()

	// use our formatter instead of entry.String()
	msg, err = hook.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	fd.Write(msg)
	return nil
}

func (hook *LoggerHook) Fire(entry *logrus.Entry) error {
	if hook.writers != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	} else if hook.paths != nil || hook.hasDefaultPath {
		return hook.fileWrite(entry)
	}

	return nil
}

func (hook *LoggerHook) ioWrite(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    []byte
		err    error
		ok     bool
	)

	hook.lock.Lock()
	defer hook.lock.Unlock()

	if writer, ok = hook.writers[entry.Level]; !ok {
		if hook.hasDefaultWriter {
			writer = hook.defaultWriter
		} else {
			return nil
		}
	}
	addStack(entry)

	// use our formatter instead of entry.String()
	msg, err = hook.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	_, err = writer.Write(msg)
	return err
}

func (hook *LoggerHook) SetFormatter(formatter logrus.Formatter) {
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}

	hook.formatter = formatter
}

// SetDefaultPath sets default path for levels that don't have any defined output path.
func (hook *LoggerHook) SetDefaultPath(defaultPath string) {
	hook.defaultPath = defaultPath
	hook.hasDefaultPath = true
}

// SetDefaultWriter sets default writer for levels that don't have any defined writer.
func (hook *LoggerHook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
}

func NewHook(output interface{}, formatter logrus.Formatter) *LoggerHook {
	hook := &LoggerHook{
		lock: new(sync.Mutex),
	}

	hook.SetFormatter(formatter)

	switch output.(type) {
	case string:
		hook.SetDefaultPath(output.(string))
		break
	case io.Writer:
		hook.SetDefaultWriter(output.(io.Writer))
		break
	case PathMap:
		hook.paths = output.(PathMap)
		for level := range output.(PathMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	case WriterMap:
		hook.writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	default:
		panic(fmt.Sprintf("unsupported level map type: %v", reflect.TypeOf(output)))
	}

	return hook
}
