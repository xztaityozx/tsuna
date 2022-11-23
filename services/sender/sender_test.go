package sender

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		failOnError   bool
		watchInterval time.Duration
		watchAttempt  int
	}
	tests := []struct {
		name string
		args args
		want *Sender
	}{
		{"Newできるべき", args{failOnError: false, watchInterval: 1 * time.Second, watchAttempt: 1}, &Sender{failOnError: false, watchInterval: 1 * time.Second, watchAttempt: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSender(tt.args.failOnError, tt.args.watchInterval, tt.args.watchAttempt)
			assert.Equal(t, tt.want.watchAttempt, got.watchAttempt)
			assert.Equal(t, tt.want.watchInterval, got.watchInterval)
			assert.NotNil(t, got.Logger)
			assert.NotNil(t, got.Watanabe)
		})
	}
}
