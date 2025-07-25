package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseVersion(ver string) ([]int, error) {
	parts := strings.Split(ver, ".")
	res := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("uncorrect version: %s", ver)
		}
		res[i] = n
	}
	return res, nil
}

func compareVersions(a, b string) int {
	va, _ := parseVersion(a)
	vb, _ := parseVersion(b)
	for i := 0; i < len(va) || i < len(vb); i++ {
		var ai, bi int
		if i < len(va) {
			ai = va[i]
		}
		if i < len(vb) {
			bi = vb[i]
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
}

func versionMatches(rule, actual string) bool {
	rule = strings.TrimSpace(rule)
	switch {
	case strings.HasPrefix(rule, ">="):
		return compareVersions(actual, rule[2:]) >= 0
	case strings.HasPrefix(rule, "<="):
		return compareVersions(actual, rule[2:]) <= 0
	case strings.HasPrefix(rule, ">"):
		return compareVersions(actual, rule[1:]) > 0
	case strings.HasPrefix(rule, "<"):
		return compareVersions(actual, rule[1:]) < 0
	case strings.HasPrefix(rule, "="):
		return compareVersions(actual, rule[1:]) == 0
	default:
		return compareVersions(actual, rule) == 0
	}
}
