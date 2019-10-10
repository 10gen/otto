package pcre

import (
	"github.com/dlclark/regexp2"
)

// Regexp represents a regular expression object that uses PCRE
type Regexp struct {
	regexp *regexp2.Regexp
}

// New creates a new PCRE Regexp object
func New(pattern string, opt regexp2.RegexOptions) (Regexp, error) {
	re, err := regexp2.Compile(pattern, opt)
	if err != nil {
		return Regexp{}, err
	}
	regexp := Regexp{
		regexp: re,
	}
	return regexp, nil
}

// FindStringIndex searches for a match given a string and returns the index pair of first occurrence
func (regexp *Regexp) FindStringIndex(input string) ([]int, error) {
	match, err := regexp.regexp.FindStringMatch(input)
	if err != nil {
		return nil, err
	}

	startIndex := match.Group.Capture.Index
	endIndex := match.Group.Capture.Index + match.Group.Capture.Length
	return []int{startIndex, endIndex}, nil
}

// FindAllStringIndex returns the index pairs of all successive occurrences of the regex in the input, from the given index
func (regexp *Regexp) FindAllStringIndex(input string, startAt int) ([][]int, error) {
	var results [][]int
	match, err := regexp.regexp.FindStringMatchStartingAt(input, startAt)
	if err != nil {
		return nil, err
	}

	for match != nil {
		startIndex := match.Group.Capture.Index
		endIndex := match.Group.Capture.Index + match.Group.Capture.Length
		results = append(results, []int{startIndex, endIndex})
		match, err = regexp.regexp.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

// FindAllStringSubmatch returns the index all successive matches by match group
func (regexp *Regexp) FindAllStringSubmatch(input string, startAt int) ([][]string, error) {
	var results [][]string
	match, err := regexp.regexp.FindStringMatchStartingAt(input, startAt)
	if err != nil {
		return nil, err
	}

	for match != nil {
		var submatchResults []string
		for _, submatch := range match.Groups() {
			submatchResults = append(submatchResults, submatch.String())
		}
		results = append(results, submatchResults)
		match, err = regexp.regexp.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

// FindAllStringSubmatchIndex returns the index all successive input pairs by match group
func (regexp *Regexp) FindAllStringSubmatchIndex(input string, startAt int) ([][]int, error) {
	var results [][]int
	match, err := regexp.regexp.FindStringMatchStartingAt(input, startAt)
	if err != nil {
		return nil, err
	}

	for match != nil {
		var submatchResults []int
		for _, submatch := range match.Groups() {
			startIndex := submatch.Index
			endIndex := submatch.Index + submatch.Length
			submatchResults = append(submatchResults, startIndex, endIndex)
		}
		results = append(results, submatchResults)
		match, err = regexp.regexp.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}
