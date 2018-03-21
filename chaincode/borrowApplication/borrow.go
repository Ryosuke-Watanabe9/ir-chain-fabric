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
	Type              string `json:"type"`
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
	} else if function == "setApplicationRoute" {
		return s.setApplicationRoute(APIstub, args)
	} else if function == "queryStockManagement" {
		return s.queryStockManagement(APIstub, args)
	} else if function == "createExchangeApplication" {
		return s.createExchangeApplication(APIstub, args)
	} else if function == "changeExchangeApplicationStatus" {
		return s.changeExchangeApplicationStatus(APIstub, args)
	} else if function == "setExchangeRoute" {
		return s.setExchangeRoute(APIstub, args)
	} else if function == "queryExchangeApplication" {
		return s.queryExchangeApplication(APIstub, args)
	} else if function == "queryDelConfirmApplication" {
		return s.queryDelConfirmApplication(APIstub, args)
	} else if function == "createDelConfirmApplication" {
		return s.createDelConfirmApplication(APIstub, args)
	} else if function == "changeDelConfirmApplicationStatus" {
		return s.changeDelConfirmApplicationStatus(APIstub, args)
	} else if function == "setDelConfirmRoute" {
		return s.setDelConfirmRoute(APIstub, args)	
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 15 {
		return shim.Error("Incorrect number of arguments. Expecting 15")
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
		Type:              "borrow",
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

func (s *SmartContract) setApplicationRoute(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	borrowApplication := borrowApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &borrowApplication)
	borrowApplication.BKRegistrant = args[1]
	borrowApplication.BKRechecker = args[2]
	borrowApplication.BKReviewer = args[3]
	borrowApplication.BKAuthorizer = args[4]
	borrowApplication.ApplicationStatus = "10"

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

func (s *SmartContract) queryExchangeApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	exchangeApplicationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(exchangeApplicationAsBytes)
}

func (s *SmartContract) queryDelConfirmApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	delConfirmApplicationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(delConfirmApplicationAsBytes)
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

func (s *SmartContract) queryStockManagement(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stockManagementAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(stockManagementAsBytes)
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

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// args[0] is applicationNo, args[1] is applicationStatus
	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	exchangeApplication := exchangeApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &exchangeApplication)
	exchangeApplication.ApplicationStatus = args[1]

	changeApplicationStatusAsBytes, _ = json.Marshal(exchangeApplication)

	// 承認ステータスを更新
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

	// 授受票のIR承認が完了したら、有高管理する
	if args[1] == "28" {

		// 対象のSystemNoの有高管理情報を取得
		stockManagementAsBytes, _ := APIstub.GetState(args[2])

		if stockManagementAsBytes == nil {
			// 有高管理を初めて行うシステムの場合
			var stockManagement = stockManagement{
				SystemNo:      args[2],
				BorrowAmount:  args[3],
				DeleteAmount:  "0",
				StockAmount:   args[3],
				OverDueAmount: "0",
			}
			// 有高管理対象のシステムを追加
			stockManagementAsBytes, _ := json.Marshal(stockManagement)
			APIstub.PutState(args[2], stockManagementAsBytes)
		} else {
			// 既に有高管理をしているシステムがある場合
			stockManagement := stockManagement{}
			json.Unmarshal(stockManagementAsBytes, &stockManagement)

			// 現在のBorrowAmountに借用データ数をプラスする
			var borrowAmountInt int
			var borrowInt int
			borrowAmountInt, _ = strconv.Atoi(stockManagement.BorrowAmount)
			borrowInt, _ = strconv.Atoi(args[3])
			borrowAmountInt = borrowAmountInt + borrowInt
			stockManagement.BorrowAmount = strconv.Itoa(borrowAmountInt)

			// 現在のStockAmountに授受分をプラスする
			var stockAmountInt int
			stockAmountInt, _ = strconv.Atoi(stockManagement.StockAmount)
			stockAmountInt = stockAmountInt + borrowInt
			stockManagement.StockAmount = strconv.Itoa(stockAmountInt)

			// 有高情報を更新
			stockManagementAsBytes, _ := json.Marshal(stockManagement)
			APIstub.PutState(args[2], stockManagementAsBytes)
		}
	}
	return shim.Success(nil)
}

// 授受票のIR承認ルート設定
func (s *SmartContract) setExchangeRoute(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	exchangeApplication := exchangeApplication{}

	// IR承認ルートをセット
	json.Unmarshal(changeApplicationStatusAsBytes, &exchangeApplication)
	exchangeApplication.IRRegistrant = args[1]
	exchangeApplication.IRRechecker = args[2]
	exchangeApplication.IRReviewer = args[3]
	exchangeApplication.IRAuthorizer = args[4]
	// 授受票のIRの承認ステータスは25～28にする
	exchangeApplication.ApplicationStatus = "25"

	changeApplicationStatusAsBytes, _ = json.Marshal(exchangeApplication)
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

	return shim.Success(nil)
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

	// applicationStatusが38(BK承認完了)となった場合、
	// stockManagementを変更
	if args[1] == "38" {

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

func (s *SmartContract) setDelConfirmRoute(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	changeApplicationStatusAsBytes, _ := APIstub.GetState(args[0])
	delConfirmApplication := delConfirmApplication{}

	json.Unmarshal(changeApplicationStatusAsBytes, &delConfirmApplication)
	delConfirmApplication.BKRegistrant = args[1]
	delConfirmApplication.BKRechecker = args[2]
	delConfirmApplication.BKReviewer = args[3]
	delConfirmApplication.BKAuthorizer = args[4]
	delConfirmApplication.ApplicationStatus = "35"

	changeApplicationStatusAsBytes, _ = json.Marshal(delConfirmApplication)
	APIstub.PutState(args[0], changeApplicationStatusAsBytes)

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
