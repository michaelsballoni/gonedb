package gonedb

import (
	"errors"
	"strconv"
	"strings"
)

type utils struct{}

// The API you interact with
var NodeUtils utils

// Turn IDs into a string used in a SELECT IN (here)
// Errors if an empty input is provided
func (u *utils) IdsToSqlIn(ids []int64) (string, error) {
	if len(ids) == 0 {
		return "", errors.New("ids is empty")
	}

	var output strings.Builder
	for _, id := range ids {
		if output.Len() > 0 {
			output.WriteRune(',')
		}
		output.WriteString(strconv.FormatInt(id, 10))
	}
	return output.String(), nil
}

// Turn a separator-delimited string into a slice of IDs
// only IDs found between separators are returned
// nothing in, nothing out
func (u *utils) StringToIds(str string) ([]int64, error) {
	separator := '/'
	ids := make([]int64, 0)
	var collector strings.Builder
	for _, c := range str {
		if c == separator {
			if collector.Len() > 0 {
				v, e := strconv.ParseInt(collector.String(), 10, 64)
				if e != nil {
					return []int64{}, e
				}
				ids = append(ids, v)
				collector.Reset()
			}
		} else {
			collector.WriteRune(rune(c))
		}
	}
	if collector.Len() > 0 {
		v, e := strconv.ParseInt(collector.String(), 10, 64)
		if e != nil {
			return []int64{}, e
		}
		ids = append(ids, v)
		collector.Reset()
	}
	return ids, nil
}

// Convert IDs into an ID path string
func (u *utils) IdsToParentsStr(ids []int64) string {
	var builder strings.Builder
	for _, id := range ids {
		if id != 0 {
			builder.WriteString(strconv.FormatInt(id, 10))
			builder.WriteRune('/')
		}
	}
	return builder.String()
}
