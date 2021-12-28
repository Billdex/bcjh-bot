package onebot

const (
	PostTypeMessageEvent = "message"
	PostTypeNoticeEvent  = "notice"
	PostTypeRequestEvent = "request"
	PostTypeMetaEvent    = "meta_event"

	MessageTypePrivate = "private"
	MessageTypeGroup   = "group"

	NoticeTypeGroupRecall = "group_recall"

	GroupSenderRoleOwner  = "owner"
	GroupSenderRoleAdmin  = "admin"
	GroupSenderRoleMember = "member"
)

// OneBot协议消息对象，详细说明请参考onebot文档
// https://github.com/botuniverse/onebot/blob/master/v11/specs/communication/http-post.md

// 事件
// 消息事件
type MessageEventPrivateReq struct {
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	TempSource  int    `json:"temp_source"`
	MessageId   int32  `json:"message_id"`
	UserId      int64  `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int32  `json:"font"`
	Sender      struct {
		UserId   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
	} `json:"sender"`
}

type MessageEventGroupReq struct {
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageId   int32  `json:"message_id"`
	GroupId     int64  `json:"group_id"`
	UserId      int64  `json:"user_id"`
	Anonymous   struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`
	Message    string `json:"message"`
	RawMessage string `json:"raw_message"`
	Font       int32  `json:"font"`
	Sender     struct {
		UserId   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
		Area     string `json:"area"`
		Level    string `json:"level"`
		Role     string `json:"role"`
		Title    string `json:"title"`
	} `json:"sender"`
}

// 通知事件
type NoticeEventGroupRecall struct {
	Time       int64  `json:"time"`
	SelfId     int64  `json:"self_id"`
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
	GroupId    int64  `json:"group_id"`
	UserId     int64  `json:"user_id"`
	OperatorId int64  `json:"operator_id"`
	MessageId  int64  `json:"message_id"`
}

// 请求API
type actionApiReq struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type actionApiResp struct {
	Status  string      `json:"status"`
	RetCode int64       `json:"retcode"`
	Data    interface{} `json:"data"`
	Echo    string      `json:"echo"`
}

type sendPrivateMsgParams struct {
	UserId     int64  `json:"user_id"`
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type sendPrivateMsgResp struct {
	Status  string `json:"status"`
	RetCode int64  `json:"retcode"`
	Data    struct {
		MessageId int32 `json:"message_id"`
	} `json:"data"`
}

type sendGroupMsgParams struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type sendGroupMsgResp struct {
	Status  string `json:"status"`
	RetCode int64  `json:"retcode"`
	Data    struct {
		MessageId int32 `json:"message_id"`
	} `json:"data"`
}

type GroupInfo struct {
	GroupId         int64  `json:"group_id"`
	GroupName       string `json:"group_name"`
	GroupMemo       string `json:"group_memo"`
	GroupCreateTime uint32 `json:"group_create_time"`
	GroupLevel      uint32 `json:"group_level"`
	MemberCount     int32  `json:"member_count"`
	MaxMemberCount  int32  `json:"max_member_count"`
}

type getGroupInfoParams struct {
	GroupId int64 `json:"group_id"`
	NoCache bool  `json:"no_cache"`
}

type getGroupInfoResp struct {
	Status  string    `json:"status"`
	RetCode int64     `json:"retcode"`
	Data    GroupInfo `json:"data"`
	Echo    string    `json:"echo"`
}

type getGroupListResp struct {
	Status  string      `json:"status"`
	RetCode int64       `json:"retcode"`
	Data    []GroupInfo `json:"data"`
	Echo    string      `json:"echo"`
}

type MsgInfo struct {
	MessageId int32 `json:"message_id"`
	RealId    int32 `json:"real_id"`
	Sender    struct {
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
	}
	Time       int32  `json:"time"`
	Message    string `json:"message"`
	RawMessage string `json:"raw_message"`
}

type getMsgInfoParams struct {
	MessageId int32 `json:"message_id"`
}

type getMsgInfoResp struct {
	Status  string  `json:"status"`
	RetCode int64   `json:"retcode"`
	Data    MsgInfo `json:"data"`
	Echo    string  `json:"echo"`
}

type getImageParams struct {
	File string `json:"file"`
}

type ImageInfo struct {
	Size     int32  `json:"size"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}
