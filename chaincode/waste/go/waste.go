// SmartContract(CC) of Safe Purging for Hamful Waste & Recycling from Factory 

package wastetransfer

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
	Name   string `json:"name"`			 
	Color  string `json:"color"`		 
	CasNo  string `json:"casno"`         
	Phase  string `json:"phase"`         
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
	} else if function == "createWaste" {
		return createWaste(APIstub, args)
	} else if function == "queryAllWastes" {
		return queryAllWastes(APIstub)
	} 
	// else if function == "purgeWaste" {
	// 	return purgeWaste(APIstub, args)
	// }
	return shim.Error("Invalid Smart Contract function name.")
}

func initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	ObjectType := "waste"
	
	wastes := []Waste{
		Waste{Phase: "Gas", Name: "Cyanogen chloride", Color: "None", CasNo: "506-77-4"},
		Waste{Phase: "Liquid", Name: "m-Crezol(Phenol)", Color: "yellow", CasNo: "108-39-4"},		
		Waste{Phase: "Oil", Name: "Cyanogen chloride", Color: "None", CasNo: "506-77-4"},
		Waste{Phase: "Solid", Name: "Carbon Steel(Metal)", Color: "None", CasNo: "scrap001"},				
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

func createWaste(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of argument Expecting 5")
	}

	ObjectType := "waste"
	waste := Waste{ObjectType: ObjectType, Phase: args[1], Name: args[2], Color: args[3], CasNo: args[4]}

	wasteAsBytes, _ := json.Marshal(waste)
	APIstub.PutState(ObjectType + args[0], wasteAsBytes)

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

// func purgeWaste(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	return shim.Success(nil)
// }


// The main function is only relevant in unit test mode. Only included here for completenes
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
