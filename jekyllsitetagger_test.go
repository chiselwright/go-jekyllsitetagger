package jekyllsitetagger_test

import (
	"reflect"
	"testing"

	"github.com/chiselwright/go-jekyllsitetagger"
)

func TestParseTagLine(t *testing.T) {
	t.Parallel()

	type testCase struct {
		line string
		want []string
	}

	testCases := []testCase{
		{
			line: "tags: one two three",
			want: []string{"one", "two", "three"},
		},
		{
			line: "tags: one         two three",
			want: []string{"one", "two", "three"},
		},
		{
			line: "tags: one two one two one three",
			want: []string{"one", "two", "three"},
		},
		{
			line: "tags: one:two one:two:one three",
			want: []string{"one:two", "one:two:one", "three"},
		},
	}

	for _, test := range testCases {
		got := jekyllsitetagger.ParseTagLine(test.line)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got %+v, want %+v", got, test.want)
		}
	}
}

func TestUnique(t *testing.T) {
	t.Parallel()

	type testCase struct {
		tags []string
		want []string
	}

	testCases := []testCase{
		{
			tags: []string{},
			want: []string{},
		},
		{
			tags: []string{"one"},
			want: []string{"one"},
		},
		{
			tags: []string{"one", "one"},
			want: []string{"one"},
		},
		{
			tags: []string{"one", "two", "one"},
			want: []string{"one", "two"},
		},
	}

	for _, test := range testCases {
		got := jekyllsitetagger.Unique(test.tags)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got %+v, want %+v", got, test.want)
		}
	}
}
