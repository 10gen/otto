package re2

import (
	"regexp"
)

// Regexp represents a regular expression object that uses re2
type Regexp struct {
	regexp *regexp.Regexp
}

// New creates a new re2 Regexp object
func New(pattern string) (Regexp, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Regexp{}, err
	}
	regexp := Regexp{
		regexp: re,
	}
	return regexp, nil
}

// FindStringIndex searches for a match given a string and returns the index pair of first occurrence
func (regexp Regexp) FindStringIndex(input string) ([]int, error) {
	return regexp.regexp.FindStringIndex(input), nil
}

// FindStringSubmatchIndex finds the indexes of the first match and group matches
func (regexp Regexp) FindStringSubmatchIndex(input string) ([]int, error) {
	return regexp.regexp.FindStringSubmatchIndex(input), nil
}

// FindAllStringIndex returns the index pairs of all successive occurrences of the regex in the input, from the given index
func (regexp Regexp) FindAllStringIndex(input string, startAt int) ([][]int, error) {
	return regexp.regexp.FindAllStringIndex(input, startAt), nil
}

// FindAllSubmatchIndex returns all successive input pairs by match group
func (regexp Regexp) FindAllSubmatchIndex(input []byte, startAt int) ([][]int, error) {
	return regexp.regexp.FindAllSubmatchIndex(input, startAt), nil
}

// FindAllStringSubmatch returns the index all successive matches by match group
func (regexp Regexp) FindAllStringSubmatch(input string, startAt int) ([][]string, error) {
	return regexp.regexp.FindAllStringSubmatch(input, startAt), nil
}

// FindAllStringSubmatchIndex returns the index all successive input pairs by match group
func (regexp Regexp) FindAllStringSubmatchIndex(input string, startAt int) ([][]int, error) {
	return regexp.regexp.FindAllStringSubmatchIndex(input, startAt), nil
}
