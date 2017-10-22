package errors

const (
	IdentificationDuplicate = "IdentificationDuplicate"
	AddressRequired         = "AddressRequired"
	CityNotExist            = "CityNotExist"
	ClientNotAddressFound   = "AddressNotExist"
	ClientNotExist          = "ClientNotExist"

	RequiredField = "RequiredField"
	InvalidField  = "InvalidField"
)

var ErrorMessages = map[string]string{
	IdentificationDuplicate: "Identification '%v' already exists",
	AddressRequired:         "At least one address is required",
	CityNotExist:            "City '%v' not exist",
	ClientNotAddressFound:   "No addresses found for client '%v'",
	ClientNotExist:          "Client '%v' not exist",

	RequiredField: "The field '%v' is mandatory",
	InvalidField:  "The field '%v' is invalid",
}
