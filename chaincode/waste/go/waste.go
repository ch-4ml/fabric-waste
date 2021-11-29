// SmartContract(CC) of Safe Purging for Hamful Waste & Recycling from Factory 
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {}

type Waste struct {
	ObjectType string `json:"docType"`
	CasNo  string `json:"casNo"`         
	Name   string `json:"name"`			 
	Color  string `json:"color"`		 
	Phase  string `json:"phase"`         
	Quantity string `json:"quantity"`
}

type EmissionRecord struct {
	ObjectType string `json:"docType"`
	CasNo	string	`json:"casNo"`
	Cost string `json:"cost"`
	Quantity string `json:"quantity"`
}

func (t *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (t *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// 핸들러 이용한 
	if function == "initLedger" {
		return initLedger(APIstub)
	} else if function == "registerWaste" { // 물질 등록
		return registerWaste(APIstub, args)
	} else if function == "createWaste" { // 물질 생성
		return createWaste(APIstub, args)
	} else if function == "queryAllWastes" { // 조회
		return queryAllWastes(APIstub)
	} else if function == "purgeWaste" { // 배출
		return purgeWaste(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	ObjectType := "waste"
	
	wastes := []Waste{
		Waste{Phase: "Gas", Name: "Cyanogen chloride", Color: "None", CasNo: "506-77-4", Quantity: "5000"},
		Waste{Phase: "Liquid", Name: "m-Crezol(Phenol)", Color: "yellow", CasNo: "108-39-4", Quantity: "5000"},		
		Waste{Phase: "Oil", Name: "Cyanogen chloride", Color: "None", CasNo: "506-77-4", Quantity: "5000"},
		Waste{Phase: "Solid", Name: "Carbon Steel(Metal)", Color: "None", CasNo: "scrap001", Quantity: "5000"},				
	}

	i := 0
	for i < len(wastes) {
		fmt.Println("i is ", i)
		wasteAsBytes, _ := json.Marshal(wastes[i])
		APIstub.PutState(ObjectType + strconv.Itoa(i), wasteAsBytes)
		fmt.Println("Added", wastes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

// 공장에서 처리할 폐기물을 등록하는 함수
func registerWaste(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of argument Expecting 5")
	}

	ObjectType := "waste"
	waste := Waste{ObjectType: ObjectType, Phase: args[1], Name: args[2], Color: args[3], CasNo: args[4], Quantity: "0"}

	wasteAsBytes, _ := json.Marshal(waste)
	APIstub.PutState(ObjectType + waste.CasNo + "-" + waste.Phase, wasteAsBytes)

	return shim.Success(nil)
}

// 공장을 가동하면서 쌓이는 폐기물을 기록하는 함수
func createWaste(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// 0: CasNo, 1: Phase, 2: Quantity to add
	if len(args) != 3 { 
		return shim.Error("Incorrect number of argument Expecting 2")
	}

	ObjectType := "waste"

	casNo := args[0]
	phase := args[1]

	wasteAsBytes, _ := APIstub.GetState(ObjectType + casNo + "-" + phase)

	waste := Waste{}
	json.Unmarshal(wasteAsBytes, &waste)
	
	q1, _ := strconv.Atoi(waste.Quantity)
	q2, _ := strconv.Atoi(args[2])

	quantity := q1 + q2
	quantityToString := strconv.Itoa(quantity)

	waste.Quantity = quantityToString
	wasteAsBytes, _ = json.Marshal(waste)
	
	APIstub.PutState(ObjectType + casNo + "-" + phase, wasteAsBytes)

	return shim.Success(nil)
}

func queryAllWastes(APIstub shim.ChaincodeStubInterface) sc.Response {

	queryString := fmt.Sprintf("{\"selector\": {\"docType\": \"waste\"}}")
	resultsIterator, err := APIstub.GetQueryResult(queryString)
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

	fmt.Printf("- queryAllWastes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// 폐기물 처리
func purgeWaste(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


	return shim.Success(nil)
}

// 폐기물 처리 내역 조회


// The main function is only relevant in unit test mode. Only included here for completenes
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
