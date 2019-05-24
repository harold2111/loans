package errors

const (
	IdentificationDuplicate = "IdentificationDuplicate"
	AddressNotExist         = "AddressNotExist"
	AddressDuplicate        = "AddressDuplicate"
	AddressFieldRequired    = "AddressFieldRequired"
	AddressToUpdateNotExist = "AddressToUpdateDoesNotExist"
	AtLeastOneAddress       = "AtLeastOneAddress"
	CityNotExist            = "CityNotExist"
	ClientNotAddressFound   = "ClientNotAddressFound"
	ClientNotExist          = "ClientNotExist"
	LoanNotExist            = "LoanNotExist"
	BillAlreadyExist        = "BillAlreadyExist"
	ToManyBillActives       = "ToManyBillActives"
	NotDataFound            = "NotDataFound"
	NotCitiesForDepartment  = "NotCitiesForDepartment"

	RequiredField = "RequiredField"
	InvalidField  = "InvalidField"
)

var ErrorMessages = map[string]string{
	IdentificationDuplicate: "Identification '%v' already exists",
	AddressNotExist:         "The address %v does not exist",
	AddressDuplicate:        "Address '%v' already exists",
	AtLeastOneAddress:       "At least one address is mandatory",
	AddressFieldRequired:    "Address field in Address is mandatory",
	AddressToUpdateNotExist: "Address '%v' to update does not exist",
	CityNotExist:            "City '%v' not exist",
	ClientNotAddressFound:   "No addresses found for client '%v'",
	ClientNotExist:          "Client '%v' not exist",
	LoanNotExist:            "Loan '%v' not exist",
	BillAlreadyExist:        "Loan bill already exist",
	ToManyBillActives:       "To many bill actives",
	NotDataFound:            "Not Data Found",
	NotCitiesForDepartment:  "There are no cities for department %v",

	RequiredField: "The field '%v' is mandatory",
	InvalidField:  "The field '%v' is invalid",
}
