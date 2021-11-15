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

func parseScript(s, path string) (ScriptExpr, error) {
	return parseScriptR(s, path, true)
}

func deriveIfXpub(s, path string) (string, error) {
	// Remove the [hex/path] origin if present.
	s = trimKeyOrigin(s)

	if !IsXPub(s) {
		return s, nil
	}

	xpub, err := NewXPub(s)
	if err != nil {
		return "", err
	}

	if path != "" {
		xpub, err = xpub.Derive(path)
		if err != nil {
			return "", err
		}
	}

	pub, err := xpub.PubKey()
	if err != nil {
		return "", nil
	}

	return pub, nil
}

func parseScriptR(s, path string, topLevel bool) (ScriptExpr, error) {
	op, args, err := splitOpAndArgs(s)
	if err != nil {
		return nil, err
	}

	switch op {
	case "sh":
		if !topLevel {
			return nil, errors.New("sh must be a top-level expression")
		}

		script, err := parseScriptR(args, path, false)
		if err != nil {
			return nil, err
		}

		return Sh(script), nil
	case "wsh":
		script, err := parseScriptR(args, path, false)
		if err != nil {
			return nil, err
		}

		return Wsh(script), nil
	case "pkh":
		der, err := deriveIfXpub(args, path)
		if err != nil {
			return nil, err
		}

		return Pkh(der), nil
	case "wpkh":
		der, err := deriveIfXpub(args, path)
		if err != nil {
			return nil, err
		}

		return Wpkh(der), nil
	case "multi", "sortedmulti":
		n, keys, err := parseMultiArgs(args, path)
		if err != nil {
			return nil, err
		}

		if op == "multi" {
			return Multi(n, keys...), nil
		}
		return Sortedmulti(n, keys...), nil
	case "tr":
		if !topLevel {
			return nil, errors.New("tr() must be a top-level expression")
		}

		var key, tree string
		split := strings.Split(args, ",")
		switch len(split) {
		case 1:
			key = args
		case 2:
			key = split[0]
			tree = split[1]
		default:
			return nil, errors.New("too many arguments for tr()")
		}

		return Tr(key, tree), nil
	}

	return nil, fmt.Errorf("invalid op '%s'", op)
}

// parseMultiArgs parsers a string with the form `N,<key1,key2...keyM>`
func parseMultiArgs(args, path string) (int, []string, error) {
	split := strings.Split(args, ",")
	if len(split) < 2 {
		return 0, nil, fmt.Errorf("invalid multi() argument")
	}

	n, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, nil, err
	}

	var keys []string
	for _, key := range split[1:] {
		der, err := deriveIfXpub(key, path)
		if err != nil {
			return 0, nil, err
		}
		keys = append(keys, der)
	}

	return n, keys, nil
}

func Parse(s string) (*Script, error) {
	return ParseWithPath(s, "")
}

func ParseWithPath(s string, path string) (*Script, error) {
	script, err := parseScript(s, path)
	if err != nil {
		return nil, err
	}
	return script.Eval()
}
