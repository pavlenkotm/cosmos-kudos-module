package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)

func TestMsgSendKudos_ValidateBasic(t *testing.T) {
	tests := []struct {
		name      string
		msg       types.MsgSendKudos
		expectErr bool
		errType   error
	}{
		{
			name: "valid message",
			msg: types.MsgSendKudos{
				FromAddress: "cosmos1from",
				ToAddress:   "cosmos1to",
				Amount:      100,
				Comment:     "Great work!",
			},
			expectErr: false,
		},
		{
			name: "invalid from address",
			msg: types.MsgSendKudos{
				FromAddress: "invalid",
				ToAddress:   "cosmos1to",
				Amount:      100,
				Comment:     "Great work!",
			},
			expectErr: true,
			errType:   types.ErrInvalidAddress,
		},
		{
			name: "invalid to address",
			msg: types.MsgSendKudos{
				FromAddress: "cosmos1from",
				ToAddress:   "invalid",
				Amount:      100,
				Comment:     "Great work!",
			},
			expectErr: true,
			errType:   types.ErrInvalidAddress,
		},
		{
			name: "same address",
			msg: types.MsgSendKudos{
				FromAddress: "cosmos1same",
				ToAddress:   "cosmos1same",
				Amount:      100,
				Comment:     "Self kudos",
			},
			expectErr: true,
			errType:   types.ErrSameAddress,
		},
		{
			name: "zero amount",
			msg: types.MsgSendKudos{
				FromAddress: "cosmos1from",
				ToAddress:   "cosmos1to",
				Amount:      0,
				Comment:     "Zero kudos",
			},
			expectErr: true,
			errType:   types.ErrInvalidAmount,
		},
		{
			name: "comment too long",
			msg: types.MsgSendKudos{
				FromAddress: "cosmos1from",
				ToAddress:   "cosmos1to",
				Amount:      100,
				Comment:     "This is a very long comment that exceeds the maximum allowed length of 140 characters. It should fail validation because it's way too long for a kudos comment.",
			},
			expectErr: true,
			errType:   types.ErrCommentTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.expectErr {
				require.Error(t, err)
				if tt.errType != nil {
					require.ErrorIs(t, err, tt.errType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
