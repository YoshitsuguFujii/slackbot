package models

type SlackMessage struct {
	Text       string `json:"text"`       //投稿内容
	Username   string `json:"username"`   //投稿者名 or Bot名（存在しなくてOK）
	Icon_emoji string `json:"icon_emoji"` //アイコン絵文字
	Icon_url   string `json:"icon_url"`   //アイコンURL（icon_emojiが存在する場合は、適応されない）
	Channel    string `json:"channel"`    //部屋名
}

type SlackMessageAttachment struct {
	Username   string `json:"username"`   //投稿者名 or Bot名（存在しなくてOK）
	Icon_emoji string `json:"icon_emoji"` //アイコン絵文字
	Icon_url   string `json:"icon_url"`   //アイコンURL（icon_emojiが存在する場合は、適応されない）
	Channel    string `json:"channel"`    //部屋名
}

type SlackAttachment struct {
	Fallback string `json:"fallback"`
	Color    string `json:"color"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}
