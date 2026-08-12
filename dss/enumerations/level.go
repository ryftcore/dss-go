// Ported from dss-enumerations/.../Level.java (DSS 6.5.RC1).
package enumerations

// Level defines a Validation Policy constraint execution level.
type Level string

const (
	// Level_FAIL stops the validation process and reports as error.
	Level_FAIL Level = "FAIL"
	// Level_WARN continues the validation process and adds a warning
	// message.
	Level_WARN Level = "WARN"
	// Level_INFORM continues the validation process and adds an
	// informative message.
	Level_INFORM Level = "INFORM"
	// Level_IGNORE continues the validation process and skips the current
	// check (equals to not present check).
	Level_IGNORE Level = "IGNORE"
)

// LevelValues returns all constants in declaration order.
func LevelValues() []Level {
	return []Level{
		Level_FAIL,
		Level_WARN,
		Level_INFORM,
		Level_IGNORE,
	}
}

// LevelValueOf returns the constant matching the given Java enum name.
func LevelValueOf(name string) (Level, error) {
	for _, v := range LevelValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &levelInvalidValueError{name}
}

type levelInvalidValueError struct {
	name string
}

func (e *levelInvalidValueError) Error() string {
	return "no enum constant Level." + e.name
}
