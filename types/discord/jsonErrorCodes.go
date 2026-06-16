package discord

// JSONErrorCode is a Discord JSON error code returned in REST API responses
// (the numeric "code" field of an error body). These are distinct from gateway
// close codes (4000–4014). The full list lives at
// https://docs.discord.com/developers/topics/opcodes-and-status-codes#json.
type JSONErrorCode int

const (
	// General error
	JSONErrorCodeGeneral JSONErrorCode = 0

	// 1xxxx — unknown resource
	JSONErrorCodeUnknownAccount                        JSONErrorCode = 10001
	JSONErrorCodeUnknownApplication                    JSONErrorCode = 10002
	JSONErrorCodeUnknownChannel                        JSONErrorCode = 10003
	JSONErrorCodeUnknownGuild                          JSONErrorCode = 10004
	JSONErrorCodeUnknownIntegration                    JSONErrorCode = 10005
	JSONErrorCodeUnknownInvite                         JSONErrorCode = 10006
	JSONErrorCodeUnknownMember                         JSONErrorCode = 10007
	JSONErrorCodeUnknownMessage                        JSONErrorCode = 10008
	JSONErrorCodeUnknownPermissionOverwrite            JSONErrorCode = 10009
	JSONErrorCodeUnknownProvider                       JSONErrorCode = 10010
	JSONErrorCodeUnknownRole                           JSONErrorCode = 10011
	JSONErrorCodeUnknownToken                          JSONErrorCode = 10012
	JSONErrorCodeUnknownUser                           JSONErrorCode = 10013
	JSONErrorCodeUnknownEmoji                          JSONErrorCode = 10014
	JSONErrorCodeUnknownWebhook                        JSONErrorCode = 10015
	JSONErrorCodeUnknownWebhookService                 JSONErrorCode = 10016
	JSONErrorCodeUnknownSession                        JSONErrorCode = 10020
	JSONErrorCodeUnknownAsset                          JSONErrorCode = 10021
	JSONErrorCodeUnknownBan                            JSONErrorCode = 10026
	JSONErrorCodeUnknownSKU                            JSONErrorCode = 10027
	JSONErrorCodeUnknownStoreListing                   JSONErrorCode = 10028
	JSONErrorCodeUnknownEntitlement                    JSONErrorCode = 10029
	JSONErrorCodeUnknownBuild                          JSONErrorCode = 10030
	JSONErrorCodeUnknownLobby                          JSONErrorCode = 10031
	JSONErrorCodeUnknownBranch                         JSONErrorCode = 10032
	JSONErrorCodeUnknownStoreDirectoryLayout           JSONErrorCode = 10033
	JSONErrorCodeUnknownRedistributable                JSONErrorCode = 10036
	JSONErrorCodeUnknownGiftCode                       JSONErrorCode = 10038
	JSONErrorCodeUnknownStream                         JSONErrorCode = 10049
	JSONErrorCodeUnknownPremiumServerSubscribeCooldown JSONErrorCode = 10050
	JSONErrorCodeUnknownGuildTemplate                  JSONErrorCode = 10057
	JSONErrorCodeUnknownDiscoverableServerCategory     JSONErrorCode = 10059
	JSONErrorCodeUnknownSticker                        JSONErrorCode = 10060
	JSONErrorCodeUnknownStickerPack                    JSONErrorCode = 10061
	JSONErrorCodeUnknownInteraction                    JSONErrorCode = 10062
	JSONErrorCodeUnknownApplicationCommand             JSONErrorCode = 10063
	JSONErrorCodeUnknownVoiceState                     JSONErrorCode = 10065
	JSONErrorCodeUnknownApplicationCommandPermissions  JSONErrorCode = 10066
	JSONErrorCodeUnknownStageInstance                  JSONErrorCode = 10067
	JSONErrorCodeUnknownGuildMemberVerificationForm    JSONErrorCode = 10068
	JSONErrorCodeUnknownGuildWelcomeScreen             JSONErrorCode = 10069
	JSONErrorCodeUnknownGuildScheduledEvent            JSONErrorCode = 10070
	JSONErrorCodeUnknownGuildScheduledEventUser        JSONErrorCode = 10071
	JSONErrorCodeUnknownTag                            JSONErrorCode = 10087
	JSONErrorCodeUnknownSound                          JSONErrorCode = 10097
	JSONErrorCodeUnknownInviteTargetUsersJob           JSONErrorCode = 10124
	JSONErrorCodeUnknownInviteTargetUsers              JSONErrorCode = 10129

	// 2xxxx — endpoint / action restrictions
	JSONErrorCodeBotsCannotUseEndpoint             JSONErrorCode = 20001
	JSONErrorCodeOnlyBotsCanUseEndpoint            JSONErrorCode = 20002
	JSONErrorCodeExplicitContentCannotBeSentToUser JSONErrorCode = 20009
	JSONErrorCodeNotAuthorizedForApplication       JSONErrorCode = 20012
	JSONErrorCodeSlowmodeRateLimit                 JSONErrorCode = 20016
	JSONErrorCodeOnlyOwnerCanPerformAction         JSONErrorCode = 20018
	JSONErrorCodeAnnouncementEditRateLimit         JSONErrorCode = 20022
	JSONErrorCodeUnderMinimumAge                   JSONErrorCode = 20024
	JSONErrorCodeChannelWriteRateLimit             JSONErrorCode = 20028
	JSONErrorCodeServerWriteRateLimit              JSONErrorCode = 20029
	JSONErrorCodeContainsDisallowedWords           JSONErrorCode = 20031
	JSONErrorCodePremiumLevelTooLow                JSONErrorCode = 20035

	// 3xxxx — resource limit reached
	JSONErrorCodeMaxGuilds                      JSONErrorCode = 30001
	JSONErrorCodeMaxFriends                     JSONErrorCode = 30002
	JSONErrorCodeMaxPinsInChannel               JSONErrorCode = 30003
	JSONErrorCodeMaxRecipientsInDM              JSONErrorCode = 30004
	JSONErrorCodeMaxGuildRoles                  JSONErrorCode = 30005
	JSONErrorCodeMaxWebhooks                    JSONErrorCode = 30007
	JSONErrorCodeMaxEmojis                      JSONErrorCode = 30008
	JSONErrorCodeMaxReactions                   JSONErrorCode = 30010
	JSONErrorCodeMaxGroupDMs                    JSONErrorCode = 30011
	JSONErrorCodeMaxGuildChannels               JSONErrorCode = 30013
	JSONErrorCodeMaxAttachments                 JSONErrorCode = 30015
	JSONErrorCodeMaxInvites                     JSONErrorCode = 30016
	JSONErrorCodeMaxAnimatedEmojis              JSONErrorCode = 30018
	JSONErrorCodeMaxServerMembers               JSONErrorCode = 30019
	JSONErrorCodeMaxGuildCategories             JSONErrorCode = 30030
	JSONErrorCodeGuildAlreadyHasTemplate        JSONErrorCode = 30031
	JSONErrorCodeMaxApplicationCommands         JSONErrorCode = 30032
	JSONErrorCodeMaxThreadParticipants          JSONErrorCode = 30033
	JSONErrorCodeMaxDailyAppCommandCreates      JSONErrorCode = 30034
	JSONErrorCodeMaxBansForNonMembers           JSONErrorCode = 30035
	JSONErrorCodeMaxBansFetched                 JSONErrorCode = 30037
	JSONErrorCodeMaxScheduledEvents             JSONErrorCode = 30038
	JSONErrorCodeMaxStickers                    JSONErrorCode = 30039
	JSONErrorCodeMaxPruneRequests               JSONErrorCode = 30040
	JSONErrorCodeMaxWidgetSettingsUpdates       JSONErrorCode = 30042
	JSONErrorCodeMaxSoundboardSounds            JSONErrorCode = 30045
	JSONErrorCodeMaxOldMessageEdits             JSONErrorCode = 30046
	JSONErrorCodeMaxPinnedThreadsInForum        JSONErrorCode = 30047
	JSONErrorCodeMaxForumTags                   JSONErrorCode = 30048
	JSONErrorCodeBitrateTooHighForChannelType   JSONErrorCode = 30052
	JSONErrorCodeMaxPremiumEmojis               JSONErrorCode = 30056
	JSONErrorCodeMaxWebhooksPerGuild            JSONErrorCode = 30058
	JSONErrorCodeMaxChannelPermissionOverwrites JSONErrorCode = 30060
	JSONErrorCodeGuildChannelsTooLarge          JSONErrorCode = 30061

	// 4xxxx — action errors
	JSONErrorCodeUnauthorized                     JSONErrorCode = 40001
	JSONErrorCodeVerifyAccountRequired            JSONErrorCode = 40002
	JSONErrorCodeOpeningDMsTooFast                JSONErrorCode = 40003
	JSONErrorCodeSendMessagesTemporarilyDisabled  JSONErrorCode = 40004
	JSONErrorCodeRequestEntityTooLarge            JSONErrorCode = 40005
	JSONErrorCodeFeatureTemporarilyDisabled       JSONErrorCode = 40006
	JSONErrorCodeUserBannedFromGuild              JSONErrorCode = 40007
	JSONErrorCodeConnectionRevoked                JSONErrorCode = 40012
	JSONErrorCodeOnlyConsumableSKUsCanBeConsumed  JSONErrorCode = 40018
	JSONErrorCodeOnlySandboxEntitlementsDeletable JSONErrorCode = 40019
	JSONErrorCodeTargetUserNotInVoice             JSONErrorCode = 40032
	JSONErrorCodeMessageAlreadyCrossposted        JSONErrorCode = 40033
	JSONErrorCodeAppCommandNameAlreadyExists      JSONErrorCode = 40041
	JSONErrorCodeAppInteractionFailedToSend       JSONErrorCode = 40043
	JSONErrorCodeCannotSendMessageInForumChannel  JSONErrorCode = 40058
	JSONErrorCodeInteractionAlreadyAcknowledged   JSONErrorCode = 40060
	JSONErrorCodeTagNamesMustBeUnique             JSONErrorCode = 40061
	JSONErrorCodeServiceResourceRateLimited       JSONErrorCode = 40062
	JSONErrorCodeNoTagsAvailableForNonModerators  JSONErrorCode = 40066
	JSONErrorCodeTagRequiredForForumPost          JSONErrorCode = 40067
	JSONErrorCodeEntitlementAlreadyGranted        JSONErrorCode = 40074
	JSONErrorCodeMaxInteractionFollowups          JSONErrorCode = 40094
	JSONErrorCodeCloudflareBlocking               JSONErrorCode = 40333

	// 5xxxx — permission / invalid access errors
	JSONErrorCodeMissingAccess                             JSONErrorCode = 50001
	JSONErrorCodeInvalidAccountType                        JSONErrorCode = 50002
	JSONErrorCodeCannotExecuteOnDMChannel                  JSONErrorCode = 50003
	JSONErrorCodeGuildWidgetDisabled                       JSONErrorCode = 50004
	JSONErrorCodeCannotEditOthersMessages                  JSONErrorCode = 50005
	JSONErrorCodeCannotSendEmptyMessage                    JSONErrorCode = 50006
	JSONErrorCodeCannotSendMessagesToThisUser              JSONErrorCode = 50007
	JSONErrorCodeCannotSendMessagesInNonTextChannel        JSONErrorCode = 50008
	JSONErrorCodeChannelVerificationLevelTooHigh           JSONErrorCode = 50009
	JSONErrorCodeOauth2ApplicationHasNoBot                 JSONErrorCode = 50010
	JSONErrorCodeOauth2ApplicationLimitReached             JSONErrorCode = 50011
	JSONErrorCodeInvalidOauth2State                        JSONErrorCode = 50012
	JSONErrorCodeMissingPermissions                        JSONErrorCode = 50013
	JSONErrorCodeInvalidAuthenticationToken                JSONErrorCode = 50014
	JSONErrorCodeNoteTooLong                               JSONErrorCode = 50015
	JSONErrorCodeInvalidBulkDeleteCount                    JSONErrorCode = 50016
	JSONErrorCodeInvalidMFALevel                           JSONErrorCode = 50017
	JSONErrorCodeCannotPinMessageInOtherChannel            JSONErrorCode = 50019
	JSONErrorCodeInviteCodeInvalidOrTaken                  JSONErrorCode = 50020
	JSONErrorCodeCannotExecuteOnSystemMessage              JSONErrorCode = 50021
	JSONErrorCodeCannotExecuteOnThisChannelType            JSONErrorCode = 50024
	JSONErrorCodeInvalidOauth2AccessToken                  JSONErrorCode = 50025
	JSONErrorCodeMissingRequiredOauth2Scope                JSONErrorCode = 50026
	JSONErrorCodeInvalidWebhookToken                       JSONErrorCode = 50027
	JSONErrorCodeInvalidRole                               JSONErrorCode = 50028
	JSONErrorCodeInvalidRecipients                         JSONErrorCode = 50033
	JSONErrorCodeMessageTooOldForBulkDelete                JSONErrorCode = 50034
	JSONErrorCodeInvalidFormBody                           JSONErrorCode = 50035
	JSONErrorCodeInviteAcceptedForGuildBotNotIn            JSONErrorCode = 50036
	JSONErrorCodeInvalidActivityAction                     JSONErrorCode = 50039
	JSONErrorCodeInvalidAPIVersion                         JSONErrorCode = 50041
	JSONErrorCodeFileUploadExceedsMaxSize                  JSONErrorCode = 50045
	JSONErrorCodeInvalidFileUploaded                       JSONErrorCode = 50046
	JSONErrorCodeCannotSelfRedeemGift                      JSONErrorCode = 50054
	JSONErrorCodeInvalidGuild                              JSONErrorCode = 50055
	JSONErrorCodeInvalidSKU                                JSONErrorCode = 50057
	JSONErrorCodeInvalidRequestOrigin                      JSONErrorCode = 50067
	JSONErrorCodeInvalidMessageType                        JSONErrorCode = 50068
	JSONErrorCodePaymentSourceRequiredToRedeemGift         JSONErrorCode = 50070
	JSONErrorCodeCannotModifySystemWebhook                 JSONErrorCode = 50073
	JSONErrorCodeCannotDeleteCommunityRequiredChannel      JSONErrorCode = 50074
	JSONErrorCodeCannotEditStickersWithinMessage           JSONErrorCode = 50080
	JSONErrorCodeInvalidStickerSent                        JSONErrorCode = 50081
	JSONErrorCodeOperationOnArchivedThread                 JSONErrorCode = 50083
	JSONErrorCodeInvalidThreadNotificationSettings         JSONErrorCode = 50084
	JSONErrorCodeBeforeValueBeforeThreadCreation           JSONErrorCode = 50085
	JSONErrorCodeCommunityChannelsMustBeTextChannels       JSONErrorCode = 50086
	JSONErrorCodeEventEntityTypeMismatch                   JSONErrorCode = 50091
	JSONErrorCodeServerNotAvailableInYourLocation          JSONErrorCode = 50095
	JSONErrorCodeServerNeedsMonetizationEnabled            JSONErrorCode = 50097
	JSONErrorCodeServerNeedsMoreBoosts                     JSONErrorCode = 50101
	JSONErrorCodeRequestBodyContainsInvalidJSON            JSONErrorCode = 50109
	JSONErrorCodeProvidedFileInvalid                       JSONErrorCode = 50110
	JSONErrorCodeProvidedFileTypeInvalid                   JSONErrorCode = 50123
	JSONErrorCodeFileDurationExceedsMax                    JSONErrorCode = 50124
	JSONErrorCodeOwnerCannotBePendingMember                JSONErrorCode = 50131
	JSONErrorCodeCannotTransferOwnershipToBot              JSONErrorCode = 50132
	JSONErrorCodeFailedToResizeAsset                       JSONErrorCode = 50138
	JSONErrorCodeCannotMixSubscriptionEmojiRoles           JSONErrorCode = 50144
	JSONErrorCodeCannotConvertBetweenPremiumAndNormalEmoji JSONErrorCode = 50145
	JSONErrorCodeUploadedFileNotFound                      JSONErrorCode = 50146
	JSONErrorCodeInvalidEmoji                              JSONErrorCode = 50151
	JSONErrorCodeVoiceMessageNoAdditionalContent           JSONErrorCode = 50159
	JSONErrorCodeVoiceMessageRequiresSingleAudioAttachment JSONErrorCode = 50160
	JSONErrorCodeVoiceMessageRequiresMetadata              JSONErrorCode = 50161
	JSONErrorCodeVoiceMessageCannotBeEdited                JSONErrorCode = 50162
	JSONErrorCodeCannotDeleteGuildSubscriptionIntegration  JSONErrorCode = 50163
	JSONErrorCodeCannotSendVoiceMessagesInChannel          JSONErrorCode = 50173
	JSONErrorCodeUserAccountMustBeVerified                 JSONErrorCode = 50178
	JSONErrorCodeFileLacksValidDuration                    JSONErrorCode = 50192
	JSONErrorCodeCannotMessageUserNoMutualGuilds           JSONErrorCode = 50278
	JSONErrorCodeMissingPermissionToSendSticker            JSONErrorCode = 50600

	// 6xxxx — authentication requirements
	JSONErrorCodeTwoFactorRequired JSONErrorCode = 60003

	// 8xxxx — user lookup
	JSONErrorCodeNoUsersWithDiscordTagExist JSONErrorCode = 80004

	// 9xxxx — reactions
	JSONErrorCodeReactionBlocked             JSONErrorCode = 90001
	JSONErrorCodeUserCannotUseBurstReactions JSONErrorCode = 90002

	// 11xxxx — availability
	JSONErrorCodeIndexNotYetAvailable       JSONErrorCode = 110000
	JSONErrorCodeApplicationNotYetAvailable JSONErrorCode = 110001

	// 13xxxx — overload
	JSONErrorCodeAPIResourceOverloaded JSONErrorCode = 130000

	// 15xxxx — stage instances
	JSONErrorCodeStageAlreadyOpen JSONErrorCode = 150006

	// 16xxxx — thread errors
	JSONErrorCodeCannotReplyWithoutReadMessageHistoryPermission JSONErrorCode = 160002
	JSONErrorCodeThreadAlreadyCreatedForMessage                 JSONErrorCode = 160004
	JSONErrorCodeThreadIsLocked                                 JSONErrorCode = 160005
	JSONErrorCodeMaxActiveThreads                               JSONErrorCode = 160006
	JSONErrorCodeMaxActiveAnnouncementThreads                   JSONErrorCode = 160007
	JSONErrorCodeCannotForwardUnreadableMessage                 JSONErrorCode = 160014

	// 17xxxx — Lottie / sticker validation
	JSONErrorCodeInvalidLottieJSON                  JSONErrorCode = 170001
	JSONErrorCodeLottieCannotContainRasterImages    JSONErrorCode = 170002
	JSONErrorCodeStickerMaxFramerateExceeded        JSONErrorCode = 170003
	JSONErrorCodeStickerFrameCountExceedsMax        JSONErrorCode = 170004
	JSONErrorCodeLottieAnimationDimensionsExceeded  JSONErrorCode = 170005
	JSONErrorCodeStickerFrameRateInvalid            JSONErrorCode = 170006
	JSONErrorCodeStickerAnimationDurationExceedsMax JSONErrorCode = 170007

	// 18xxxx — scheduled events / stages
	JSONErrorCodeCannotUpdateFinishedEvent   JSONErrorCode = 180000
	JSONErrorCodeFailedToCreateStageForEvent JSONErrorCode = 180002

	// 20xxxx — automatic moderation
	JSONErrorCodeMessageBlockedByAutoMod JSONErrorCode = 200000
	JSONErrorCodeTitleBlockedByAutoMod   JSONErrorCode = 200001

	// 22xxxx — forum webhooks
	JSONErrorCodeWebhookForumThreadNameOrIDRequired   JSONErrorCode = 220001
	JSONErrorCodeWebhookForumThreadNameAndIDExclusive JSONErrorCode = 220002
	JSONErrorCodeWebhookCanOnlyCreateThreadsInForum   JSONErrorCode = 220003
	JSONErrorCodeWebhookServicesNotAllowedInForum     JSONErrorCode = 220004

	// 24xxxx — harmful links
	JSONErrorCodeMessageBlockedByHarmfulLinksFilter JSONErrorCode = 240000

	// 35xxxx — onboarding
	JSONErrorCodeCannotEnableOnboardingRequirementsNotMet     JSONErrorCode = 350000
	JSONErrorCodeCannotUpdateOnboardingWhileBelowRequirements JSONErrorCode = 350001

	// 40xxxx — uploads
	JSONErrorCodeFileUploadsLimitedForGuild JSONErrorCode = 400001

	// 50xxxx — bans
	JSONErrorCodeFailedToBanUsers JSONErrorCode = 500000

	// 52xxxx — polls
	JSONErrorCodePollVotingBlocked              JSONErrorCode = 520000
	JSONErrorCodePollExpired                    JSONErrorCode = 520001
	JSONErrorCodeInvalidChannelTypeForPoll      JSONErrorCode = 520002
	JSONErrorCodeCannotEditPollMessage          JSONErrorCode = 520003
	JSONErrorCodeCannotUseEmojiIncludedWithPoll JSONErrorCode = 520004
	JSONErrorCodeCannotExpireNonPollMessage     JSONErrorCode = 520006

	// 53xxxx — provisional accounts / OpenID
	JSONErrorCodeProvisionalAccountsNotGranted  JSONErrorCode = 530000
	JSONErrorCodeIDTokenExpired                 JSONErrorCode = 530001
	JSONErrorCodeIDTokenIssuerMismatch          JSONErrorCode = 530002
	JSONErrorCodeIDTokenAudienceMismatch        JSONErrorCode = 530003
	JSONErrorCodeIDTokenIssuedTooLongAgo        JSONErrorCode = 530004
	JSONErrorCodeFailedToGenerateUniqueUsername JSONErrorCode = 530006
	JSONErrorCodeInvalidClientSecret            JSONErrorCode = 530007
)
