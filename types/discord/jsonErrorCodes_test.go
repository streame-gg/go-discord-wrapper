package discord

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// jsonErrorCodesSuite locks in the numeric values that were previously wrong
// (an off-by-one in the 10010–10013 range plus several mislabelled blocks) and
// confirms the deprecated GatewayErrorCode alias still resolves.
type jsonErrorCodesSuite struct {
	suite.Suite
}

func TestJSONErrorCodesSuite(t *testing.T) {
	suite.Run(t, new(jsonErrorCodesSuite))
}

func (s *jsonErrorCodesSuite) TestCorrectedValues() {
	cases := map[string]struct {
		got  JSONErrorCode
		want int
	}{
		// Previously off-by-one (10010 was missing; role/token/user shifted down).
		"UnknownProvider": {JSONErrorCodeUnknownProvider, 10010},
		"UnknownRole":     {JSONErrorCodeUnknownRole, 10011},
		"UnknownToken":    {JSONErrorCodeUnknownToken, 10012},
		"UnknownUser":     {JSONErrorCodeUnknownUser, 10013},
		// Previously swapped.
		"UnknownSticker":     {JSONErrorCodeUnknownSticker, 10060},
		"UnknownStickerPack": {JSONErrorCodeUnknownStickerPack, 10061},
		// 4xxxx block was shifted.
		"SendMessagesTemporarilyDisabled": {JSONErrorCodeSendMessagesTemporarilyDisabled, 40004},
		"UserBannedFromGuild":             {JSONErrorCodeUserBannedFromGuild, 40007},
		"ConnectionRevoked":               {JSONErrorCodeConnectionRevoked, 40012},
		"TargetUserNotInVoice":            {JSONErrorCodeTargetUserNotInVoice, 40032},
		"MessageAlreadyCrossposted":       {JSONErrorCodeMessageAlreadyCrossposted, 40033},
		"AppCommandNameAlreadyExists":     {JSONErrorCodeAppCommandNameAlreadyExists, 40041},
		// Previously mislabelled / misplaced blocks.
		"CannotTransferOwnershipToBot": {JSONErrorCodeCannotTransferOwnershipToBot, 50132},
		"UploadedFileNotFound":         {JSONErrorCodeUploadedFileNotFound, 50146},
		"UserAccountMustBeVerified":    {JSONErrorCodeUserAccountMustBeVerified, 50178},
		"ReactionBlocked":              {JSONErrorCodeReactionBlocked, 90001},
		"APIResourceOverloaded":        {JSONErrorCodeAPIResourceOverloaded, 130000},
		"StageAlreadyOpen":             {JSONErrorCodeStageAlreadyOpen, 150006},
		"ThreadIsLocked":               {JSONErrorCodeThreadIsLocked, 160005},
		"InvalidLottieJSON":            {JSONErrorCodeInvalidLottieJSON, 170001},
		"MessageBlockedByAutoMod":      {JSONErrorCodeMessageBlockedByAutoMod, 200000},
		"FailedToBanUsers":             {JSONErrorCodeFailedToBanUsers, 500000},
		// Newly added current-feature codes.
		"PollExpired":                   {JSONErrorCodePollExpired, 520001},
		"ProvisionalAccountsNotGranted": {JSONErrorCodeProvisionalAccountsNotGranted, 530000},
	}
	for name, c := range cases {
		s.Equalf(c.want, int(c.got), "%s", name)
	}
}

func (s *jsonErrorCodesSuite) TestGatewayErrorCodeAlias() {
	// The deprecated alias must remain assignable to/from JSONErrorCode.
	var g GatewayErrorCode = JSONErrorCodeUnknownChannel
	s.Equal(JSONErrorCodeUnknownChannel, g)
	s.Equal(10003, int(g))
}
