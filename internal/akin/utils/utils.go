package utils

import (
	"os"
	"path/filepath"
)

func MakeSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, s := range items {
		set[s] = true
	}
	return set
}

func MergeStrArr(a, b []string) []string {
	check := make(map[string]int)
	c := append(a, b...)
	res := make([]string, 0)
	for _, val := range c {
		check[val] = 1
	}
	for key := range check {
		res = append(res, key)
	}
	return res
}

func AbsoluteDirPath(paths ...string) string {
	abs, _ := os.Getwd()
	for _, path := range paths {
		abs = filepath.Join(abs, path)
	}
	return filepath.Dir(abs)
}
