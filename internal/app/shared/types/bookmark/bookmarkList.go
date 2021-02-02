package bookmark

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

const bookamrksCap = 256

var re, _ = regexp.Compile(`^[A-z0-9\-\_]{3,32}$`)

// List - bookmarks container
type List struct {
	List map[string]string `json:"list"`
}

// Add - appends new bookmark to the list
func (bl *List) Add(alias string, URL string) error {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return fmt.Errorf("Invalid URL: %s", err.Error())
	}
	if !re.MatchString(alias) {
		return errors.New("Invalid alias name. It should only contain letters, numbers, - and _ signs, having min length = 3 and max length = 32")
	}

	bl.List[alias] = URL

	return nil
}

// NewList - returns new list
func NewList() *List {
	return &List{
		List: make(map[string]string),
	}
}
