package regexp

// Regexp represents a type that can parse regular expressions and perform finds using them
type Regexp interface {
	FindStringIndex(input string) ([]int, error)
	FindStringSubmatchIndex(input string) ([]int, error)
	FindAllStringIndex(input string, startAt int) ([][]int, error)
	FindAllSubmatchIndex(input []byte, startAt int) ([][]int, error)
	FindAllStringSubmatch(input string, startAt int) ([][]string, error)
	FindAllStringSubmatchIndex(input string, startAt int) ([][]int, error)
}
