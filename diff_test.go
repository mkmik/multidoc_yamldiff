package main

import (
	"testing"
)

func TestContext(t *testing.T) {
	src := ` a
 b
 c
+d
 e
-f
 g
 h
 i
 l
`

	want := `@@ .. @@
 b
 c
+d
 e
-f
 g
 h
`

	if got := compact(src, 2); got != want {
		t.Errorf("got:\n%s, want:\n%s", got, want)
	}

}
