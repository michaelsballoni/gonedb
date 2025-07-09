package test

import (
	"fmt"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestParseCmds(t *testing.T) {
	validate_cmd("", []string{})
	validate_cmd(" ", []string{})
	validate_cmd("  ", []string{})
	validate_cmd("   ", []string{})

	validate_cmd("foo", []string{"foo"})
	validate_cmd("foo ", []string{"foo"})
	validate_cmd(" foo", []string{"foo"})
	validate_cmd(" foo ", []string{"foo"})

	validate_cmd("foo bar", []string{"foo", "bar"})
	validate_cmd("foo  bar", []string{"foo", "bar"})
	validate_cmd("foo   bar", []string{"foo", "bar"})

	validate_cmd("foo bar blet", []string{"foo", "bar", "blet"})
	validate_cmd("foo  bar  blet", []string{"foo", "bar", "blet"})

	validate_cmd("\"\"", []string{""})
	validate_cmd("\" \"", []string{" "})
	validate_cmd("\"fred\"", []string{"fred"})

	validate_cmd("\"fred some\"", []string{"fred some"})
	validate_cmd("\"fred some\" blet", []string{"fred some", "blet"})
	validate_cmd("\"fred some\" \"blet\"", []string{"fred some", "blet"})
	validate_cmd("\"fred some\" \"blet\" \"some monkey\"", []string{"fred some", "blet", "some monkey"})
}

func validate_cmd(cmd string, matches []string) {
	fmt.Printf("validate_cmd: %s\n", cmd)
	gotten := gonedb.ParseCmds(cmd)
	fmt.Printf("gotten: %d vs. matches: %d\n", len(gotten), len(matches))
	if len(matches) != len(gotten) {
		panic("len mismatch")
	}

	for i := range matches {
		fmt.Printf("gotten[%d]: %s vs. matches[%d]: %s\n", i, gotten[i], i, matches[i])
		if gotten[i] != matches[i] {
			panic("value mismatch")
		}
	}

	fmt.Printf("validate_cmd: %s - PASS\n", cmd)
}
