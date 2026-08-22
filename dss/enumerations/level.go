// Ported from dss-enumerations/.../Level.java (DSS 6.5.RC1).
package enumerations

// Level defines a Validation Policy constraint execution level.
type Level string

const (
	// LevelFail stops the validation process and reports as error.
	LevelFail Level = "FAIL"
	// LevelWarn continues the validation process and adds a warning
	// message.
	LevelWarn Level = "WARN"
	// LevelInform continues the validation process and adds an
	// informative message.
	LevelInform Level = "INFORM"
	// LevelIgnore continues the validation process and skips the current
	// check (equals to not present check).
	LevelIgnore Level = "IGNORE"
)

// LevelValues returns all constants in declaration order.
func LevelValues() []Level {
	return []Level{
		LevelFail,
		LevelWarn,
		LevelInform,
		LevelIgnore,
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
