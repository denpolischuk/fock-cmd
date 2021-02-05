package bookmark_test

import (
	"fmt"
	"testing"

	"github.com/denpolischuk/fock-cli/internal/app/shared/types/bookmark"
)

func spamMapWithDummyData(m map[string]string, l int, s string) {
	for i := 0; i < l; i++ {
		m[fmt.Sprintf("1test-%d", i)] = s
	}
}

func TestBookmarkList(t *testing.T) {
	t.Run("List.Add", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			list := bookmark.NewList()
			if err := list.Add("test", "http://localhost:3000"); err != nil && len(list.List) == 1 {
				t.Fail()
			}
		})

		t.Run("fail - url doesn't match", func(t *testing.T) {
			list := bookmark.NewList()
			if err := list.Add("test", "bad url"); err == nil || len(list.List) == 1 {
				t.Fail()
			}
		})

		t.Run("fail - alias length is too short", func(t *testing.T) {
			list := bookmark.NewList()
			if err := list.Add("t", "http://localhost:3000"); err == nil || len(list.List) == 1 {
				t.Fail()
			}
		})

		t.Run("fail - list is too big", func(t *testing.T) {
			list := bookmark.NewList()
			spamMapWithDummyData(list.List, bookmark.BookmarksCap, "http://localhost:3000")
			if err := list.Add("test", "http://localhost:3000"); err == nil || len(list.List) == bookmark.BookmarksCap+1 {
				t.Fail()
			}
		})
	})
}
