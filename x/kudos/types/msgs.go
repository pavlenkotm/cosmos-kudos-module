package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendKudos{}

// ValidateBasic performs stateless validation on MsgSendKudos
func (msg *MsgSendKudos) ValidateBasic() error {
	// Validate from address
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid from address: %s", err)
	}

	// Validate to address
	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid to address: %s", err)
	}

	// Cannot send to yourself
	if msg.FromAddress == msg.ToAddress {
		return ErrSameAddress
	}

	// Amount must be greater than 0
	if msg.Amount == 0 {
		return ErrInvalidAmount
	}

	// Comment must not exceed 140 characters
	if len(msg.Comment) > 140 {
		return ErrCommentTooLong
	}

	return nil
}

// GetSigners returns the expected signers for MsgSendKudos
func (msg *MsgSendKudos) GetSigners() []sdk.AccAddress {
	fromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{fromAddress}
}
