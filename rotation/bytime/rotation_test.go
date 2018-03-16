package bytime_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/rotation/bytime"
	"github.com/modern-go/test/must"
	"io/ioutil"
	"os"
)

func TestRotation(t *testing.T) {
	t.Run("simple write", test.Case(func(ctx context.Context) {
		os.Remove("/tmp/test.log")
		writer, err := bytime.NewRotationWriter(func(writer *bytime.Writer) {
			writer.WritePath = "/tmp/test.log"
		})
		must.Nil(err)
		writer.Write([]byte("hello"))
		content, err := ioutil.ReadFile("/tmp/test.log")
		must.Nil(err)
		must.Equal("hello", string(content))
	}))
}