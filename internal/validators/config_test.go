package validators

import "testing"

func Test_GetValidatorsFuncs(t *testing.T) {
	type tsubcase struct {
		input   string
		wantErr bool
	}

	type tcase struct {
		name    string
		cfg     Validate
		expects []tsubcase
	}

	cases := []tcase{
		{
			name: "no validation",
			cfg:  Validate{},
			expects: []tsubcase{
				{input: "", wantErr: false},
				{input: "test", wantErr: false},
			},
		},
		{
			name: "required",
			cfg:  Validate{Required: true},
			expects: []tsubcase{
				{input: "", wantErr: true},
				{input: "test", wantErr: false},
			},
		},
		{
			name: "min length",
			cfg:  Validate{MinLength: 3},
			expects: []tsubcase{
				{input: "", wantErr: true},
				{input: "te", wantErr: true},
				{input: "test", wantErr: false},
			},
		},
		{
			name: "max length",
			cfg:  Validate{MaxLength: 3},
			expects: []tsubcase{
				{input: "test", wantErr: true},
				{input: "te", wantErr: false},
				{input: "", wantErr: false},
			},
		},
		{
			name: "match regex",
			cfg: Validate{
				Match: ValidateMatch{
					Regex:   "^[a-z]+$",
					Message: "must be lowercase letters only",
				},
			},
			expects: []tsubcase{
				{input: "TEST", wantErr: true},
				{input: "test", wantErr: false},
				{input: "test123", wantErr: true},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vs := Combine(GetValidatorFuncs[string](tc.cfg)...)

			for _, sub := range tc.expects {
				t.Run(sub.input, func(t *testing.T) {
					got := vs(sub.input)

					if sub.wantErr && got == nil {
						t.Errorf("expected error, got nil")
					}

					if !sub.wantErr && got != nil {
						t.Errorf("expected no error, got %v", got)
					}
				})
			}
		})
	}
}
