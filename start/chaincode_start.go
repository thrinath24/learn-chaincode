

package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	

var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)



type MilkContainer struct{

        ContainerID string `json:"containerid"`
        User string        `json:"user"`

        Litres string        `json:"litres"`

}



type SupplyCoin struct{

        CoinID string `json:"coinid"`
        User string        `json:"user"`
}

type Order struct{
        OrderID string `json:"orderid"`
       User string `json:"user"`
       Status string `json:"status"`
       Litres string    `json:"litres"`
}

type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}
/*
type Asset struct{
	  User string        `json:"user"`
	containerIDs []string `json:"containerids"`
	coinIds []string `json:"coinids"`
}

*/

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var err error

       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }
	
        /* Reset container index list - Making sure the value corresponding to containerIndexStr  is empty */

       var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(containerIndexStr, jsonAsBytes)                     //Resetting - Making milk container list as empty 
       if err != nil {
		return nil, err
        }  
	
	/* Resetting the order list - Making sure the value corresponding to openOrdersStr is empty */
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	// Resetting the Assets of Supplier for test case- later on we can do for market and logistics also
	/*
	var emptyasset Asset
	
	
	jsonAsBytes, _ = json.Marshal(emptyasset)
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // Supplier assets are empty now
	
	*/
        return nil, nil

}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		return t.Create_milkcontainer(stub, args)      
	}else if function == "Create_coin" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coin(stub, args)	
        }
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}



func (t *SimpleChaincode) Create_milkcontainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error

// "1x22" "supplier" 20 
// args[0] args[1] args[2] 

id := args[0]
user := args[1]
litres :=args[2] 
	
// Checking if the container already exists in the network
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details og given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

//If not present, create it and Update ledger, containerIndexStr, Assets of Supplier
//Creation
res.ContainerID = id
res.User = user
res.Litres = litres
milkAsBytes, _ =json.Marshal(res)

stub.PutState(res.ContainerID,milkAsBytes)
	
//Update containerIndexStr	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string                                        //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(containerAsBytes, &containerIndex)	
	
	
	containerIndex = append(containerIndex, res.ContainerID)          //append the newly created container to the global container list									//add marble name to index list
	fmt.Println("! container index: ", containerIndex)
	jsonAsBytes, _ := json.Marshal(containerIndex)
        err = stub.PutState(containerIndexStr, jsonAsBytes)

// append the container ID to the existing assets of the Supplier
	/*
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	supplierasset.containerIDs = append(supplierasset.containerIDs, res.ContainerID)
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
*/
	//t.read(stub,"SupplierAssets")
	
	
return nil,nil

}



/******************Asset Coin creation****************/

func (t *SimpleChaincode) Create_coin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"1x245" "Market/Logistics"
id := args[0]
user:= args[1]
//Check if coin already exists in network
coinAsBytes , err := stub.GetState(id)
if err != nil{
              return nil, errors.New("Failed to get details of given id")
} 

res :=SupplyCoin{}

json.Unmarshal(coinAsBytes, &res)

if res.CoinID == id{

          fmt.Println("Coin already exists")
          fmt.Println(res)
          return nil,errors.New("This coin already exists")
}
// Proceed to create if not der in ntwrk
res.CoinID = id
res.User = user

coinAsBytes, _ = json.Marshal(res)
stub.PutState(id,coinAsBytes)
//t.read(stub,"res.CoinID")
return nil,nil
}



/*****************ORDER MILK****************/


func (t *SimpleChaincode) Order_milk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//"20"
//litres
var err error
Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "pending"
Openorder.OrderID = "abcd"
Openorder.Litres = args[0]
orderAsBytes,_ := json.Marshal(Openorder)
	
err = stub.PutState(Openorder.OrderID,orderAsBytes)
	
if err != nil {
		return nil, err
}

//Add the new order to the orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println("! appended Openorder to orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}
	//t.read(stub,"openOrdersStr")
return nil,nil
}



















// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
