package bytime_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/rotation/bytime"
	"github.com/modern-go/test/must"
	"io/ioutil"
	"os"
	"time"
)

func TestRotation(t *testing.T) {
	t.Run("simple write", test.Case(func(ctx context.Context) {
		os.Remove("/tmp/test.log")
		writer, err := bytime.NewRotationWriter(func(writer *bytime.Writer) {
			writer.WritePath = "/tmp/test.log"
		})
		must.Nil(err)
		_, err = writer.Write([]byte("hello"))
		must.Nil(err)
		content, err := ioutil.ReadFile("/tmp/test.log")
		must.Nil(err)
		must.Equal("hello", string(content))
	}))
	t.Run("rotate every second", test.Case(func(ctx context.Context) {
		os.RemoveAll("/tmp/testlog/")
		os.Mkdir("/tmp/testlog/", 0755)
		writer, err := bytime.NewRotationWriter(func(writer *bytime.Writer) {
			writer.WritePath = "/tmp/testlog/test.log"
			writer.ArchiveFilePattern = "test-{time,goTime,2006-01-02T15:04:05.999999999}.log"
			writer.ArchiveKeepDuration = time.Millisecond * 30
			writer.Interval = time.Millisecond
		})
		must.Nil(err)
		for i := 0; i < 100; i++ {
			_, err = writer.Write([]byte("hello"))
			must.Nil(err)
			time.Sleep(time.Millisecond)
		}
		files, _ := ioutil.ReadDir("/tmp/testlog/")
		must.Pass(len(files) > 0)
		must.Pass(len(files) < 30)
	}))
}
