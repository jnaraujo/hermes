package data

import (
	"errors"
	"fmt"
	"strings"
)

type Data struct {
	Cmd     string
	Content string
}

func New(cmd, content string) *Data {
	return &Data{cmd, content}
}

func FromStr(raw string) (*Data, error) {
	lines := strings.Split(strings.Replace(raw, `\n`, "\n", -1), "\n")

	if len(lines) <= 1 {
		return nil, errors.New("invalid data")
	}

	return New(lines[0], lines[1]), nil
}

func (d *Data) Decode() string {
	return fmt.Sprintf("%s\n%s", d.Cmd, d.Content)
}
