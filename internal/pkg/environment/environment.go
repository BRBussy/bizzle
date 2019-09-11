package environment

const Development Environment = "Development"
const Production Environment = "Production"

type Environment string

func (e Environment) String() string {
	return string(e)
}
