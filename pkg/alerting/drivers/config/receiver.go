package config

/*
Contains the specifications of transactions on an AlertManager config.
*/

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aity-cloud/monty/pkg/alerting/shared"
	"github.com/aity-cloud/monty/pkg/alerting/templates"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	"github.com/aity-cloud/monty/pkg/util"
	amCfg "github.com/prometheus/alertmanager/config"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

const missingBody = "<missing body>"
const missingTitle = "<missing title>"

// Extends the receiver configs of AlertManager, e.g. SlackConfig, EmailConfig...
// panics when the receiver type
type MontyReceiver interface {
	InternalId() string
	// extract non-specific receiver info
	ExtractInfo() *alertingv1.EndpointImplementation
	// set non-specific receiver info
	StoreInfo(details *alertingv1.EndpointImplementation)
	// configure receiver specific info
	Configure(*alertingv1.AlertEndpoint) MontyReceiver
	Clone() MontyReceiver
	yaml.Unmarshaler
}

func ExtractReceiver(unmarshall func(interface{}) error, _ /*data*/ interface{}) (MontyReceiver, error) {
	type slack SlackConfig
	type email EmailConfig
	type pagerduty PagerdutyConfig
	type webhook WebhookConfig
	type opsGenie OpsGenieConfig
	type victorOps VictorOpsConfig
	type wechat WechatConfig
	type pushover PushoverConfig
	type sns SNSConfig
	type telegram TelegramConfig

	slackCfg := &SlackConfig{}
	emailCfg := &EmailConfig{}
	pagerdutyCfg := &PagerdutyConfig{}
	webhookCfg := &WebhookConfig{}
	opsGenieCfg := &OpsGenieConfig{}
	victorOpsCfg := &VictorOpsConfig{}
	wechatCfg := &WechatConfig{}
	pushoverCfg := &PushoverConfig{}
	snsCfg := &SNSConfig{}
	telegramCfg := &TelegramConfig{}
	if err := unmarshall((*slack)(slackCfg)); err == nil {
		return slackCfg, nil
	}
	if err := unmarshall((*email)(emailCfg)); err == nil {
		return emailCfg, nil
	}
	if err := unmarshall((*pagerduty)(pagerdutyCfg)); err == nil {
		return pagerdutyCfg, nil
	}
	if err := unmarshall((*webhook)(webhookCfg)); err == nil {
		return webhookCfg, nil
	}
	if err := unmarshall((*opsGenie)(opsGenieCfg)); err == nil {
		return opsGenieCfg, nil
	}
	if err := unmarshall((*victorOps)(victorOpsCfg)); err == nil {
		return victorOpsCfg, nil
	}
	if err := unmarshall((*wechat)(wechatCfg)); err == nil {
		return wechatCfg, nil
	}
	if err := unmarshall((*pushover)(pushoverCfg)); err == nil {
		return pushoverCfg, nil
	}
	if err := unmarshall((*sns)(snsCfg)); err == nil {
		return snsCfg, nil
	}
	if err := unmarshall((*telegram)(telegramCfg)); err == nil {
		return telegramCfg, nil
	}
	return nil, fmt.Errorf("unknown receiver type")
}

