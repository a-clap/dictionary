package translate

type Language int64

const (
	Polish Language = iota
	English
)

type Translate interface {
	Translate(word string, lang Language) (translation string)
}
