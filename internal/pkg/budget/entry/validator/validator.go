package validator

// Validator provides validation functions for budget entries
type Validator interface {
	ValidateForCreate(*ValidateForCreateRequest) (*ValidateForCreateResponse, error)
	ValidateForUpdate(*ValidateForUpdateRequest) (*ValidateForUpdateResponse, error)
}

// ValidateForCreateRequest is the request object for the ValidateForCreate service
type ValidateForCreateRequest struct {

}

// ValidateForCreateResponse is the response object for the ValidateForCreate service
type ValidateForCreateResponse struct {

}

// ValidateForUpdateRequest is the request object for the ValidateForUpdate service
type ValidateForUpdateRequest struct {

}

// ValidateForUpdateResponse is the response object for the ValidateForUpdate service
type ValidateForUpdateResponse struct {

}
