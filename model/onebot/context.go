package onebot

// OneBot协议消息对象，注释说明仅供参考，在各事件类型中含义不一样，以官方文档为准
type Context struct {
	Time          int    `json:"time"`            // 通用，事件发生时间戳
	SelfId        int    `json:"self_id"`         // 通用，收到事件的机器人QQ号
	PostType      string `json:"post_type"`       // 通用， 上报消息类型
	MessageType   string `json:"message_type"`    // 消息事件，具体消息类型
	NoticeType    string `json:"notice_type"`     // 通知事件，具体通知类型
	RequestType   string `json:"request_type"`    // 请求事件，具体请求类型
	MetaEventType string `json:"meta_event_type"` // 元事件，具体请求类型
	HonorType     string `json:"honor_type"`      // 通知事件，群荣誉类型
	SubType       string `json:"sub_type"`        // 通用，消息子类型
	MessageId     int    `json:"message_id"`      // 通用，消息ID
	GroupId       int    `json:"group_id"`        // 通用，群号
	OperatorId    int    `json:"operator_id"`     // 通知事件，操作者QQ号
	UserId        int    `json:"user_id"`         // 通用，事件相关QQ号
	TargetId      int    `json:"target_id"`       // 通知事件，目标QQ号
	Comment       string `json:"comment"`         // 请求事件，验证信息
	Flag          string `json:"flag"`            // 请求事件，请求flag，调用处理请求API时传入
	Message       string `json:"message"`         // 消息事件，消息内容
	RawMessage    string `json:"raw_message"`     // 消息事件，原始消息内容
	Font          int    `json:"font"`            // 消息事件，字体
	Duration      int    `json:"duration"`        // 通知事件，禁言时长(秒)
	Status        string `json:"status"`          // 元事件，状态信息
	Interval      int    `json:"interval"`        // 元事件，到下次心跳事件时间间隔(毫秒)
	Anonymous     struct {
		Id   int    `json:"id"`   // 匿名用户ID
		Name string `json:"name"` // 匿名用户名称
		Flag string `json:"flag"` // 匿名用户flag，调用禁言API时传入
	} `json:"anonymous"` // 消息事件，匿名信息
	Sender struct {
		UserId   int    `json:"user_id"`  // 发送者QQ号
		Nickname string `json:"nickname"` // 昵称
		Card     string `json:"card"`     // 群名片、备注
		Sex      string `json:"sex"`      // 性别
		Age      int    `json:"age"`      // 年龄
		Area     string `json:"area"`     // 地区
		Level    string `json:"level"`    // 成员等级
		Role     string `json:"role"`     // 群内角色
		Title    string `json:"title"`    // 群内专属头衔
	} `json:"sender"` // 消息事件，发送人信息
	File struct {
		Id    string `json:"id"`    // 文件ID
		Name  string `json:"name"`  // 文件名
		Size  int    `json:"size"`  // 文件大小(字节数)
		BusId int    `json:"busid"` // busid
	} `json:"file"` // 通知事件，群文件信息

}
