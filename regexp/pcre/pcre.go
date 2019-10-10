package pcre

import (
	"fmt"

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
func (regexp Regexp) FindStringIndex(input string) ([]int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in r: %+v", r)
		}
	}()

	match, err := regexp.regexp.FindStringMatch(input)
	if err != nil {
		return nil, err
	}

	startIndex := match.Group.Capture.Index
	endIndex := match.Group.Capture.Index + match.Group.Capture.Length
	return []int{startIndex, endIndex}, nil
}

// FindStringSubmatchIndex finds the indexes of the first match and group matches
func (regexp Regexp) FindStringSubmatchIndex(input string) ([]int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in r: %+v", r)
		}
	}()

	match, err := regexp.regexp.FindStringMatch(input)
	if err != nil {
		return nil, err
	}

	var submatchResults []int
	for _, submatch := range match.Groups() {
		startIndex := submatch.Index
		endIndex := submatch.Index + submatch.Length
		submatchResults = append(submatchResults, startIndex, endIndex)
	}
	return submatchResults, nil
}

// FindAllStringIndex returns the index pairs of all successive occurrences of the regex in the input, from the given index
func (regexp Regexp) FindAllStringIndex(input string, index int) ([][]int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in r: %+v", r)
		}
	}()

	var results [][]int
	match, err := regexp.regexp.FindStringMatch(input)
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
	if index >= 0 {
		return results[:index], nil
	}

	return results, nil
}

// FindAllSubmatchIndex returns all successive input pairs by match group
func (regexp Regexp) FindAllSubmatchIndex(input []byte, index int) ([][]int, error) {
	return regexp.FindAllStringSubmatchIndex(string(input), index)
}

// FindAllStringSubmatch returns the index all successive matches by match group
func (regexp Regexp) FindAllStringSubmatch(input string, index int) ([][]string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in r: %+v", r)
		}
	}()

	var results [][]string
	match, err := regexp.regexp.FindStringMatch(input)
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
	if index >= 0 {
		return results[:index], nil
	}

	return results, nil
}

// FindAllStringSubmatchIndex returns all successive input pairs by match group
func (regexp Regexp) FindAllStringSubmatchIndex(input string, index int) ([][]int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in r: %+v", r)
		}
	}()

	var results [][]int
	match, err := regexp.regexp.FindStringMatch(input)
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
	if index >= 0 {
		return results[:index], nil
	}

	return results, nil
}
