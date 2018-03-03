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
type delConfirm struct {
        ApplicationNo     string `json:"applicationNo"`
	delConfirmNo      string `json:"delConfirmNo"`
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
	DataNum           string `json:"DataNum"`
	DelDate           string `json:"DelDate"`
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
        if function == "createApplication" {
                return s.createApplication(APIstub, args)
        } else if function == "queryApplication" {
                return s.queryApplication(APIstub, args)
        } else if function == "queryAllApplications" {
                return s.queryAllApplications(APIstub)
        }
        return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

        if len(args) != 18 {
                return shim.Error("Incorrect number of arguments. Expecting 18")
        }

        //get maxNumber of application
        maxNumberAsBytes, _ := APIstub.GetState("maxApplicationNo")

        maxNumber := MaxNumber{}
        json.Unmarshal(maxNumberAsBytes, &maxNumber)

        var tmpNumber int
        tmpNumber, _ = strconv.Atoi(maxNumber.MaxApplicationNo)
        tmpNumber = tmpNumber + 1

        maxNumber.MaxApplicationNo = strconv.Itoa(tmpNumber)

        var delConfirm = delConfirm{
                ApplicationNo:     maxNumber.MaxApplicationNo,

                delConfirmNo:       args[0],
                Type:   "delConfirm",
                SystemNo:           args[1],
                ProjectName:        args[2],
                DataQuantity:       args[3],
                ReportingDate:      args[4],
                StartDate:          args[5],
                EndDate:            args[6],
                ApplicationStatus:  "30",
                IRRegistrant:       args[7],
                IRRechecker:        args[8],
                IRReviewer:         args[9],
                IRAuthorizer:       args[10],
                BKRegistrant:       args[11],
                BKRechecker:        args[12],
                BKReviewer:         args[13],
                BKAuthorizer:       args[14],
                EvidenceValue:      args[15],
                DataNum:            args[16],
                DelDate:            args[17],
		
        }

        maxNumberAsBytes, _ = json.Marshal(maxNumber)
        delConfirmAsBytes, _ := json.Marshal(delConfirm)

        APIstub.PutState("maxApplicationNo", maxNumberAsBytes)
        APIstub.PutState(maxNumber.MaxApplicationNo, delConfirmAsBytes)

        return shim.Success(nil)
}


func (s *SmartContract) queryApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        delConfirmAsBytes, _ := APIstub.GetState(args[0])
        return shim.Success(delConfirmAsBytes)
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