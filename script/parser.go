package script

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// parses a expr(args) form.
var exprRegex = regexp.MustCompile(`(\w+)\((.+)\)`)

func splitOpAndArgs(s string) (string, string, error) {
	// remove non printable characters
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)

	matches := exprRegex.FindStringSubmatch(s)
	if matches == nil {
		return "", "", errors.New("invalid script")
	}

	return matches[1], matches[2], nil

}

func parseScript(s string) (ScriptExpr, error) {
	return parseScriptR(s, true)
}

func parseScriptR(s string, topLevel bool) (ScriptExpr, error) {
	op, args, err := splitOpAndArgs(s)
	if err != nil {
		return nil, err
	}

	switch op {
	case "sh":
		if !topLevel {
			return nil, errors.New("sh must be a top-level expression")
		}

		script, err := parseScriptR(args, false)
		if err != nil {
			return nil, err
		}

		return Sh(script), nil
	case "wsh":
		script, err := parseScriptR(args, false)
		if err != nil {
			return nil, err
		}

		return Wsh(script), nil
	case "pkh":
		return Pkh(args), nil
	case "wpkh":
		return Wpkh(args), nil
	case "multi", "sortedmulti":
		n, keys, err := parseMultiArgs(args)
		if err != nil {
			return nil, err
		}

		if op == "multi" {
			return Multi(n, keys...), nil
		} else {
			return Sortedmulti(n, keys...), nil
		}
	}

	return nil, fmt.Errorf("invalid op '%s'", op)
}

// parseMultiArgs parsers a string with the form `N,<key1,key2...keyM>`
func parseMultiArgs(args string) (int, []string, error) {
	split := strings.Split(args, ",")
	if len(split) < 2 {
		return 0, nil, fmt.Errorf("invalid multi() argument")
	}

	n, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, nil, err
	}

	return n, split[1:], nil
}
