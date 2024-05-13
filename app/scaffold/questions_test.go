package scaffold

import (
	"reflect"
	"testing"
)

func Test_QuestionGroupBy(t *testing.T) {
	type tcase struct {
		name   string
		input  []Question
		expect [][]Question
	}

	tests := []tcase{
		{
			name: "no groups",
			input: []Question{
				{Name: "name1", Group: "", Prompt: AnyPrompt{}},
				{Name: "name2", Group: "", Prompt: AnyPrompt{}},
			},
			expect: [][]Question{
				{{Name: "name1", Group: "", Prompt: AnyPrompt{}}},
				{{Name: "name2", Group: "", Prompt: AnyPrompt{}}},
			},
		},
		{
			name: "one group",
			input: []Question{
				{Name: "name1", Group: "group1", Prompt: AnyPrompt{}},
				{Name: "name2", Group: "group1", Prompt: AnyPrompt{}},
			},
			expect: [][]Question{
				{
					{Name: "name1", Group: "group1", Prompt: AnyPrompt{}},
					{Name: "name2", Group: "group1", Prompt: AnyPrompt{}},
				},
			},
		},
		{
			name: "two groups, and one question without group",
			input: []Question{
				{Name: "name1", Group: "", Prompt: AnyPrompt{}},
				{Name: "name2", Group: "group1", Prompt: AnyPrompt{}},
				{Name: "name3", Group: "group1", Prompt: AnyPrompt{}},
				{Name: "name4", Group: "group2", Prompt: AnyPrompt{}},
				{Name: "name5", Group: "group2", Prompt: AnyPrompt{}},
				{Name: "name6", Group: "", Prompt: AnyPrompt{}},
			},
			expect: [][]Question{
				{{Name: "name1", Group: "", Prompt: AnyPrompt{}}},
				{
					{Name: "name2", Group: "group1", Prompt: AnyPrompt{}},
					{Name: "name3", Group: "group1", Prompt: AnyPrompt{}},
				},
				{
					{Name: "name4", Group: "group2", Prompt: AnyPrompt{}},
					{Name: "name5", Group: "group2", Prompt: AnyPrompt{}},
				},
				{{Name: "name6", Group: "", Prompt: AnyPrompt{}}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := QuestionGroupBy(tc.input)
			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("expected %v, got %v", tc.expect, got)
			}
		})
	}
}
