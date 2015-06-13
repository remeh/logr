package logr

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	// SuffixTimeFormat is the default format used for the suffix date and time on each rotated log.
	SuffixTimeFormat = "2006-01-02_1504"
)

// RotatingWriter is a io.Writer which wraps a *os.File, suitable for log rotation.
type RotatingWriter struct {
	lock     sync.Mutex
	filename string
	file     *os.File

	schedule struct {
		hour   int
		minute int
	}

	currentSize int64
	maxSize     int64
}

// NewWriter creates a new file and returns a rotating writer.
func NewWriter(filename string) (*RotatingWriter, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return NewWriterFromFile(file)
}

// NewWriterFromFile creates a rotating writer using the provided file as base.
//
// The caller must take care to not close the file it provides here, as the RotatingWriter
// will do it automatically when rotating.
func NewWriterFromFile(file *os.File) (*RotatingWriter, error) {
	w := &RotatingWriter{
		filename: file.Name(),
		file:     file,
		maxSize:  -1,
	}

	w.schedule.hour = -1
	w.schedule.minute = -1

	if err := w.readCurrentSize(); err != nil {
		return nil, err
	}

	return w, nil
}

// readCurrentSize reads the current size from the file
func (w *RotatingWriter) readCurrentSize() error {
	fi, err := w.file.Stat()
	if err != nil {
		return err
	}

	w.currentSize = fi.Size()

	return nil
}

// Schedule set the time at which to rotate, each day
func (w *RotatingWriter) Schedule(hour, minute int) *RotatingWriter {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.schedule.hour = hour
	w.schedule.minute = minute

	return w
}

// MaxSize set the size at which to rotate the file
func (w *RotatingWriter) MaxSize(s int64) *RotatingWriter {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.maxSize = s

	return w
}

func (w *RotatingWriter) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.schedule.hour > -1 {
		now := time.Now()
		if now.Hour() == w.schedule.hour {
			if err := w.rotate(); err != nil {
				return -1, err
			}
		}
	}

	if w.maxSize > -1 {
		if w.currentSize >= w.maxSize {
			if err := w.rotate(); err != nil {
				return -1, err
			}
		}
	}

	n, err := w.file.Write(b)
	w.currentSize += int64(n)

	return n, err
}

// rotate rotates the file. must be called while having the file lock
func (w *RotatingWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	{
		destName := w.filename + "." + time.Now().Format(SuffixTimeFormat)
		_, err := os.Stat(destName)
		if err != nil && !os.IsNotExist(err) {
			fmt.Println(err.Error())
			return err
		}

		if err := os.Rename(w.filename, destName); err != nil {
			return err
		}
	}

	{
		file, err := os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		w.file = file
		w.currentSize = 0
	}

	return nil
}
