package scaffoldrc

type RunHooksOption string

var (
	RunHooksNever  RunHooksOption = "never"
	RunHooksAlways RunHooksOption = "always"
	RunHooksPrompt RunHooksOption = "prompt"
)

func ParseRunHooksOption(s string) RunHooksOption {
	zero := RunHooksOption("")
	ptr := &zero
	_ = ptr.UnmarshalText([]byte(s))
	return *ptr
}

func (r *RunHooksOption) UnmarshalText(text []byte) error {
	switch string(text) {
	case "never", "no", "false":
		*r = RunHooksNever
	case "always", "yes", "true":
		*r = RunHooksAlways
	case "prompt", "": // if left empty, default to prompt
		*r = RunHooksPrompt
	default:
		// fallback to whatever they input so we can log the incorrect value
		*r = RunHooksOption(string(text))
	}

	return nil
}

func (r RunHooksOption) IsValid() bool {
	switch r {
	case RunHooksNever, RunHooksAlways, RunHooksPrompt:
		return true
	default:
		return false
	}
}

func (r RunHooksOption) String() string {
	return string(r)
}
