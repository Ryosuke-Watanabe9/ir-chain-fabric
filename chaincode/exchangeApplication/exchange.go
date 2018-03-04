package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract structure
type SmartContract struct {
}

// MaxNumber structure
type MaxNumber struct {
	MaxApplicationNo string `json:"maxApplicationNo"`
}

// exchangeApplication Chaincode implementation
type exchangeApplication struct {
	ApplicationNo     string `json:"applicationNo"`
	Type              string `json:"type"`
	BorrowNo          string `json:"borrowNo"`
	SystemNo          string `json:"systemNo"`
	ProjectName       string `json:"projectName"`
	DataQuantity      string `json:"dataQuantity"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	ApplicationStatus string `json:"applicationStatus"`
	OSRegistrant      string `json:"osRegistrant"`
	OSRechecker       string `json:"osRechecker"`
	OSReviewer        string `json:"osReviewer"`
	OSAuthorizer      string `json:"osAuthorizer"`
	IRRegistrant      string `json:"irRregistrant"`
	IRRechecker       string `json:"irRechecker"`
	IRReviewer        string `json:"irReviewer"`
	IRAuthorizer      string `json:"irAuthorizer"`
	DataContent       string `json:"dataContent"`
}

// stockManagement Chaincode implementation
type stockManagement struct {
	SystemNo      string `json:"systemNo"`
	BorrowAmount  string `json:"borrowAmount"`
	DeleteAmount  string `json:"deleteAmount"`
	StockAmount   string `json:"stockAmount"`
	OverDueAmount string `json:"overDueAmount"`
}

/*
//Numbering is a function to decide applicationNo
func (s *SmartContract) Numbering(APIstub shim.ChaincodeStubInterface) sc.Response {
	//	return shim.Success(maxNumber)
}
*/

//Init method is called as a result of deployment "exchangeApplication"
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	var maxNumber = MaxNumber{
		MaxApplicationNo: "0",
	}
	maxNumberAsBytes, _ := json.Marshal(maxNumber)
	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "exchangeApplication"
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "createExchangeApplication" {
		return s.createExchangeApplication(APIstub, args)
	} else if function == "queryApplication" {
		return s.queryApplication(APIstub, args)
	} else if function == "queryAllApplications" {
		return s.queryAllApplications(APIstub)
	} else if function == "changeExchangeApplicationStatus" {
		return s.changeExchangeApplicationStatus(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createExchangeApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// 15args retrieved
	if len(args) != 15 {
		return shim.Error("Incorrect number of arguments. Expecting 15")
	}

	// get maxNumber of application
	maxNumberAsBytes, _ := APIstub.GetState("maxApplicationNo")

	maxNumber := MaxNumber{}
	json.Unmarshal(maxNumberAsBytes, &maxNumber)

	var tmpNumber int
	tmpNumber, _ = strconv.Atoi(maxNumber.MaxApplicationNo)
	tmpNumber = tmpNumber + 1

	maxNumber.MaxApplicationNo = strconv.Itoa(tmpNumber)

	// exchange data define
	var exchangeApplication = exchangeApplication{
		ApplicationNo:     maxNumber.MaxApplicationNo,
		Type:              "exchange",
		BorrowNo:          args[0],
		SystemNo:          args[1],
		ProjectName:       args[2],
		DataQuantity:      args[3],
		StartDate:         args[4],
		EndDate:           args[5],
		ApplicationStatus: "20",
		OSRegistrant:      args[6],
		OSRechecker:       args[7],
		OSReviewer:        args[8],
		OSAuthorizer:      args[9],
		IRRegistrant:      args[10],
		IRRechecker:       args[11],
		IRReviewer:        args[12],
		IRAuthorizer:      args[13],
		DataContent:       args[14],
	}

	maxNumberAsBytes, _ = json.Marshal(maxNumber)
	exchangeApplicationAsBytes, _ := json.Marshal(exchangeApplication)

	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)
	APIstub.PutState(maxNumber.MaxApplicationNo, exchangeApplicationAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeExchangeApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// args[0] is applicationNo, args[1] is applicationStatus
	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	exchangeApplication := exchangeApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &exchangeApplication)
	exchangeApplication.ApplicationStatus = args[1]

	changeApplicationStatusAsBytes, _ = json.Marshal(exchangeApplication)

	// update application status
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

	// stockManagement data define
	var stockManagement = stockManagement{
		SystemNo:      args[2],
		BorrowAmount:  args[3],
		DeleteAmount:  args[4],
		StockAmount:   args[5],
		OverDueAmount: args[6],
	}

	stockManagementAsBytes, _ := json.Marshal(stockManagement)

	// if exchangeApplication is completed, stockManagement write
	if APIstub.GetState("SystemNo:"+args[2]) != NULL {
		APIstub.PutState("SystemNo:"+args[2], stockManagementAsBytes)
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	exchangeApplicationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(exchangeApplicationAsBytes)
}

func (s *SmartContract) queryAllApplications(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllApplications:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
