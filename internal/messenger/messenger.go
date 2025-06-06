package messenger

type IMessenger interface {
	SendEmail(string) error
}

type Messenger struct {
}

func New() *Messenger {
	return &Messenger{}
}

func (m *Messenger) SendEmail(string) error {
	return nil
}
