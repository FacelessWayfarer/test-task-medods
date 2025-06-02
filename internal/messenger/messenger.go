package messenger

type Messenger struct {
}

func (m *Messenger) SendEmail(msg string) error {
	_ = msg
	return nil
}
