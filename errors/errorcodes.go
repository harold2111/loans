package errors

const (
	IdentificationDuplicate = "IdentificationDuplicate"
	AddressDuplicate        = "AddressDuplicate"
	AddressRequired         = "AddressRequired"
	CityNotExist            = "CityNotExist"
	ClientNotAddressFound   = "AddressNotExist"
	ClientNotExist          = "ClientNotExist"
	LoanNotExist            = "LoanNotExist"
	BillAlreadyExist        = "BillAlreadyExist"
	ToManyBillActives       = "ToManyBillActives"
	NotDataFound            = "NotDataFound"

	RequiredField = "RequiredField"
	InvalidField  = "InvalidField"
)

var ErrorMessages = map[string]string{
	IdentificationDuplicate: "Identification '%v' already exists",
	AddressDuplicate:        "Address '%v' already exists",
	AddressRequired:         "Address is required",
	CityNotExist:            "City '%v' not exist",
	ClientNotAddressFound:   "No addresses found for client '%v'",
	ClientNotExist:          "Client '%v' not exist",
	LoanNotExist:            "Loan '%v' not exist",
	BillAlreadyExist:        "Loan bill already exist",
	ToManyBillActives:       "To many bill actives",
	NotDataFound:            "Not Data Found",

	RequiredField: "The field '%v' is mandatory",
	InvalidField:  "The field '%v' is invalid",
}
