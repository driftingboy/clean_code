package template

import "fmt"

// ISMS ISMS
type ISMS interface {
	send(content string, phone int) error
}

// SMS 短信发送基类
type sms struct {
	ISMS
}

// Send 发送短信
func (s *sms) SendTemplate(content string, phone int) error {
	if err := s.valid(content); err != nil {
		return err
	}
	// do something...

	// 调用子类的方法发送短信
	return s.send(content, phone)
}

// Valid 校验短信字数
func (s *sms) valid(content string) error {
	fmt.Println("sms valid...")
	if len(content) > 63 {
		return fmt.Errorf("content is too long")
	}
	return nil
}

// TelecomSms 走电信通道
type TelecomSms struct {
	*sms
}

// NewTelecomSms NewTelecomSms
func NewTelecomSms() *TelecomSms {
	tel := &TelecomSms{}
	tel.sms = &sms{ISMS: tel}
	return tel
}

func (tel *TelecomSms) send(content string, phone int) error {
	fmt.Println("send by telecom success")
	return nil
}
