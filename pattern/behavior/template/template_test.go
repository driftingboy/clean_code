package template

func Example_sms_Send() {
	tel := NewTelecomSms()
	_ = tel.SendTemplate("sms tele", 111111)
	//output:
	//sms valid...
	//send by telecom success
}
