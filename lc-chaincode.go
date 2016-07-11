/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var lcIndexStr = "lcIndex"

type LC struct {
	Id           string `json:"cusip"`
	Name         string `json:"name"`
	ContractType string `json:"conracttype"`
	Vendor       string `json:"vendor"`
	Price        string `json:"price"`
	Bank         string `json:"bank"`
	Date         string `json:"date"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize index list
	err := stub.PutState(lcIndexStr, []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "create" {
		return t.create(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) create(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	fmt.Println("start creating a l/c")

	var lc LC
	var err error

	err = json.Unmarshal([]byte(args[0]), &lc)
	if err != nil {
		fmt.Println("error invalid L/C")
		return nil, errors.New("Error Unmarshal L/C")
	}

	lcBytes, err := json.Marshal(&lc)
	if err != nil {
		return nil, errors.New("Error Marshal L/C")
	}

	err = stub.PutState(lc.Id, lcBytes)

	if err != nil {
		return nil, errors.New("Error PutState L/C")
	}
	/*
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New("2nd argument must be a numeric string")
		}

		str := "{id: '" + args[0] + "',amount: '" + strconv.Itoa(amount) + "'}"
		err = stub.PutState(args[0], []byte(str))
		if err != nil {
			return nil, err
		}

		// lc index
		lcsAsBytes, err := stub.GetState(lcIndexStr)
	*/

	lcsAsBytes, err := stub.GetState(lcIndexStr)

	if err != nil {
		return nil, errors.New("Failed to get index")
	}
	var lcIndex []string
	json.Unmarshal(lcsAsBytes, &lcIndex)

	// append
	lcIndex = append(lcIndex, lc.Id)
	fmt.Println("! lc index: ", lcIndex)
	jsonAsBytes, err := json.Marshal(lcIndex)
	// store id of lc
	err = stub.PutState(lcIndexStr, jsonAsBytes)
	if err != nil {
		return nil, errors.New("Failed to put index")
	}

	fmt.Println("end creating marble")

	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
