package logr

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMakeDestName(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Hour * 24)

	rw := RotatingWriter{
		filename:  "/var/log/logr.log",
		prefix:    false,
		startDate: now,
	}
	n := rw.makeDestName()

	expected := fmt.Sprintf("/var/log/logr.log.%s", now.Format(TimeFormat))
	require.Equal(t, expected, n)

	rw.prefix = true
	n = rw.makeDestName()

	expected = fmt.Sprintf("/var/log/logr.%s.log", now.Format(TimeFormat))
	require.Equal(t, expected, n)
}
