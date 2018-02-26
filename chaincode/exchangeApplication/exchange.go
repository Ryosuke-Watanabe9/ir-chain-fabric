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
	BorrowNo          string `json:"borrowNo"`
	ProjectName       string `json:"projectName"`
	DataQuantity      string `json:"dataQuantity"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	OSRegistrant      string `json:"osRegistrant"`
	OSRechecker       string `json:"osRechecker"`
	OSReviewer        string `json:"osReviewer"`
	OSAuthorizer      string `json:"osAuthorizer"`
	IRRegistrant      string `json:"irRregistrant"`
	IRRechecker       string `json:"irRechecker"`
	IRReviewer        string `json:"irReviewer"`
	IRAuthorizer      string `json:"irAuthorizer"`
	DataContent       string `json:"dataContent"`
	ApplicationStatus string `json:"applicationStatus"`
	ApplicationType   string `json:"applicationType"`
}

type borrowData struct {
	BorrowNo string `json:"borrowNo"`
	Type     string `json:"type"`
	SystemNo string `json:"systemNo"`
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
	if function == "createApplication" {
		return s.createApplication(APIstub, args)
	} else if function == "queryApplication" {
		return s.queryApplication(APIstub, args)
	} else if function == "queryAllApplications" {
		return s.queryAllApplications(APIstub)
	} else if function == "changeApplicationStatus" {
		return s.changeApplicationStatus(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 14 {
		return shim.Error("Incorrect number of arguments. Expecting 14")
	}

	//get maxNumber of application
	maxNumberAsBytes, _ := APIstub.GetState("maxApplicationNo")

	maxNumber := MaxNumber{}
	json.Unmarshal(maxNumberAsBytes, &maxNumber)

	var tmpNumber int
	tmpNumber, _ = strconv.Atoi(maxNumber.MaxApplicationNo)
	tmpNumber = tmpNumber + 1

	maxNumber.MaxApplicationNo = strconv.Itoa(tmpNumber)

	var borrow = borrowData{BorrowNo: "999", Type: "borrow", SystemNo: "49"}

	var exchangeApplication = exchangeApplication{
		ApplicationNo:     maxNumber.MaxApplicationNo,
		BorrowNo:          args[0],
		ProjectName:       args[1],
		DataQuantity:      args[2],
		StartDate:         args[3],
		EndDate:           args[4],
		OSRegistrant:      args[5],
		OSRechecker:       args[6],
		OSReviewer:        args[7],
		OSAuthorizer:      args[8],
		IRRegistrant:      args[9],
		IRRechecker:       args[10],
		IRReviewer:        args[11],
		IRAuthorizer:      args[12],
		DataContent:       args[13],
		ApplicationStatus: "20",
		ApplicationType:   "exchange",
	}

	maxNumberAsBytes, _ = json.Marshal(maxNumber)
	exchangeApplicationAsBytes, _ := json.Marshal(exchangeApplication)
	borrowAsBytes, _ := json.Marshal(borrow)

	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)
	APIstub.PutState(maxNumber.MaxApplicationNo, exchangeApplicationAsBytes)
	APIstub.PutState(maxNumber.MaxApplicationNo, borrowAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	exchangeApplication := exchangeApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &exchangeApplication)
	exchangeApplication.ApplicationStatus = args[1]

	changeApplicationStatusAsBytes, _ = json.Marshal(exchangeApplication)
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

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
