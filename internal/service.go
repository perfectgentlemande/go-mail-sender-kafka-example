package service

type Letter struct {
	EmailAddress string
	Contents     string
}

type LettersBroker interface {
	GetLetter() (Letter, error)
}
