package logr_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vrischmann/logr"
)

func makeBuf(b byte) []byte {
	buf := make([]byte, 1024)
	for i := 0; i < len(buf); i++ {
		buf[i] = b
	}

	return buf
}

func readFile(t testing.TB, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	require.Nil(t, err)

	return data
}

func checkEqual(t testing.TB, buf []byte, b byte) error {
	for _, v := range buf {
		if v != b {
			return fmt.Errorf("%v != %v", v, b)
		}
	}

	return nil
}

func TestRotateMaxSize(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "logr")
	require.Nil(t, err)

	rw, err := logr.NewWriterFromFile(f)
	require.Nil(t, err)

	now := time.Now()
	{
		n, err := rw.Write(makeBuf(0xFF))
		require.Nil(t, err)
		require.Equal(t, 1024, n)

		rw.MaxSize(512)

		n, err = rw.Write(makeBuf(0xFE))
		require.Nil(t, err)
		require.Equal(t, 1024, n)
	}

	newData := readFile(t, f.Name())
	require.Nil(t, checkEqual(t, newData, 0xFE))

	rotatedData := readFile(t, f.Name()+"."+now.Format(logr.SuffixTimeFormat))
	require.Nil(t, checkEqual(t, rotatedData, 0xFF))
}
