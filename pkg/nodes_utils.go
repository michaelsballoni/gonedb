package gonedb

import (
	"errors"
	"strconv"
	"strings"
)

// The API you interact with
type utils struct{}

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
	sep := "/"
	if str == "" || str == sep {
		return []int64{}, nil
	}

	strs := strings.Split(str, sep)
	// DEBUG
	//fmt.Printf("strings.Split(str, sep): str:%s sep:%s strs:%d\n", str, sep, len(strs))
	ids := make([]int64, len(strs))
	for i, v := range strs {
		n, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return []int64{}, e
		} else {
			ids[i] = n
		}
	}
	return ids, nil
}

// Convert IDs into an ID path string
func (u *utils) IdsToParentsStr(ids []int64) string {
	if len(ids) == 0 || (len(ids) == 1 && ids[0] == 0) {
		return ""
	}

	strs := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != 0 {
			strs = append(strs, strconv.FormatInt(id, 10))
		}
	}
	return strings.Join(strs, string('/')) + "/"
}
