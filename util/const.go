package util

const (
	OneBotStatusOK         = 200
	OneBotNoQuickReply     = 204
	OneBotTextFormatError  = 400
	OneBotContentTypeError = 406
	OneBotAPINotFound      = 404
	OneBotTokenEmpty       = 401
	OneBotTokenWrong       = 403

	OneBotMessageEvent = "message"
	OneBotNoticeEvent  = "notice"
	OneBotRequestEvent = "request"
	OneBotMetaEvent    = "meta_event"

	OneBotMessagePrivate      = "private"
	OneBotMessageGroup        = "group"
	OneBotSubMessageFriend    = "friend"
	OneBotSubMessageGroup     = "group"
	OneBotSubMessageOther     = "other"
	OneBotSubMessageNormal    = "normal"
	OneBotSubMessageAnonymous = "anonymous"
	OneBotSubMessageNotice    = "notice"
	OneBotSenderSexMale       = "male"
	OneBotSenderSexFemale     = "female"
	OneBotSenderSexUnknown    = "unknown"
	OneBotSenderRoleOwner     = "owner"
	OneBotSenderRoleAdmin     = "admin"
	OneBotSenderRoleMember    = "member"

	OneBotNoticeGroupUpload   = "group_upload"
	OneBotNoticeGroupAdmin    = "group_admin"
	OneBotNoticeGroupDecrease = "group_decrease"
	OneBotNoticeGroupIncrease = "group_increase"
	OneBotNoticeGroupBan      = "group_ban"
	OneBotNoticeFriendAdd     = "friend_add"
	OneBotNoticeGroupRecall   = "group_recall"
	OneBotNoticeFriendRecall  = "friend_recall"
	OneBotNoticeNotify        = "notify"
	OneBotSubNoticeSetAdmin   = "set"
	OneBotSubNoticeUnsetAdmin = "unset"
	OneBotSubNoticeLeaveGroup = "leave"
	OneBotSubNoticeKickMember = "kick"
	OneBotSubNoticeKickMe     = "kick_me"
	OneBotSubNoticeApprove    = "approve"
	OneBotSubNoticeInvite     = "invite"
	OneBotSubNoticeBan        = "ban"
	OneBotSubNoticeLiftBan    = "lift_ban"
	OneBotSubNoticePoke       = "poke"
	OneBotSubNoticeLuckyKing  = "lucky_king"
	OneBotSubNoticeHonor      = "honor"
	OneBotHonorTalkative      = "talkative"
	OneBotHonorPerformer      = "performer"
	OneBotHonorEmotion        = "emotion"

	OneBotRequestFriend    = "friend"
	OneBotRequestGroup     = "group"
	OneBotSubRequestAdd    = "add"
	OneBotSubRequestInvite = "invite"

	OneBotMetaEventLifecycle = "lifecycle"
	OneBotMetaEventHeartbeat = "heartbeat"
	OneBotSubMetaEnable      = "enable"
	OneBotSubMetaDisable     = "disable"
	OneBotSubMetaConnect     = "connect"
	OneBotSubMetaHeartbeat   = "heartbeat"

	FoodGameDataURL      = "https://foodgame.gitee.io/data/data.min.json"
	FoodGameImageCSSURL  = "https://foodgame.gitee.io/css/image.css"
	ChefImageRetinaURL   = "https://foodgame.gitee.io/images/chef_retina.png"
	RecipeImageRetinaURL = "https://foodgame.gitee.io/images/recipe_retina.png"
	EquipImageRetinaURL  = "https://foodgame.gitee.io/images/equip_retina.png"

	ArgsSplitCharacter   = " "
	ArgsConnectCharacter = "-"
	MaxQueryListLength   = 10

	QueryParamWrongNote  = "查询格式错了哦"
	SystemErrorNote      = "唔，系统开小差啦"
	PermissionDeniedNote = "我不听我不听我不听!"
)

var PrefixCharacters = []string{"#", "＃"}
