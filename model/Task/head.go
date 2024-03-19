package Task

//type Item struct {
//	Id       string
//	Type     string
//	//CheckOut bool
//	Body     []byte
//}

type Item struct {
	Id   string
	Body []byte
}

type QQRegisterRecaptchaItem struct {
	Item
}

type QQRegisterRecaptchaSmsItem struct {
	Item
}
