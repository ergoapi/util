package bark

import (
	"fmt"
	"strings"

	"github.com/ergoapi/util/exhttp"
)

type Level string

const (
	DefaultLevel       Level = "active"        // 默认, 系统立即亮屏
	TimeSensitiveLevel Level = "timeSensitive" // 时效性, 可在专注模式下显示通知
	PassiveLevel       Level = "passive"       // 系统不亮屏, 仅放到通知列表
)

type Bark struct {
	Client    *exhttp.Client `json:"-"`          // http client
	APIUrl    string         `json:"-"`          // bark url
	Title     string         `json:"title"`      // 标题
	Body      string         `json:"body"`       // 内容
	DeviceKey string         `json:"device_key"` // 设备key

	// optional args
	Level             Level  `json:"level,omitempty"`             // 级别
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty"` // 是否自动复制, 设置值只能为1
	Copy              string `json:"copy,omitempty"`              // 复制内容
	Sound             string `json:"sound,omitempty"`             // 铃声
	Icon              string `json:"icon,omitempty"`              // 图标
	Group             string `json:"group,omitempty"`             // 分组
	IsArchive         string `json:"isArchive,omitempty"`         // 是否存档，设置值只能为1
	Url               string `json:"url,omitempty"`               // 点击跳转的url
}

type Core struct {
	Title string `json:"title"`           // 标题
	Body  string `json:"body"`            // 内容
	Url   string `json:"url,omitempty"`   // 点击跳转的url
	Copy  string `json:"copy,omitempty"`  // 复制内容
	Group string `json:"group,omitempty"` // 分组
}

type Result struct {
	Code      int64  `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (b *Bark) api() string {
	if strings.HasSuffix(b.APIUrl, "/") {
		return fmt.Sprintf("%spush", b.APIUrl)
	}
	return fmt.Sprintf("%s/push", b.APIUrl)
}

func (b *Bark) SendEvent(c Core) error {
	var res Result
	if len(c.Title) > 0 {
		b.Title = c.Title
	}
	if len(c.Body) > 0 {
		b.Body = c.Body
	}
	if len(c.Url) > 0 {
		b.Url = c.Url
	}
	if len(c.Group) > 0 {
		b.Group = c.Group
	}
	if len(c.Copy) > 0 {
		b.Copy = c.Copy
	}
	resp, err := b.Client.R().
		SetBody(b).
		SetSuccessResult(&res).
		Post(b.api())
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("bark send event failed, bad response status: %s", resp.Status)
	}
	return nil
}

func NewBark(url, device string, httpClient *exhttp.Client) (*Bark, error) {
	b := &Bark{
		Client:            httpClient,
		APIUrl:            url,
		Title:             "默认标题",
		Body:              "默认正文",
		DeviceKey:         device,
		Level:             DefaultLevel,
		AutomaticallyCopy: "1",
		Group:             "默认",
		IsArchive:         "1",
	}
	if b.Client == nil {
		b.Client, _ = exhttp.GetClient()
	}
	return b, nil
}
