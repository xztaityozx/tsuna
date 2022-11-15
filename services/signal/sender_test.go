package signal

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		logger *zerolog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Sender
	}{
		{"Newできるべき", args{logger: nil}, &Sender{logger: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewSender(tt.args.logger), tt.want)
		})
	}
}
