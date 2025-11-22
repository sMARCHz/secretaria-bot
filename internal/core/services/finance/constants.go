package finance

var commandPrefixSet = map[string]struct{}{
	"!p":        {},
	"!e":        {},
	"!t":        {},
	"balance":   {},
	"statement": {},
}

const (
	invalidCommandMsg = "Invalid command"
)
