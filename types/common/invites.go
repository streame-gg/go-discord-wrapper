package common

type InviteTargetUserType int

const (
	InviteTargetUserTypeStream              InviteTargetUserType = 1
	InviteTargetUserTypeEmbeddedApplication InviteTargetUserType = 2
)

type InviteFlags int

const (
	IsGuestInvite InviteFlags = 1 << 0
)
