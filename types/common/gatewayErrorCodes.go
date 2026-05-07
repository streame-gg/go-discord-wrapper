package common

type GatewayErrorCode int

const (
	// General errors
	GatewayErrorCodeGeneral GatewayErrorCode = 0

	//1xxxx
	GatewayErrorCodeUnknownAccount                       GatewayErrorCode = 10001
	GatewayErrorCodeUnknownApplication                   GatewayErrorCode = 10002
	GatewayErrorCodeUnknownChannel                       GatewayErrorCode = 10003
	GatewayErrorCodeUnknownGuild                         GatewayErrorCode = 10004
	GatewayErrorCodeUnknownIntegration                   GatewayErrorCode = 10005
	GatewayErrorCodeUnknownInvite                        GatewayErrorCode = 10006
	GatewayErrorCodeUnknownMember                        GatewayErrorCode = 10007
	GatewayErrorCodeUnknownMessage                       GatewayErrorCode = 10008
	GatewayErrorCodeUnknownPermissionOverwrite           GatewayErrorCode = 10009
	GatewayErrorCodeUnknownRole                          GatewayErrorCode = 10010
	GatewayErrorCodeUnknownToken                         GatewayErrorCode = 10011
	GatewayErrorCodeUnknownUser                          GatewayErrorCode = 10012
	GatewayErrorCodeUnknownEmoji                         GatewayErrorCode = 10014
	GatewayErrorCodeUnknownWebhook                       GatewayErrorCode = 10015
	GatewayErrorCodeUnknownWebhookService                GatewayErrorCode = 10016
	GatewayErrorCodeUnknownSession                       GatewayErrorCode = 10020
	GatewayErrorCodeUnknownAsset                         GatewayErrorCode = 10021
	GatewayErrorCodeUnknownBan                           GatewayErrorCode = 10026
	GatewayErrorCodeUnknownSku                           GatewayErrorCode = 10027
	GatewayErrorCodeUnknownStoreListing                  GatewayErrorCode = 10028
	GatewayErrorCodeUnknownEntitlement                   GatewayErrorCode = 10029
	GatewayErrorCodeUnknownBuild                         GatewayErrorCode = 10030
	GatewayErrorCodeUnknownLobby                         GatewayErrorCode = 10031
	GatewayErrorCodeUnknownBranch                        GatewayErrorCode = 10032
	GatewayErrorCodeUnknownRedistributable               GatewayErrorCode = 10036
	GatewayErrorCodeUnknownGiftCode                      GatewayErrorCode = 10038
	GatewayErrorCodeUnknownStream                        GatewayErrorCode = 10049
	GatewayErrorCodeUnknownPremiumSubscriptionOption     GatewayErrorCode = 10050
	GatewayErrorCodeUnknownGuildTemplate                 GatewayErrorCode = 10057
	GatewayErrorCodeUnknownDiscoverableServer            GatewayErrorCode = 10059
	GatewayErrorCodeUnknownStickerPack                   GatewayErrorCode = 10060
	GatewayErrorCodeUnknownSticker                       GatewayErrorCode = 10061
	GatewayErrorCodeUnknownInteraction                   GatewayErrorCode = 10062
	GatewayErrorCodeUnknownApplicationCommand            GatewayErrorCode = 10063
	GatewayErrorCodeUnknownVoiceState                    GatewayErrorCode = 10065
	GatewayErrorCodeUnknownApplicationCommandPermissions GatewayErrorCode = 10066
	GatewayErrorCodeUnknownStageInstance                 GatewayErrorCode = 10067
	GatewayErrorCodeUnknownGuildMemberVerificationForm   GatewayErrorCode = 10068
	GatewayErrorCodeUnknownGuildWelcomeScreen            GatewayErrorCode = 10069
	GatewayErrorCodeUnknownScheduledEvent                GatewayErrorCode = 10070
	GatewayErrorCodeUnknownScheduledEventUser            GatewayErrorCode = 10071
	GatewayErrorCodeUnknownTag                           GatewayErrorCode = 10087
	GatewayErrorCodeUnknownSound                         GatewayErrorCode = 10097
	GatewayErrorCodeUnknownInviteTargetUserJob           GatewayErrorCode = 10124
	GatewayErrorCodeUnknownInviteTargetUsers             GatewayErrorCode = 10129

	//2xxxx
	GatewayErrorCodeBotsCannotUseEndpoint                     GatewayErrorCode = 20001
	GatewayErrorCodeOnlyBotsCanUseEndpoint                    GatewayErrorCode = 20002
	GatewayErrorCodeExplicitContentCannotBeSentToTheRecipient GatewayErrorCode = 20009
	GatewayErrorCodeNotAuthorizedToMessageUser                GatewayErrorCode = 20012
	GatewayErrorCodeActionCannotBePerformedSlowmodeRateLimit  GatewayErrorCode = 20016
	GatewayErrorCodeOnlyOwnerCanPerformAction                 GatewayErrorCode = 20018
	GatewayErrorCodeMessageEditFailedAnnouncementRateLimit    GatewayErrorCode = 20022
	GatewayErrorCodeUnderMinimumAge                           GatewayErrorCode = 20024
	GatewayErrorCodeChannelHasWriteRateLimit                  GatewayErrorCode = 20028
	GatewayErrorCodeServerHasWriteRateLimit                   GatewayErrorCode = 20029
	GatewayErrorCodeContainsDisallowedWords                   GatewayErrorCode = 20031
	GatewayErrorCodePremiumLevelTooLow                        GatewayErrorCode = 20035

	//3xxxx — resource limit reached
	GatewayErrorCodeMaxGuilds                      GatewayErrorCode = 30001
	GatewayErrorCodeMaxFriends                     GatewayErrorCode = 30002
	GatewayErrorCodeMaxPinsInChannel               GatewayErrorCode = 30003
	GatewayErrorCodeMaxRecipientsInDM              GatewayErrorCode = 30004
	GatewayErrorCodeMaxGuildRoles                  GatewayErrorCode = 30005
	GatewayErrorCodeMaxWebhooks                    GatewayErrorCode = 30007
	GatewayErrorCodeMaxEmojis                      GatewayErrorCode = 30008
	GatewayErrorCodeMaxReactions                   GatewayErrorCode = 30010
	GatewayErrorCodeMaxGroupDMs                    GatewayErrorCode = 30011
	GatewayErrorCodeMaxGuildChannels               GatewayErrorCode = 30013
	GatewayErrorCodeMaxAttachments                 GatewayErrorCode = 30015
	GatewayErrorCodeMaxInvites                     GatewayErrorCode = 30016
	GatewayErrorCodeMaxAnimatedEmojis              GatewayErrorCode = 30018
	GatewayErrorCodeMaxServerMembers               GatewayErrorCode = 30019
	GatewayErrorCodeMaxGuildCategories             GatewayErrorCode = 30030
	GatewayErrorCodeGuildAlreadyHasTemplate        GatewayErrorCode = 30031
	GatewayErrorCodeMaxApplicationCommands         GatewayErrorCode = 30032
	GatewayErrorCodeMaxThreadParticipants          GatewayErrorCode = 30033
	GatewayErrorCodeMaxDailyAppCommandCreates      GatewayErrorCode = 30034
	GatewayErrorCodeMaxBansForNonMembers           GatewayErrorCode = 30035
	GatewayErrorCodeMaxBansFetched                 GatewayErrorCode = 30037
	GatewayErrorCodeMaxScheduledEvents             GatewayErrorCode = 30038
	GatewayErrorCodeMaxStickers                    GatewayErrorCode = 30039
	GatewayErrorCodeMaxPruneRequests               GatewayErrorCode = 30040
	GatewayErrorCodeMaxWidgetSettingsUpdates       GatewayErrorCode = 30042
	GatewayErrorCodeMaxOldMessageEdits             GatewayErrorCode = 30046
	GatewayErrorCodeMaxPinnedThreadsInForum        GatewayErrorCode = 30047
	GatewayErrorCodeMaxForumTags                   GatewayErrorCode = 30048
	GatewayErrorCodeBitrateTooHighForChannelType   GatewayErrorCode = 30052
	GatewayErrorCodeMaxPremiumEmojis               GatewayErrorCode = 30056
	GatewayErrorCodeMaxWebhooksPerGuild            GatewayErrorCode = 30058
	GatewayErrorCodeMaxChannelPermissionOverwrites GatewayErrorCode = 30060
	GatewayErrorCodeGuildChannelsTooLarge          GatewayErrorCode = 30061

	//4xxxx — action errors
	GatewayErrorCodeUnauthorized                      GatewayErrorCode = 40001
	GatewayErrorCodeVerifyAccountRequired             GatewayErrorCode = 40002
	GatewayErrorCodeOpeningDMsTooFast                 GatewayErrorCode = 40003
	GatewayErrorCodeRequestEntityTooLarge             GatewayErrorCode = 40004
	GatewayErrorCodeFeatureTemporarilyDisabled        GatewayErrorCode = 40005
	GatewayErrorCodeUserBannedFromGuild               GatewayErrorCode = 40006
	GatewayErrorCodeConnectionRevoked                 GatewayErrorCode = 40007
	GatewayErrorCodeTargetUserNotInVoice              GatewayErrorCode = 40012
	GatewayErrorCodeMessageAlreadyCrossposted         GatewayErrorCode = 40032
	GatewayErrorCodeAppCommandNameAlreadyExists       GatewayErrorCode = 40033
	GatewayErrorCodeAppInteractionAlreadyAcknowledged GatewayErrorCode = 40034
	GatewayErrorCodeMissingRequiredFields             GatewayErrorCode = 40035
	GatewayErrorCodeOauth2TokenError                  GatewayErrorCode = 40036
	GatewayErrorCodeInteractionAlreadyAcknowledged    GatewayErrorCode = 40060
	GatewayErrorCodeTagNamesMustBeUnique              GatewayErrorCode = 40061
	GatewayErrorCodeServiceResourceRateLimited        GatewayErrorCode = 40062
	GatewayErrorCodeNoTagsAvailableForNonModerators   GatewayErrorCode = 40066
	GatewayErrorCodeTagRequiredForForumPost           GatewayErrorCode = 40067

	//5xxxx — permission/invalid access errors
	GatewayErrorCodeMissingAccess                                        GatewayErrorCode = 50001
	GatewayErrorCodeInvalidAccountType                                   GatewayErrorCode = 50002
	GatewayErrorCodeCannotExecuteOnDMChannel                             GatewayErrorCode = 50003
	GatewayErrorCodeGuildWidgetDisabled                                  GatewayErrorCode = 50004
	GatewayErrorCodeCannotEditOthersMessages                             GatewayErrorCode = 50005
	GatewayErrorCodeCannotSendEmptyMessage                               GatewayErrorCode = 50006
	GatewayErrorCodeCannotSendMessagesToThisUser                         GatewayErrorCode = 50007
	GatewayErrorCodeCannotSendMessagesInNonTextChannel                   GatewayErrorCode = 50008
	GatewayErrorCodeChannelVerificationLevelTooHigh                      GatewayErrorCode = 50009
	GatewayErrorCodeOauth2ApplicationHasNoBot                            GatewayErrorCode = 50010
	GatewayErrorCodeOauth2ApplicationLimitReached                        GatewayErrorCode = 50011
	GatewayErrorCodeInvalidOauth2State                                   GatewayErrorCode = 50012
	GatewayErrorCodeMissingPermissions                                   GatewayErrorCode = 50013
	GatewayErrorCodeInvalidAuthenticationToken                           GatewayErrorCode = 50014
	GatewayErrorCodeNoteTooLong                                          GatewayErrorCode = 50015
	GatewayErrorCodeInvalidBulkDeleteCount                               GatewayErrorCode = 50016
	GatewayErrorCodeInvalidMFALevel                                      GatewayErrorCode = 50017
	GatewayErrorCodeCannotPinMessageInOtherChannel                       GatewayErrorCode = 50019
	GatewayErrorCodeInviteCodeInvalidOrTaken                             GatewayErrorCode = 50020
	GatewayErrorCodeCannotExecuteOnSystemMessage                         GatewayErrorCode = 50021
	GatewayErrorCodeCannotExecuteOnThisChannelType                       GatewayErrorCode = 50024
	GatewayErrorCodeInvalidOauth2AccessToken                             GatewayErrorCode = 50025
	GatewayErrorCodeMissingRequiredOauth2Scope                           GatewayErrorCode = 50026
	GatewayErrorCodeInvalidWebhookToken                                  GatewayErrorCode = 50027
	GatewayErrorCodeInvalidRole                                          GatewayErrorCode = 50028
	GatewayErrorCodeInvalidRecipients                                    GatewayErrorCode = 50033
	GatewayErrorCodeMessageTooOldForBulkDelete                           GatewayErrorCode = 50034
	GatewayErrorCodeInvalidFormBody                                      GatewayErrorCode = 50035
	GatewayErrorCodeInviteAcceptedForGuildBotNotIn                       GatewayErrorCode = 50036
	GatewayErrorCodeInvalidActivityAction                                GatewayErrorCode = 50039
	GatewayErrorCodeInvalidAPIVersion                                    GatewayErrorCode = 50041
	GatewayErrorCodeFileUploadExceedsMaxSize                             GatewayErrorCode = 50045
	GatewayErrorCodeInvalidFileUploaded                                  GatewayErrorCode = 50046
	GatewayErrorCodeCannotSelfRedeemGift                                 GatewayErrorCode = 50054
	GatewayErrorCodeInvalidGuild                                         GatewayErrorCode = 50055
	GatewayErrorCodeInvalidRequestOrigin                                 GatewayErrorCode = 50067
	GatewayErrorCodeInvalidMessageType                                   GatewayErrorCode = 50068
	GatewayErrorCodePaymentSourceRequiredToRedeemGift                    GatewayErrorCode = 50070
	GatewayErrorCodeCannotModifySystemWebhook                            GatewayErrorCode = 50073
	GatewayErrorCodeCannotDeleteCommunityRequiredChannel                 GatewayErrorCode = 50074
	GatewayErrorCodeCannotEditStickersWithinMessage                      GatewayErrorCode = 50080
	GatewayErrorCodeInvalidStickerSent                                   GatewayErrorCode = 50081
	GatewayErrorCodeOperationOnArchivedThread                            GatewayErrorCode = 50083
	GatewayErrorCodeInvalidThreadNotificationSettings                    GatewayErrorCode = 50084
	GatewayErrorCodeBeforeValueBeforeThreadCreation                      GatewayErrorCode = 50085
	GatewayErrorCodeCommunityChannelsMustBeTextChannels                  GatewayErrorCode = 50086
	GatewayErrorCodeEventEntityTypeMismatch                              GatewayErrorCode = 50091
	GatewayErrorCodeServerNotAvailableInYourLocation                     GatewayErrorCode = 50095
	GatewayErrorCodeServerNeedsMonetizationEnabled                       GatewayErrorCode = 50097
	GatewayErrorCodeServerNeedsMoreBoosts                                GatewayErrorCode = 50101
	GatewayErrorCodeRequestBodyContainsInvalidJSON                       GatewayErrorCode = 50109
	GatewayErrorCodeInvalidFileForAvatar                                 GatewayErrorCode = 50110
	GatewayErrorCodeCannotTransferOwnershipToBot                         GatewayErrorCode = 50123
	GatewayErrorCodeFailedToResizeAsset                                  GatewayErrorCode = 50132
	GatewayErrorCodeCannotMixSubscriptionAndNonSubscriptionRolesForEmoji GatewayErrorCode = 50138
	GatewayErrorCodeCannotConvertBetweenPremiumAndNormalEmoji            GatewayErrorCode = 50144
	GatewayErrorCodeUploadAssetExceedsMaxSize                            GatewayErrorCode = 50146
	GatewayErrorCodeCannotAddOrRemovePollVote                            GatewayErrorCode = 50178

	//6xxxx — authentication requirements
	GatewayErrorCodeTwoFactorRequired GatewayErrorCode = 60003

	//7xxxx — reactions
	GatewayErrorCodeReactionBlocked GatewayErrorCode = 70001

	//8xxxx — API overloaded
	GatewayErrorCodeAPIResourceOverloaded GatewayErrorCode = 80004

	//9xxxx — stage channels
	GatewayErrorCodeStageAlreadyOpen GatewayErrorCode = 90001

	//11xxxx — thread errors
	GatewayErrorCodeCannotReplyWithoutReadMessageHistoryPermission GatewayErrorCode = 110001
	GatewayErrorCodeThreadAlreadyCreatedForMessage                 GatewayErrorCode = 110002
	GatewayErrorCodeThreadIsLocked                                 GatewayErrorCode = 110003
	GatewayErrorCodeMaxActiveThreads                               GatewayErrorCode = 110004
	GatewayErrorCodeMaxActiveAnnouncementThreads                   GatewayErrorCode = 110005

	//13xxxx — Lottie/sticker validation
	GatewayErrorCodeInvalidLottieJSON                 GatewayErrorCode = 130000
	GatewayErrorCodeLottieCannotContainRasterImages   GatewayErrorCode = 130001
	GatewayErrorCodeStickerMaxFramerateExceeded       GatewayErrorCode = 130002
	GatewayErrorCodeStickerFrameCountExceedsMax       GatewayErrorCode = 130003
	GatewayErrorCodeLottieAnimationDimensionsExceeded GatewayErrorCode = 130004
	GatewayErrorCodeStickerFrameRateInvalid           GatewayErrorCode = 130005
	GatewayErrorCodeStickerAnimationDurationExceeds5s GatewayErrorCode = 130006

	//15xxxx — AutoMod
	GatewayErrorCodeMessageBlockedByAutoMod GatewayErrorCode = 150006

	//16xxxx — webhooks
	GatewayErrorCodeForumWebhookMissingRequiredFields GatewayErrorCode = 160002

	//18xxxx — channel permissions
	GatewayErrorCodeNoPermissionToSendInChannel GatewayErrorCode = 180000

	//20xxxx — ban failures
	GatewayErrorCodeFailedToBanUsers GatewayErrorCode = 200000
)
