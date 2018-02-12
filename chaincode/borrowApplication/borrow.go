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

// borrowApplication Chaincode implementation
type borrowApplication struct {
	ApplicationNo     string `json:"applicationNo"`
	SystemNo          string `json:"systemNo"`
	ProjectName       string `json:"projectName"`
	DataQuantity      string `json:"dataQuantity"`
	ReportingDate     string `json:"reportingDate"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	ApplicationStatus string `json:"applicationStatus"`
	IRRegistrant      string `json:"irRegistrant"`
	IRRechecker       string `json:"irRechecker"`
	IRReviewer        string `json:"irReviewer"`
	IRAuthorizer      string `json:"irAuthorizer"`
	BKRegistrant      string `json:"bkRregistrant"`
	BKRechecker       string `json:"bkRechecker"`
	BKReviewer        string `json:"bkReviewer"`
	BKAuthorizer      string `json:"bkAuthorizer"`
	EvidenceValue     string `json:"evidenceValue"`
}

/*
//Numbering is a function to decide applicationNo
func (s *SmartContract) Numbering(APIstub shim.ChaincodeStubInterface) sc.Response {

	//	return shim.Success(maxNumber)
}
*/

//Init method is called as a result of deployment "borrowApplication"
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	var maxNumber = MaxNumber{
		MaxApplicationNo: "0",
	}
	maxNumberAsBytes, _ := json.Marshal(maxNumber)
	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)

	return shim.Success(nil)
}

//Invoke method is called as a result of an application request to run the Smart Contract "borrowApplication"
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

	if len(args) != 15 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	//get maxNumber of application
	maxNumberAsBytes, _ := APIstub.GetState("maxApplicationNo")

	maxNumber := MaxNumber{}
	json.Unmarshal(maxNumberAsBytes, &maxNumber)

	var tmpNumber int
	tmpNumber, _ = strconv.Atoi(maxNumber.MaxApplicationNo)
	tmpNumber = tmpNumber + 1

	maxNumber.MaxApplicationNo = strconv.Itoa(tmpNumber)

	var borrowApplication = borrowApplication{
		ApplicationNo:     maxNumber.MaxApplicationNo,
		SystemNo:          args[0],
		ProjectName:       args[1],
		DataQuantity:      args[2],
		ReportingDate:     args[3],
		StartDate:         args[4],
		EndDate:           args[5],
		IRRegistrant:      args[6],
		IRRechecker:       args[7],
		IRReviewer:        args[8],
		IRAuthorizer:      args[9],
		BKRegistrant:      args[10],
		BKRechecker:       args[11],
		BKReviewer:        args[12],
		BKAuthorizer:      args[13],
		EvidenceValue:     args[14],
		ApplicationStatus: "0",
	}

	maxNumberAsBytes, _ = json.Marshal(maxNumber)
	borrowApplicationAsBytes, _ := json.Marshal(borrowApplication)

	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)
	APIstub.PutState(maxNumber.MaxApplicationNo, borrowApplicationAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	borrowApplication := borrowApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &borrowApplication)
	borrowApplication.ApplicationStatus = args[1]

	changeApplicationStatusAsBytes, _ = json.Marshal(borrowApplication)
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	borrowApplicationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(borrowApplicationAsBytes)
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