// Takes a collection of MontyReceivers and converts them to a single AlertManager
// this function will panic if it cannot convert the receivers
func BuildReceiver(receiverId string, recvs []MontyReceiver) (Receiver, error) {
	var r Receiver
	slackCfg := []*SlackConfig{}
	emailCfg := []*EmailConfig{}
	pagerdutyCfg := []*PagerdutyConfig{}
	webhookCfg := []*WebhookConfig{}
	opsgenieCfg := []*OpsGenieConfig{}
	victoropsCfg := []*VictorOpsConfig{}
	wechatCfg := []*WechatConfig{}
	pushoverCfg := []*PushoverConfig{}
	snsCfg := []*SNSConfig{}
	telegramCfg := []*TelegramConfig{}

	if len(recvs) == 0 {
		return r, fmt.Errorf("no receivers to build")
	}

	for _, recv := range recvs {
		switch recv.InternalId() {
		case shared.InternalSlackId:
			slackCfg = append(slackCfg, recv.(*SlackConfig))
		case shared.InternalEmailId:
			emailCfg = append(emailCfg, recv.(*EmailConfig))
		case shared.InternalPagerdutyId:
			pagerdutyCfg = append(pagerdutyCfg, recv.(*PagerdutyConfig))
		case shared.InternalWebhookId:
			webhookCfg = append(webhookCfg, recv.(*WebhookConfig))
		case shared.InternalOpsGenieId:
			fallthrough
		case shared.InternalVictorOpsId:
			fallthrough
		case shared.InternalWechatId:
			fallthrough
		case shared.InternalPushoverId:
			fallthrough
		case shared.InternalSNSId:
			fallthrough
		case shared.InternalTelegramId:
			fallthrough
		default:
			return r, fmt.Errorf("unknown receiver type %s", recv.InternalId())
		}
	}
	if len(slackCfg)+len(emailCfg)+len(pagerdutyCfg)+len(webhookCfg) == 0 {
		return r, fmt.Errorf("no receivers to configs parsed")
	}

	return Receiver{
		Name:             receiverId,
		SlackConfigs:     slackCfg,
		EmailConfigs:     emailCfg,
		PagerdutyConfigs: pagerdutyCfg,
		WebhookConfigs:   webhookCfg,
		OpsGenieConfigs:  opsgenieCfg,
		VictorOpsConfigs: victoropsCfg,
		WechatConfigs:    wechatCfg,
		PushoverConfigs:  pushoverCfg,
		SNSConfigs:       snsCfg,
		TelegramConfigs:  telegramCfg,
	}, nil
}

var _ MontyReceiver = (*EmailConfig)(nil)
var _ MontyReceiver = (*SlackConfig)(nil)
var _ MontyReceiver = (*PagerdutyConfig)(nil)
var _ MontyReceiver = (*WebhookConfig)(nil)
var _ MontyReceiver = (*OpsGenieConfig)(nil)
var _ MontyReceiver = (*VictorOpsConfig)(nil)
var _ MontyReceiver = (*WechatConfig)(nil)
var _ MontyReceiver = (*PushoverConfig)(nil)
var _ MontyReceiver = (*SNSConfig)(nil)
var _ MontyReceiver = (*TelegramConfig)(nil)

// var _ MontyReceiver = (*DiscordConfig)(nil)
// var _ MontyReceiver = (*WebexConfig)(nil)

// --- MontyReceiver Implementations

func (c *EmailConfig) InternalId() string {
	return shared.InternalEmailId
}

func (c *EmailConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	res := &alertingv1.EndpointImplementation{}
	res.SendResolved = &c.VSendResolved
	return res
}

func (c *EmailConfig) StoreInfo(details *alertingv1.EndpointImplementation) {
	if def := details.SendResolved; def != nil {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: *def,
		}
	} else {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: false,
		}
	}
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	c.Headers["Subject"] = details.Body
	c.HTML = details.Body
}

func (c *EmailConfig) Configure(endp *alertingv1.AlertEndpoint) MontyReceiver {
	emailSpec := endp.GetEmail()
	c.To = emailSpec.To
	if emailSpec.SmtpFrom != nil {
		c.From = *emailSpec.SmtpFrom
	}
	if emailSpec.SmtpSmartHost != nil {
		arr := strings.Split(*emailSpec.SmtpSmartHost, ":")
		c.Smarthost = HostPort{
			Host: arr[0],
			Port: arr[1],
		}
	}
	if emailSpec.SmtpAuthUsername != nil {
		c.AuthUsername = *emailSpec.SmtpAuthUsername
	}
	if emailSpec.SmtpAuthPassword != nil {
		c.AuthPassword = *emailSpec.SmtpAuthPassword
	}
	if emailSpec.SmtpAuthIdentity != nil {
		c.AuthIdentity = *emailSpec.SmtpAuthIdentity
	}
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	c.Headers["Subject"] = templates.HeaderTemplate()
	c.HTML = templates.BodyTemplate()
	return c
}

func (c *EmailConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *SlackConfig) InternalId() string {
	return shared.InternalSlackId
}

func (c *SlackConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	res := &alertingv1.EndpointImplementation{}
	res.SendResolved = &c.VSendResolved
	return res
}

