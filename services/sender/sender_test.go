package sender

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		failOnError bool
	}
	tests := []struct {
		name string
		args args
		want *Sender
	}{
		{"Newできるべき", args{failOnError: false}, &Sender{failOnError: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewSender(tt.args.failOnError), tt.want)
		})
	}
}
