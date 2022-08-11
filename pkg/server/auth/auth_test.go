//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth_test

import (
	"github.com/a-clap/dictionary/pkg/server/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestToken_Validate(t *testing.T) {
	type args struct {
		name     string
		duration time.Duration
	}
	type token struct {
		err bool
	}
	type validate struct {
		err       bool
		validated bool
	}
	tests := []struct {
		name     string
		args     args
		token    token
		validate validate
	}{
		{
			name: "validation test",
			args: args{
				name:     "adam",
				duration: 3 * time.Second,
			},
			token: token{
				err: false,
			},
			validate: validate{
				err:       false,
				validated: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.Token(tt.args.name, tt.args.duration)

			if tt.token.err {
				assert.NotNil(t, err, tt.name)
			} else {
				assert.Nil(t, err, tt.name)
			}

			validated, err := auth.Validate(got)

			if tt.validate.err {
				assert.NotNil(t, err, tt.name)
			} else {
				assert.Nil(t, err, tt.name)
			}

			require.Equal(t, tt.validate.validated, validated)
		})
	}
}
