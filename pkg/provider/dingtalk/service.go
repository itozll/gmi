package dingtalk

type Handler interface {
	Send(typ string, body interface{}) error
}

type service struct {
	token string
}

func New(token string) Handler {
	if len(token) == 0 {
		return defaultService
	}

	return &service{
		token: token,
	}
}

func (svc *service) Send(typ string, body interface{}) error {
	return Send(svc.token, typ, body)
}
