package validators

// Validate is a struct the holds the configuration for a validator.
type Validate struct {
	Required  bool          `yaml:"required"`
	MinLength int           `yaml:"min"`
	MaxLength int           `yaml:"max"`
	Match     ValidateMatch `yaml:"match"`
}

type ValidateMatch struct {
	Regex   string `yaml:"regex"`
	Message string `yaml:"message"`
}

// GetValidatorFuncs converts a Validate struct into a slice of validator functions.
func GetValidatorFuncs[T Validatable](v Validate) []Validator[T] {
	var vals []Validator[T]

	if v.MinLength > 0 {
		vals = append(vals, MinLength[T](v.MinLength))
	} else if v.Required {
		vals = append(vals, NotZero[T])
	}

	if v.MaxLength > 0 {
		vals = append(vals, MaxLength[T](v.MaxLength))
	}

	if v.Match.Regex != "" {
		vals = append(vals, Match[T](v.Match.Regex, v.Match.Message))
	}

	return vals
}
