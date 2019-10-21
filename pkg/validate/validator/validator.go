package validator

type Validator interface {
	ValidateRequest(interface{}) error
}
