package auth

type Handler interface {
	Type() string
	Check(token string) (*AuthUser, error)
}

type Handlers []Handler
