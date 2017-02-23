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
	
	 return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		
		return t.Create_milkcontainer(stub, args)
		
	}else if function == "Create_coin" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coin(stub, args)	
        }else if function == "Order_milk"{                      // To order something - invoked by market - params - litres
		res,err :=  t.Order_milk(stub,args)
		jsonresp,_ := View_order(stub,args)
		fmt.Println(jsonresp)
		
		jsonresp,_ = init_logistics(stub,args)
		
		jsonresp,_ = set_user(stub,args)
		
		jsonresp,_ = checktheproduct(stub,args)
		
		jsonresp,_ = cointransfer(stub,args)
		
               
		printdetails(stub, 3)
		return res,err
	}
	fmt.Println("invoke did not find func: " + function)

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
	 

	return nil,nil

}





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



func printdetails(stub  shim.ChaincodeStubInterface, a int)(err error) {

	
	fmt.Println(a)
	err = stub.PutState("hello cts",[]byte("inside printdetails")) //Just to check the network whether we can read and write
       if err != nil {
		return  err
       }
	return err

}



func  View_order(stub shim.ChaincodeStubInterface, args []string) ([]string, error) {
	// This will be invoked by Supplier- think of UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty-but just for syntax-pass something as parameter
        a := args[0]
	fmt.Println(a)
	
	/* fetching the Orders*/
	
	
	
	ordersAsBytes, _ := stub.GetState(openOrdersStr)
	
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	
	
/*fetching the containers*/	
	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string             //an array to clone container indices
	json.Unmarshal(containerAsBytes, &containerIndex)
	
// From the list of Id's , picking up one Id and fetching its details
	
	containerAsBytes,_ = stub.GetState(containerIndex[0])
	
	res := MilkContainer{} 
        json.Unmarshal(containerAsBytes, &res)

// If ordered quantity and container quantity , then proceed and trigger logistics(How ever this is not automated here,we will do it 
	
	if (res.Litres == orders.OpenOrders[0].Litres) {
		fmt.Println("Found a suitable container")
		
		orders.OpenOrders[0].Status = "Ready to be Shipped"
		//t.init_logistics(stub,orders.OpenOrders[0].OrderId, containerIndex[0])
		ordersAsBytes,_ = json.Marshal(orders)
		stub.PutState("inside view order",[]byte("Hope this works"))
		stub.PutState(openOrdersStr,ordersAsBytes)
		
		OrderID := orders.OpenOrders[0].OrderID
		orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ShipOrder.Status = "Ready to be Shipped"
	 
	orderAsBytes,err = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,orderAsBytes)
		
		//t.read(stub,openOrdersStr)
	}else{
                stub.PutState("sorry",[]byte("we couldn't find a product for your choice of requirements"))
        }


	
	
	return nil,nil	
}

func init_logistics(stub shim.ChaincodeStubInterface, args []string) ([]string, error) {
	
	
	
	//args[0] args[1]
	// OrderId, ContainerID
	
	fmt.Println("Inside Init logistics function")
	OrderID := "abcd"
	//ContainerID := args[1]
	
	// fetch the order details and update status as "in transit"
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ShipOrder.Status = "In transit"
	 
	orderAsBytes,err = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,orderAsBytes)
	
	
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	orders.OpenOrders[0].Status = "In transit"
	ordersAsBytes,_ = json.Marshal(orders)
	stub.PutState(openOrdersStr,ordersAsBytes)
	
	
	//t.read(stub, OrderID)
	//t.set_user(stub,OrderID,ContainerID)
	
	
	
return nil,nil
}

func  set_user(stub shim.ChaincodeStubInterface, args []string) ([]string, error) {
	
// OrderId  ContainerID
//args[0] args[1]
	
//So here we will set the user name in container ID to the one in Order ID and Status to Delivered - Asset Transfer
	
	OrderID := "abcd"
	ContainerID := "1x223"
//fetch order details
       orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details	
	assetAsBytes,err := stub.GetState(ContainerID)
	container := MilkContainer{}
	json.Unmarshal(assetAsBytes, &container)

	if (container.User == "Supplier"){
	
	container.User = ShipOrder.User             //ASSET TRANSFER
	
       
	assetAsBytes,err = json.Marshal(container)
	stub.PutState(ContainerID, assetAsBytes)    //Pushing the updated container  back to the ledger
	
	ShipOrder.Status = "Delivered"
	orderAsBytes,err = json.Marshal(ShipOrder) 
	
	stub.PutState(OrderID,orderAsBytes)      //pushing the updated Order back to ledger
		
	//Updating the orders list 
		ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
		
		orders.OpenOrders[0].Status = ShipOrder.Status
		ordersAsBytes,_ = json.Marshal(orders)
		stub.PutState(openOrdersStr,ordersAsBytes)
	}else
        {
                stub.PutState("setuser",[]byte("failure in this function"))
                //t.read(stub,"setuser")
                return nil,nil
        }


return nil,nil
	//t.checktheproduct(stub,OrderID,ContainerID)
}


func  checktheproduct(stub shim.ChaincodeStubInterface, args []string) ([]string, error) {

// args[0] args[1]
// OrderID, ContainerID
	
	OrderID := "abcd"
	ContainerID := "1x223"
//fetch order details
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details
       assetAsBytes,_ := stub.GetState(ContainerID)
	Deliveredcontainer := MilkContainer{}
	json.Unmarshal(assetAsBytes, &Deliveredcontainer)

//check and transfer coin
	if (Deliveredcontainer.User == "Market" && Deliveredcontainer.Litres == ShipOrder.Litres) {
		
		fmt.Println("Thanks, I got the product")
		stub.PutState("Market Response",[]byte("Product received"))
		//t.cointransfer(stub,coinid) coinid -hard code it and send the coin id created by market
		return nil,nil
       }else{
                stub.PutState("checktheproduct",[]byte("failure"))
               // t.read(stub,"checktheproduct")
                return nil,nil
        }

	
return nil,nil


}


func cointransfer( stub shim.ChaincodeStubInterface, args []string) ([]string, error) {
	
//args[0] 
//coinID  
	//lets keep it simple for now, just fetch the coin from ledger, change username to Supplier and End of Story
	CoinID := "1x245"
	
	
	
	assetAsBytes,err := stub.GetState(CoinID)
	if err != nil{
		fmt.Println("Something wrog happened")
	}
	
	Transfercoin := SupplyCoin{}
	json.Unmarshal(assetAsBytes, &Transfercoin)
	
	if (Transfercoin.User == "Market") {   // check if the market guy actually holds coin in his name
	
		Transfercoin.User = "Supplier"
		assetAsBytes,err = json.Marshal(Transfercoin)
		stub.PutState(CoinID, assetAsBytes)
		return nil,nil
		
	}else{
	
		fmt.Println("There was some issue in transferring")
		//stub.PutState("cointransfer",[]byte("problem in coin transfer"))
		//t.read(stub,"cointransfer")
		return nil,nil
	}

	
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

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
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
