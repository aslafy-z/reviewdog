package reviewdog

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/reviewdog/errorformat"
)

var _ CommentService = &testWriter{}

type testWriter struct {
	FakePost func(c *Comment) error
}

func (s *testWriter) Post(_ context.Context, c *Comment) error {
	return s.FakePost(c)
}

func ExampleReviewdog() {
	difftext := `diff --git a/golint.old.go b/golint.new.go
index 34cacb9..a727dd3 100644
--- a/golint.old.go
+++ b/golint.new.go
@@ -2,6 +2,12 @@ package test
 
 var V int
 
+var NewError1 int
+
 // invalid func comment
 func F() {
 }
+
+// invalid func comment2
+func F2() {
+}
`
	lintresult := `golint.new.go:3:5: exported var V should have comment or be unexported
golint.new.go:5:5: exported var NewError1 should have comment or be unexported
golint.new.go:7:1: comment on exported function F should be of the form "F ..."
golint.new.go:11:1: comment on exported function F2 should be of the form "F2 ..."
`
	efm, _ := errorformat.NewErrorformat([]string{`%f:%l:%c: %m`})
	p := NewErrorformatParser(efm)
	c := NewRawCommentWriter(os.Stdout)
	d := NewDiffString(difftext, 1)
	app := NewReviewdog("tool name", p, c, d)
	app.Run(context.Background(), strings.NewReader(lintresult))
	// Unordered output:
	// golint.new.go:5:5: exported var NewError1 should have comment or be unexported
	// golint.new.go:11:1: comment on exported function F2 should be of the form "F2 ..."
}

func TestReviewdog_Run_clean_path(t *testing.T) {
	difftext := `diff --git a/golint.old.go b/golint.new.go
index 34cacb9..a727dd3 100644
--- a/golint.old.go
+++ b/golint.new.go
@@ -2,6 +2,12 @@ package test
 
 var V int
 
+var NewError1 int
+
 // invalid func comment
 func F() {
 }
+
+// invalid func comment2
+func F2() {
+}
`
	lintresult := `./golint.new.go:3:5: exported var V should have comment or be unexported
./golint.new.go:5:5: exported var NewError1 should have comment or be unexported
./golint.new.go:7:1: comment on exported function F should be of the form "F ..."
./golint.new.go:11:1: comment on exported function F2 should be of the form "F2 ..."
`

	want := "golint.new.go"

	c := &testWriter{
		FakePost: func(c *Comment) error {
			if got := c.Path; got != want {
				t.Errorf("path: got %v, want %v", got, want)
			}
			return nil
		},
	}

	efm, _ := errorformat.NewErrorformat([]string{`%f:%l:%c: %m`})
	p := NewErrorformatParser(efm)
	d := NewDiffString(difftext, 1)
	app := NewReviewdog("tool name", p, c, d)
	app.Run(context.Background(), strings.NewReader(lintresult))
}
