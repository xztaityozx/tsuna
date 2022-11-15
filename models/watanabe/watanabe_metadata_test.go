package watanabe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Metadata
	}{
		{"Newできるべき", Metadata{First: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := New()
			assert.NotNil(t, actual)
			assert.NotEmpty(t, actual.First)
		})
	}
}
