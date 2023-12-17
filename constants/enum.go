package constants

type Environment int

const (
	EnvProduction Environment = iota
	EnvDevelopment
)

func (e Environment) String() string {
	switch e {
	case EnvProduction:
		return "production"
	case EnvDevelopment:
		return "development"
	default:
		return "development"
	}
}
