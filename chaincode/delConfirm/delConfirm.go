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

// delConfirm Chaincode implementation
type delConfirmApplication struct {
	ApplicationNo     string `json:"applicationNo"`
	Type              string `json:"type"`
	delConfirmNo      string `json:"delConfirmNo"`
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
	DataNum           string `json:"DataNum"`
	DelDate           string `json:"DelDate"`
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
        //      return shim.Success(maxNumber)
}
*/

//Init method is called as a result of deployment "delConfirm"
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	var maxNumber = MaxNumber{
		MaxApplicationNo: "0",
	}
	maxNumberAsBytes, _ := json.Marshal(maxNumber)
	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)

	return shim.Success(nil)
}

//Invoke method is called as a result of an application request to run the Smart Contract "delConfirm"
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "createDelConfirmApplication" {
		return s.createDelConfirmApplication(APIstub, args)
	} else if function == "queryDelConfirmApplication" {
		return s.queryDelConfirmApplication(APIstub, args)
	} else if function == "queryAllDelConfirmApplications" {
		return s.queryAllDelConfirmApplications(APIstub)
	} else if function == "changeDelConfirmApplicationStatus" {
		return s.changeDelConfirmApplicationStatus(APIstub, args)
	} else if function == "createStockManagement" {
		return s.createStockManagement(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createDelConfirmApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 17 {
		return shim.Error("Incorrect number of arguments. Expecting 17")
	}

	//get maxNumber of application
	maxNumberAsBytes, _ := APIstub.GetState("maxApplicationNo")

	maxNumber := MaxNumber{}
	json.Unmarshal(maxNumberAsBytes, &maxNumber)

	var tmpNumber int
	tmpNumber, _ = strconv.Atoi(maxNumber.MaxApplicationNo)
	tmpNumber = tmpNumber + 1

	maxNumber.MaxApplicationNo = strconv.Itoa(tmpNumber)

	var delConfirmApplication = delConfirmApplication{
		ApplicationNo:     maxNumber.MaxApplicationNo,
		delConfirmNo:      args[0],
		Type:              "delConfirm",
		SystemNo:          args[1],
		ProjectName:       args[2],
		DataQuantity:      args[3],
		ReportingDate:     args[4],
		StartDate:         args[5],
		EndDate:           args[6],
		ApplicationStatus: "30",
		IRRegistrant:      args[7],
		IRRechecker:       args[8],
		IRReviewer:        args[9],
		IRAuthorizer:      args[10],
		BKRegistrant:      args[11],
		BKRechecker:       args[12],
		BKReviewer:        args[13],
		BKAuthorizer:      args[14],
		DataNum:           args[15],
		DelDate:           args[16],
	}

	maxNumberAsBytes, _ = json.Marshal(maxNumber)
	delConfirmApplicationAsBytes, _ := json.Marshal(delConfirmApplication)

	APIstub.PutState("maxApplicationNo", maxNumberAsBytes)
	APIstub.PutState(maxNumber.MaxApplicationNo, delConfirmApplicationAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeDelConfirmApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// args[0] is applicationNo, args[1] is applicationStatus
	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	delConfirmApplication := delConfirmApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &delConfirmApplication)
	delConfirmApplication.ApplicationStatus = args[1]

	changeApplicationStatusAsBytes, _ = json.Marshal(delConfirmApplication)

	// update application status
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

	// applicationStatusが13(BK承認完了)となった場合、
	// stockManagementを変更
	if args[1] == "13" {

		// args[2] is systemNo
		stockManagementAsBytes, _ := APIstub.GetState(args[2])
		stockManagement := stockManagement{}

		json.Unmarshal(stockManagementAsBytes, &stockManagement)

		// DeleteAmountの加算
		// args[3] is deleteAmount
		var tmpDeleteAmount int
		var tmpAddDeleteAmount int
		tmpDeleteAmount, _ = strconv.Atoi(stockManagement.DeleteAmount)
		tmpAddDeleteAmount, _ = strconv.Atoi(args[3])
		tmpDeleteAmount = tmpDeleteAmount + tmpAddDeleteAmount

		stockManagement.DeleteAmount = strconv.Itoa(tmpDeleteAmount)

		// StockAmountの算出
		var tmpBorrowAmount int
		var tmpStockAmount int
		tmpBorrowAmount, _ = strconv.Atoi(stockManagement.StockAmount)
		tmpStockAmount = tmpBorrowAmount - tmpDeleteAmount

		stockManagement.StockAmount = strconv.Itoa(tmpStockAmount)

		stockManagementAsBytes, _ = json.Marshal(stockManagement)

		// update application status
		APIstub.PutState(args[0], stockManagementAsBytes)
	}

	// if exchangeApplication is completed, stockManagement write
	//if APIstub.GetState("SystemNo:"+args[2]) != NULL {
	//	APIstub.PutState("SystemNo:"+args[2], stockManagementAsBytes)
	//}

	return shim.Success(nil)
}

func (s *SmartContract) queryDelConfirmApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	delConfirmApplicationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(delConfirmApplicationAsBytes)
}

func (s *SmartContract) queryAllDelConfirmApplications(APIstub shim.ChaincodeStubInterface) sc.Response {

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

	fmt.Printf("- queryAllDelConfirmApplications:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
func (s *SmartContract) createStockManagement(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 授受票のOS承認が完了したら、有高管理する
	stockManagementAsBytes, err := APIstub.GetState("SystemNo:" + args[0])

	if stockManagementAsBytes == nil {
		// 有高管理を初めて行うシステムの場合
		var stockManagement = stockManagement{
			SystemNo:      args[1],
			BorrowAmount:  args[2],
			DeleteAmount:  "0",
			StockAmount:   args[3],
			OverDueAmount: "0",
		}
		// 有高管理対象のシステムを追加
		stockManagementAsBytes, _ := json.Marshal(stockManagement)
		APIstub.PutState("SystemNo:"+args[0], stockManagementAsBytes)
	}
	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
