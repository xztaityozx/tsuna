package watanabe

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Newできるべき"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := New()
			assert.NotNil(t, actual)
			assert.NotNil(t, actual.Logger)
			assert.Regexp(t, "渡(邉|辺|邊)", actual.Last)
			assert.NotEmpty(t, actual.Yomi)
			assert.NotEmpty(t, actual.First)
		})
	}
}

func TestWatanabe_FullName(t *testing.T) {
	type fields struct {
		Logger *zerolog.Logger
		Last   string
		First  string
		yomi   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{nil, "Last", "First", "Yomi"}, "LastFirst"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Watanabe{
				Logger: tt.fields.Logger,
				Last:   tt.fields.Last,
				First:  tt.fields.First,
				Yomi:   tt.fields.yomi,
			}
			assert.Equal(t, tt.want, w.FullName())
		})
	}
}

func TestWatanabe_String(t *testing.T) {
	type fields struct {
		Logger *zerolog.Logger
		Last   string
		First  string
		yomi   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{nil, "Last", "First", "Yomi"}, "Last First (ワタナベ Yomi)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Watanabe{
				Logger: tt.fields.Logger,
				Last:   tt.fields.Last,
				First:  tt.fields.First,
				Yomi:   tt.fields.yomi,
			}
			assert.Equal(t, tt.want, w.String())
		})
	}
}
