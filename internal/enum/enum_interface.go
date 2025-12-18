package enum

// интерфейс для строковых enum
type EnumInterface interface {
	IsValid() bool
	String() string
	GetDefault() string
}
