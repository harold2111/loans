package errors

const (
	IdentificationDuplicate = "IdentificationDuplicate"
	AddressDuplicate        = "AddressDuplicate"
	AddressFieldRequired    = "AddressFieldRequired"
	CityNotExist            = "CityNotExist"
	ClientNotAddressFound   = "AddressNotExist"
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
	AddressDuplicate:        "Address '%v' already exists",
	AddressFieldRequired:    "Address field in Address is mandatory",
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
