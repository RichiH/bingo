package gobin

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/bwplotka/gobin/pkg/testutil"
)

func TestParse(t *testing.T) {
	testPackages := []string{
		"github.com/bwplotka/gobin",
		"github.com/bwplotka/gobin/a",
		"github.com/bwplotka/gobin/a/b",
		"github.com/bwplotka/gobin2/and/yolo",
		"github.com/bwplotka/gobin2/yolo",
		"github.com/bwplotka/gobin3",
		"golang.org/x/tools/cmd/goimports",
	}
	for _, tcase := range []struct {
		input io.Reader

		expected    []string
		expectedErr error
	}{
		{
			input:       strings.NewReader(""),
			expectedErr: errors.New("read /test/something/.gobin/binaries.go: binaries.go:1:1: expected 'package', found 'EOF'"),
		},
		{
			input: func() io.Reader {
				b := bytes.Buffer{}
				testutil.Ok(t, Recreate("test/something/.gobin/binaries.go", &b, testPackages))
				return &b
			}(),
			expected: testPackages,
		},
	} {
		t.Run("", func(t *testing.T) {
			bs, err := Parse("/test/something/.gobin/binaries.go", tcase.input)
			if tcase.expectedErr != nil {
				testutil.NotOk(t, err)
				testutil.Equals(t, tcase.expectedErr.Error(), err.Error())
				return
			}

			testutil.Ok(t, err)
			testutil.Equals(t, tcase.expected, bs)
		})
	}
}

func TestRecreate(t *testing.T) {
	t.Run("empty moddir", func(t *testing.T) {
		err := Recreate(".", nil, []string{"github.com/bwplotka/gobin"})
		testutil.NotOk(t, err)
		testutil.Equals(t, "moddir cannot be empty, got filePath: .", err.Error())
	})
	t.Run("ok", func(t *testing.T) {
		b := &strings.Builder{}
		testutil.Ok(t, Recreate("/test/something/.gobin/binaries.go", b, []string{
			"github.com/bwplotka/gobin2/yolo",
			"github.com/bwplotka/gobin3",
			"github.com/bwplotka/gobin2/and/yolo",
			"github.com/bwplotka/gobin/a",
			"github.com/bwplotka/gobin",
			"golang.org/x/tools/cmd/goimports",
			"github.com/bwplotka/gobin/a/b",
		}))
		testutil.Equals(t, `// Code generated by https://github.com/bwplotka/gobin . DO NOT EDIT.
// NOTE: Actually you can edit just fine, just don't be surprised if the file will be rewritten at some point by 
// gobin tool.
//
// This file is and extended version of somethings that is called "tools.go package" described here: https://github.com/golang/go/issues/25922#issuecomment-590529870
// It allows go modules to maintain certain version of binaries you or your project use.
// Main extension that bwplotka/gobin adds, is that the file has to be stored as the separate go module file, which allows 
// separation of dev tools from critical, production code which is th key.
//
// Read more on https://github.com/bwplotka/gobin .

// +build tools
package gobin

import (
	_ "github.com/bwplotka/gobin"
	_ "github.com/bwplotka/gobin/a"
	_ "github.com/bwplotka/gobin/a/b"
	_ "github.com/bwplotka/gobin2/and/yolo"
	_ "github.com/bwplotka/gobin2/yolo"
	_ "github.com/bwplotka/gobin3"
	_ "golang.org/x/tools/cmd/goimports"
)

const (
	// FileVersion represents version of this file. This will tell future versions of gobin how to parse this file.
	FileVersion = "v1.0.0"
	// Gobin version used to generate this file. Used for debugging only.
	GobinVersion = "??"
)
`, b.String(), "does not match, got %v", b.String())
	})
}
