package gonedb

import (
	"errors"
	"fmt"
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

// Get the ID of the node at the end of an ID path
// Just a little string math
func (u *utils) GetLasttPathId(path string) (int64, error) {
	if path == "" || path == "/" {
		return 0, nil
	}

	start_idx := len(path) - 1
	if path[start_idx] == '/' {
		start_idx -= 1
	}
	end_idx := start_idx + 1

	for start_idx >= 0 && path[start_idx] != '/' {
		start_idx -= 1
	}
	if start_idx < 0 {
		return -1, fmt.Errorf("path does not have last /: %s", path)
	}

	last_node_id_str := path[start_idx+1 : end_idx]
	if len(last_node_id_str) == 0 {
		return -1, fmt.Errorf("path does not have last ID value: %s", path)
	}

	var last_node_id int64
	var err error
	last_node_id, err = strconv.ParseInt(last_node_id_str, 10, 64)
	if err != nil {
		return -1, err
	} else {
		return last_node_id, err
	}
}