func (c *SlackConfig) StoreInfo(details *alertingv1.EndpointImplementation) {
	if def := details.SendResolved; def != nil {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: *def,
		}
	} else {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: false,
		}
	}
}

func (c *SlackConfig) Configure(endp *alertingv1.AlertEndpoint) MontyReceiver {
	slackSpec := endp.GetSlack()
	parsedURL, err := url.Parse(slackSpec.WebhookUrl)
	if err != nil {
		panic(err)
	}
	c.APIURL = &amCfg.URL{
		URL: parsedURL,
	}
	c.Channel = slackSpec.Channel
	c.Title = templates.HeaderTemplate()
	c.Text = templates.BodyTemplate()
	return c
}

func (c *SlackConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *PagerdutyConfig) InternalId() string {
	return shared.InternalPagerdutyId
}

func (c *PagerdutyConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	res := &alertingv1.EndpointImplementation{}
	res.SendResolved = &c.VSendResolved
	return res
}

func (c *PagerdutyConfig) StoreInfo(details *alertingv1.EndpointImplementation) {
	if def := details.SendResolved; def != nil {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: *def,
		}
	} else {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: false,
		}
	}
}

func (c *PagerdutyConfig) Configure(endp *alertingv1.AlertEndpoint) MontyReceiver {
	pagerdutySpec := endp.GetPagerDuty()
	if pagerdutySpec.ServiceKey != "" {
		c.ServiceKey = pagerdutySpec.ServiceKey
	}
	if pagerdutySpec.IntegrationKey != "" {
		c.ServiceKey = pagerdutySpec.IntegrationKey
	}
	if c.Details == nil {
		c.Details = map[string]string{}
	}
	c.Details = lo.Assign(c.Details,
		DefaultPagerdutyDetails,
		map[string]string{
			"title": templates.HeaderTemplate(),
		})
	c.Description = templates.HeaderTemplate() + "\n" + templates.BodyTemplate()
	return c
}

func (c *PagerdutyConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *WebhookConfig) InternalId() string {
	return shared.InternalWebhookId
}

func (c *WebhookConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	res := &alertingv1.EndpointImplementation{}
	res.SendResolved = &c.VSendResolved
	return res
}

func (c *WebhookConfig) StoreInfo(details *alertingv1.EndpointImplementation) {
	if def := details.SendResolved; def != nil {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: *def,
		}
	} else {
		c.NotifierConfig = NotifierConfig{
			VSendResolved: false,
		}
	}
}

func (c *WebhookConfig) Configure(endp *alertingv1.AlertEndpoint) MontyReceiver {
	webhookSpec := endp.GetWebhook()
	return ToWebhook(webhookSpec)
}

func (c *WebhookConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *OpsGenieConfig) InternalId() string {
	return shared.InternalOpsGenieId
}

func (c *OpsGenieConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *OpsGenieConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *OpsGenieConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *OpsGenieConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *VictorOpsConfig) InternalId() string {
	return shared.InternalVictorOpsId
}

func (c *VictorOpsConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *VictorOpsConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *VictorOpsConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *VictorOpsConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *WechatConfig) InternalId() string {
	return shared.InternalWechatId
}

func (c *WechatConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *WechatConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *WechatConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *WechatConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *PushoverConfig) InternalId() string {
	return shared.InternalPushoverId
}

func (c *PushoverConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *PushoverConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *PushoverConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *PushoverConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *SNSConfig) InternalId() string {
	return shared.InternalSNSId
}

func (c *SNSConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *SNSConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *SNSConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *SNSConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}

func (c *TelegramConfig) InternalId() string {
	return shared.InternalTelegramId
}

func (c *TelegramConfig) ExtractInfo() *alertingv1.EndpointImplementation {
	//TODO
	return nil
}

func (c *TelegramConfig) StoreInfo(_ *alertingv1.EndpointImplementation) {
	//TODO
}

func (c *TelegramConfig) Configure(*alertingv1.AlertEndpoint) MontyReceiver {
	return c
}

func (c *TelegramConfig) Clone() MontyReceiver {
	return util.DeepCopy(c)
}
