
package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	


var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)

var customerOrdersStr = "_customerorders"    // This will  be the key, value will be a list of orders placed by customer - wil be called by Customer

var supplierOrdersStr = "_supplierorders"     // this will be key, value will be a list of orders placed by supplier to logistics

type userandlitres struct{
	User string        `json:"user"`
	Litres int       `json:"litres"`
}

type MilkContainer struct{

        ContainerID string `json:"containerid"`
	Userlist  [2]userandlitres    `json:"userlist"`

}

type Order struct{
       OrderID string                  `json:"orderid"`
       User string                     `json:"user"`
       Status string                   `json:"status"`
       Litres int                      `json:"litres"`
}


type SupplierOrder struct {
   
        OrderID string                `json:"orderid"`
	Towhom string                 `json:"towhom"`
	ContainerID string            `json:"containerid"`
	
}


type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}


type AllSupplierOrders struct {
        SupplierOrdersList []SupplierOrder  `supplierOrdersList`
}
	

type Asset struct{
	User string        `json:"user"`
	ContainerIDs []string `json:"containerIDs"`
	LitresofMilk int `json:"litresofmilk"`
	Supplycoins int `json:"supplycoins"`
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	fmt.Printf("every time we enter main function")
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	
	var err error
	
	fmt.Println("Welcome  to  Supply chain management , Deployment has been started...")
	fmt.Printf("Hope for best, Plan for the worst")
	fmt.Printf("ready")
	
 
       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }
	
/* Resetting the container list - Making sure the value corresponding to openOrdersStr is empty */
	
       var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(containerIndexStr, jsonAsBytes)                     //Resetting - Making milk container list as empty 
       if err != nil {
		return nil, err
        }  
	
	
/* Resetting the customer and market order list  */
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	err = stub.PutState(customerOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	
/* Resetting the supplier order list  */
	var suporders AllSupplierOrders
	suporderAsBytes,_ := json.Marshal(suporders)
	err = stub.PutState(supplierOrdersStr, suporderAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
// Resetting the Assets of Supplier,Market, Logistics, Customer
	
	var emptyasset Asset
	
	emptyasset.User = "Supplier"
	jsonAsBytes, _ = json.Marshal(emptyasset)                // this is the byte format format of empty Asset structure
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // key -Supplier assets and value is empty now --> Supplier has no assets
	emptyasset.User = "Market"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("MarketAssets", jsonAsBytes)         // key -Market assets and value is empty now --> Market has no assets
	emptyasset.User = "Logistics"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("LogisticsAssets", jsonAsBytes)      // key - Logistics assets and value is empty now --> Logistic has no assets
	emptyasset.User = "Customer"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("CustomerAssets", jsonAsBytes)      // key - Customer assets and value is empty now --> Customer has no assets
	
	if err != nil {       
		return nil, err
}
	fmt.Println("Successfully deployed the code and orders and assets are reset")
	fmt.Printf("Go a head and play around")
	
return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "Create_coins" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coins(stub, args)	
        }else if function == "BuyMilkfrom_Retailer" { //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.BuyMilkfrom_Retailer(stub, args)	
        }else if function == "Vieworderby_Market" {  //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Vieworderby_Market(stub, args)	
        }else if function == "Checkstockby_Market" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Checkstockby_Market(stub, args)	
        }else if function == "Ordermilkto_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Ordermilkto_Supplier(stub, args)	
        }else if function == "Vieworderby_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Vieworderby_Supplier(stub, args)	
        }else if function == "Checkstockby_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Checkstockby_Supplier(stub,args)	
        }else if function == "Call_Logistics" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Call_Logistics(stub, args)	
        }else if function == "Vieworderby_Logistics"{
                return t.Vieworderby_Logistics(stub,args)
        }else if function == "pickuptheproduct" {
                return t.pickuptheproduct(stub,args)
        }else if function == "Deliverto_Market" {
                return t.Deliverto_Market(stub,args)
        }/*else if function == "Deliverto_customer" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Deliverto_customer(stub, args)	
        }*/
	
	fmt.Println("invoke did not find func: " + function)

return nil, errors.New("Received unknown function invocation: " + function)
}



func  Create_milkcontainer(stub shim.ChaincodeStubInterface, args [3]string) ( error) {
var err error

// "1x223" "supplier" "20" 
// args[0] args[1] args[2] 
	
	if len(args) != 3{
		return  errors.New("Please enter all the details")
        }
	fmt.Println("Hold on, we are Creating milkcontainer asset for you")
	
id := args[0]
user := args[1]
litres,err:=strconv.Atoi(args[2])
	if err != nil {
		return  errors.New("Litres argument must be a numeric string")
	}
	
// Checking if the container already exists in the network
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return  errors.New("Failed to get details of given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println("%+v\n",res)
        return errors.New("This container already exists")
}

//If not present, create it and Update ledger, containerIndexStr, Assets of Supplier
//Creation
        res.ContainerID = id
	res.Userlist[0].User=user
	res.Userlist[0].Litres = litres
	milkAsBytes, _ =json.Marshal(res)
        stub.PutState(res.ContainerID,milkAsBytes)
	fmt.Printf("Container created successfully, details are %+v\n", res)

//Update containerIndexStr	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return  errors.New("Failed to get container index")
	}
	var containerIndex []string                                        //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(containerAsBytes, &containerIndex)	
	
	
	containerIndex = append(containerIndex, res.ContainerID)          //append the newly created container to the global container list									//add marble name to index list
	fmt.Println("container indices in the network: ", containerIndex)
	jsonAsBytes, _ := json.Marshal(containerIndex)
        err = stub.PutState(containerIndexStr, jsonAsBytes)
	
// append the container ID to the existing assets of the Supplier
	
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	
	supplierasset.ContainerIDs = append(supplierasset.ContainerIDs, res.ContainerID)
	supplierasset.LitresofMilk += res.Userlist[0].Litres
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	fmt.Println("Balance of Supplier")
        fmt.Printf("%+v\n", supplierasset)
    //double checking
	supplierassetAsBytes,_ = stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	fmt.Printf("%+v\n", supplierasset)
	
	
	
	return nil

}


func (t *SimpleChaincode) Create_coins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"Market/Logistics/Customer",                  "100"
//args[0]                                     args[1]
//targeted owner                         No of supplycoins     
var err error
	user:= args[0]
	userAssets := user +"Assets"
        assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)

	asset.Supplycoins,err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New(" No of coins must be a numeric string")
	}
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	fmt.Println("Balance of " , user)
        fmt.Printf("%+v\n", asset)


return nil,nil
}



func (t *SimpleChaincode) BuyMilkfrom_Retailer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//args[0]      args[1]
//"cus123"       "10"
	var err error
	fmt.Println("Hello customer, welcome ")

	
	Openorder := Order{}
        Openorder.User = "customer"
        Openorder.Status = "Order received by Market"
        Openorder.OrderID = args[0]
        Openorder.Litres, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New(" No of coins must be a numeric string")
	}
	fmt.Println("Hello customer, your order has been generated successfully, you can track it with id in the following details")
	fmt.Println("%+v\n",Openorder)
        orderAsBytes,_ := json.Marshal(Openorder)
	stub.PutState(Openorder.OrderID,orderAsBytes)
	
	customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended",  Openorder.OrderID,"to existing customer orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(customerOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}

	return nil,nil
}

func(t *SimpleChaincode)  Vieworderby_Market(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET- think of UI-View orders- does he pass any parameter there...
// so here also no need of any arguments.
	
	fmt.Printf("Hello Market, these are the orders placed to  you by customer")
	
	
	ordersAsBytes, _ := stub.GetState(customerOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}



func (t *SimpleChaincode)  Checkstockby_Market(stub shim.ChaincodeStubInterface, args[]string) ([]byte, error){
	// In UI, beside each order one button to ship to customer, one button to check stock
	// we will extract details of orderId
	//we will exract asset balance of Market
	// if enough balance is der to deliver display "yes", if not der "no"
	//no tirggering is needed
	//OrderID should be passed in UI
//fetching order details
	OrderID := args[0]
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	quantity := ShipOrder.Litres
//fetching assets of market	
	marketassetAsBytes, _ := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset )
	
//checking if market has the stock	
	if (Marketasset.LitresofMilk >= quantity ){
		fmt.Println("Enough stock is available, Go ahead and deliver for customer")
		
//Call Deliver to customer function here
		b,_:= Deliverto_Customer(stub,ShipOrder.OrderID)
		fmt.Println(string(b))
		str := "Delivered to customer"
		return []byte(str), nil
		
	}else{
	        fmt.Println("Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer")
		str :=  "Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer"
	        ShipOrder.Status = "In transit to customer" // No matter, where the order placed by market is , for customer we will show it is "in transit"
	        orderAsBytes,err = json.Marshal(ShipOrder)
                stub.PutState(OrderID,orderAsBytes)  
		
		customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)	
	
	
		for i :=0; i<len(orders.OpenOrders);i++{
			if (orders.OpenOrders[i].OrderID == ShipOrder.OrderID){
			orders.OpenOrders[i].Status = "In transit to customer"
		         customerordersAsBytes , _ = json.Marshal(orders)
                        stub.PutState(customerOrdersStr,  customerordersAsBytes)
			}
	       }
	  return []byte(str), nil
		
		//Now we should send details of updated order status to customer, should be done in UI

		
        }
	
	return nil,nil

}


func Deliverto_Customer(stub shim.ChaincodeStubInterface ,args string) ([]byte,error){

	//args[0] 
	//OrderID  
	
	fmt.Println("Inside deliver to customer function")
//customer order
	OrderID := args
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	fmt.Println("%+v\n", ShipOrder)
	quantity := ShipOrder.Litres
	fmt.Println(quantity)
//market and customer assets
        marketassetAsBytes, _ := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset)
	fmt.Printf("%+v\n", Marketasset) 
	customerassetAsBytes, _ := stub.GetState("CustomerAssets")
	Customerasset := Asset{}             
	json.Unmarshal(customerassetAsBytes, &Customerasset)
	fmt.Printf("%+v\n", Customerasset) 
if (Marketasset.LitresofMilk >= quantity ){
	fmt.Println("Inside deliver to customer, market has quantity")
	
	id := Marketasset.ContainerIDs[0]
	
	
	milkAsBytes, err := stub.GetState(id) 
        if err != nil {
		return nil, errors.New("Failed to get details of given id") 
        }

        res := MilkContainer{} 
        json.Unmarshal(milkAsBytes, &res)
		
	fmt.Printf("%+v\n", res)
	
	
	
	
	
   // here we are assuming only one container is der and it has enough stock to provide
	if ( res.Userlist[0].Litres - quantity >0) {
		fmt.Println("yo yo..its about to complete")
                    
   //updating the container details, bcz it is shared now
		res.Userlist[0].Litres -= quantity // bringing down the market share of it
		res.Userlist[1].User = "Customer"
		res.Userlist[1].Litres = quantity
		fmt.Printf("%+v\n", res)
		milkAsBytes, _ =json.Marshal(res)
                stub.PutState(res.ContainerID,milkAsBytes)
		
  //updating customer assets
		
	              Customerasset.LitresofMilk += quantity
		if ( len(Customerasset.ContainerIDs) == 0){
		      fmt.Println("This is the first container of customer")
	   Customerasset.ContainerIDs = append(Customerasset.ContainerIDs ,id)
		}
		fmt.Printf("%+v\n", Customerasset)
			    Marketasset.LitresofMilk -= quantity
	
	              customerassetAsBytes,_ = json.Marshal(Customerasset)
	              stub.PutState("CustomerAssets",customerassetAsBytes)
	
	               marketassetAsBytes,_ = json.Marshal(Marketasset)
	               stub.PutState("MarketAssets",marketassetAsBytes)
	
	               ShipOrder.Status ="Delivered to Customer"
	               fmt.Printf("%+v\n", ShipOrder)
	               orderAsBytes,err = json.Marshal(ShipOrder)
                       stub.PutState(OrderID,orderAsBytes)
	
        customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)				
	
		for i :=0; i<len(orders.OpenOrders);i++{
			if (orders.OpenOrders[i].OrderID == ShipOrder.OrderID){
			orders.OpenOrders[i].Status = "Delivered to customer"
		         customerordersAsBytes , _ = json.Marshal(orders)
                        stub.PutState(customerOrdersStr,  customerordersAsBytes)
			}
	       }
	  
		b := [3]string{"30", "Customer", "Market"}
	           transfer(stub,b)        //Transfer should be automated. So it can't be invoked from UI..Loop hole
	               fmt.Println("FINALLLLLYYYY, END OF THE STORY")
         
                      return nil,nil
	}else{
	       return nil, errors.New("On a whole market has quantity, but it is divided into container, right now we are not going to that level")
	}
}else{
         return nil, errors.New(" No stock, give order to supplier")
 }

}


func(t *SimpleChaincode) Ordermilkto_Supplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
// "cus123"           "abcd"   
// CustomerOrderID    MarketOrderID 

var err error

//fetching the customer order details and ordering 5 times to the litres customer asked
CustomerOrderID := args[0]
orderAsBytes, err := stub.GetState(CustomerOrderID)
	if err != nil {
		return  nil, errors.New("Failed to get details of customer order, please make sure your id is correct")
	}
CustomerOrder := Order{} 
json.Unmarshal(orderAsBytes, &CustomerOrder)
quantity := CustomerOrder.Litres
	
//Generating market order

Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "Order placed to Supplier "
Openorder.OrderID = args[1]
Openorder.Litres = 5 * quantity

orderAsBytes,_ = json.Marshal(Openorder)
stub.PutState(Openorder.OrderID,orderAsBytes)
fmt.Println("your Order has been generated successfully")
fmt.Printf("%+v\n", Openorder)
	
//Add the new market order to market orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get  existing list of Market orders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)				
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended ",Openorder.OrderID,"to existing market orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
        }
	
	
return nil,nil
}



func(t *SimpleChaincode)  Vieworderby_Supplier(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET- think of UI-View orders- does he pass any parameter there...
// so here also no need of any arguments.
	
	fmt.Printf("Hello Supplier, these are the orders placed to  you by Market")
	
	
	ordersAsBytes, _ := stub.GetState(openOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}




func(t *SimpleChaincode)  Checkstockby_Supplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
/***FUNCTIONALITY EXPLAINED*******/
// In UI, beside each order one button to call logistics, one button to check stock
// we will extract details of orderId
//we will exract asset balance of Market
// if enough balance is der --> find a container and show it, if not create a new container (automated) and show it
//At the end of this function we will end up with a container
/*******/
	

//OrderID should be passed in UI
//fetching order details
//Market OrderID
//args[0]
	OrderID := args[0]
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	quantity := ShipOrder.Litres
//fetching assets of market	
	supplierassetAsBytes, _ := stub.GetState("SupplierAssets")
	supplierasset := Asset{}             
	json.Unmarshal(supplierassetAsBytes, &supplierasset )
	fmt.Printf("%+v\n", supplierasset)
//checking if Supplier has the stock	
if (supplierasset.LitresofMilk >= quantity ){
		fmt.Println("Enough stock is available, finding a suitable container.....")
		//length := len(Supplierasset.containerIDs)
 		/*for i:0 ; i<length; i++{
		
        // fetching the container details one by one
		       containerassetAsBytes, err := stub.GetState(supplierasset.containerIDs[i])
		       res := MilkContainer{} 
		       json.Unmarshal(containerassetAsBytes,&res)
        // Checking if the present container in loop has the quantity of Market asked
	              if (res.Userlist[0].Litres == ShipOrder.Litres){
		              fmt.Println("Found a suitable container, details are ")
			      fmt.Printf("%+v\n", res)
			      //Updating the status of market order
		              ShipOrder.Status = "Ready to be shipped to market"   // send the updated status to market
	                      orderAsBytes,err = json.Marshal(ShipOrder)
                              stub.PutState(OrderID,orderAsBytes)
		              return nil,nil
	              }
	  }
*/
	 fmt.Printf("%+v\n", supplierasset)
	cid := supplierasset.ContainerIDs[0]
	containerassetAsBytes, _ := stub.GetState(cid)
	res := MilkContainer{} 
	json.Unmarshal(containerassetAsBytes,&res)
        // Checking if the present container in loop has the quantity of Market asked
	/*              
	if (res.Userlist[0].Litres == ShipOrder.Litres){
	
	fmt.Println("Found a suitable container, below is the ID of the container, use it while placing order to Logistics")
	fmt.Println(Supplierasset.containerIDs[0])
	}
	*/
	fmt.Println("Found a suitable container, below is the ID of the container, use it while placing order to Logistics")
	fmt.Printf("%+v\n", res)
	   // return nil, errors.New("Supplier has the quantity but not all in one container, this will be covered in next phase")
}else{
	        fmt.Println("Right now there isn't sufficient quantity , Create a new container")
		var b [3]string
		b[0] = "1x223"
		b[1] = "Supplier"
		b[2] = strconv.Itoa(ShipOrder.Litres)
	 Create_milkcontainer(stub,b)

		
	       // fmt.Println("Successfully created container, check stock again to know your container details ") 
	        // can't call function again..loop hole
		//return nil,nil
}
	return nil,nil
}

func (t *SimpleChaincode)  Call_Logistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
//args[0]   //ToWhom  //Container ID
//OrderID   //Market   //"1x223"
	
// I think its fair only, in practical case, we will tell adrress for a postman to deliver, same thing here also
//Here Postman is Logistics guy, Receiver is market, letter is Container
	
	ShipOrder := SupplierOrder{}
	ShipOrder.OrderID = args[0]
	ShipOrder.Towhom = args[1]
	ShipOrder.ContainerID = args[2]
	
	orderAsBytes, _ :=json.Marshal(ShipOrder)
	stub.PutState( ShipOrder.OrderID, orderAsBytes)
	
	fmt.Println("Successfully placed order to Logistics")
	fmt.Println("%+v\n", ShipOrder)
	
	
	//Add the new Supplier order to market orders list
	ordersAsBytes, err := stub.GetState(supplierOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get  existing list of  orders placed by Supplier to logistics")
	}
	var suporders AllSupplierOrders
	json.Unmarshal(ordersAsBytes, &suporders)				
	suporders.SupplierOrdersList  = append(suporders.SupplierOrdersList, ShipOrder);		//append the new order - Openorder
	fmt.Println(" appended ",ShipOrder.OrderID,"to existing orders placed by Supplier to logistics")
	jsonAsBytes, _ := json.Marshal(suporders)
	err = stub.PutState(supplierOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
        }
	
	
	return nil,nil

}



func(t *SimpleChaincode) Vieworderby_Logistics(stub shim.ChaincodeStubInterface, args []string) ( []byte , error) {
	
	// This will be invoked by Supplier in UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty-but just for syntax-pass something as parameter in angular js
      
	
//fetching the Orders
	fmt.Printf("Hello Logistics, here are the orders placed to you by Supplier")
	fmt.Printf("Go a head and do your business")

	
	ordersAsBytes, _ := stub.GetState(supplierOrdersStr)
	
	var orders AllSupplierOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	fmt.Println(orders)
	return nil,nil
}



func(t *SimpleChaincode) pickuptheproduct(stub shim.ChaincodeStubInterface, args []string) ( []byte , error) {

// So in view order, he will see his orders, clicking on the order will show to whom and which container
//There will be a button "pickuptheproduct" which is equivalent to real life pick up --status will be in transit
//There will be one more button there only "Delivertheproduct"
//As of march 3, lets pass market order Id only as argument
//How can we update the order placed by market...without a notification
//here we are passing market order id only
	//args[0] args[1]
	// MarketOrderID, ContainerID
	MarketOrderID := args[0]
	
	
	// fetch the order details and update status as "in transit"
	orderAsBytes, err := stub.GetState(MarketOrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	MarketOrder := Order{} 
	json.Unmarshal(orderAsBytes, &MarketOrder)
	
	MarketOrder.Status = "In transit to market"
	 
	orderAsBytes,err = json.Marshal(MarketOrder)
	
	stub.PutState(MarketOrderID,orderAsBytes)
	
	fmt.Printf("%+v\n", MarketOrder)
	
	fmt.Println("Container is in transit")
	
return nil,nil
}



func(t *SimpleChaincode)  Deliverto_Market(stub shim.ChaincodeStubInterface, args []string) ([]byte , error) {
	
// SupplierOrderID      //MarketOrderID
//args[0]               //args[1]
	
//So here we will set the user name in container ID to the one in Order ID and Status to Delivered - Asset Transfer
// Why should logistics guy check if the supplier actually holds the container?????????
	fmt.Println("Delivering the container to Market")
	OrderID := args[0]
	
//fetch order details
       orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	ShipOrder := SupplierOrder{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ContainerID := ShipOrder.ContainerID
//fetch container details	
	assetAsBytes,err := stub.GetState(ContainerID)
	container := MilkContainer{}
	json.Unmarshal(assetAsBytes, &container)

	if (container.Userlist[0].User == "Supplier"){
	
	container.Userlist[0].User = "Market"         //ASSET TRANSFER
	fmt.Println("%+v\n", container)
	fmt.Println("pushing the updated container back to ledger")
	assetAsBytes,err = json.Marshal(container)
	stub.PutState(ContainerID, assetAsBytes)    //Pushing the updated container  back to the ledger
	
//fetch supplier assets
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
//fetch market assets
	userAssets := "MarketAssets"
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
//update market assets
	fmt.Println("Updating ",userAssets)
	asset.LitresofMilk += container.Userlist[0].Litres
	fmt.Println("appending", ContainerID,"to Market container id list")
        asset.ContainerIDs = append(asset.ContainerIDs,ContainerID)
       fmt.Printf("%+v\n", asset)
//update supplierassets
	
	fmt.Println("Updating Supplier assets..")
	supplierasset.LitresofMilk -= container.Userlist[0].Litres
	
	//WRITE A CODE  to remove that container id from supplier id list
		
		for i := 0 ;i < len(supplierasset.ContainerIDs);i++{
	
            if(supplierasset.ContainerIDs[i] == ContainerID){

            supplierasset.ContainerIDs =  append(supplierasset.ContainerIDs[:i],supplierasset.ContainerIDs[i+1:]...)
           break
       }	
}
	fmt.Printf("%+v\n", supplierasset)
	
//pushing updated ledger back to ledger
        supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
		
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
		
//double check
	supplierassetAsBytes,_ = stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
        fmt.Printf("%+v\n", supplierasset)
		
	assetAsBytes,_ = stub.GetState(userAssets)        // The same key which we used in Init function 
	json.Unmarshal( assetAsBytes, &asset)
	 fmt.Printf("%+v\n", asset)
//update the MarketOrder and push back to ledger
		
	MarketOrderID := args[1]
        orderAsBytes, err = stub.GetState(MarketOrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	MarketOrder := Order{} 
	json.Unmarshal(orderAsBytes, &MarketOrder)
	MarketOrder.Status = "Delivered to market"
	orderAsBytes,err = json.Marshal(MarketOrder) 
	stub.PutState(MarketOrderID,orderAsBytes)      
	fmt.Printf("%+v\n", ShipOrder)
		

	
	var b [2]string
	b[0] = args[1]
	b[1] = ContainerID

	checktheproduct(stub,b)
		
	}else
        {
                stub.PutState("delivertomarket",[]byte("failure in this function"))
                //t.read(stub,"setuser")
                return nil,nil
        }


return nil,nil
}



func  checktheproduct(stub shim.ChaincodeStubInterface, args [2]string) ( error) {

// args[0] args[1]
// MarketOrderID, ContainerID
	fmt.Println("Let us check the product")
	OrderID := args[0]
	ContainerID := args[1]
//fetch order details
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details
       assetAsBytes,_ := stub.GetState(ContainerID)
	Deliveredcontainer := MilkContainer{}
	json.Unmarshal(assetAsBytes, &Deliveredcontainer)

//check and transfer coins
	if (Deliveredcontainer.Userlist[0].User == "Market" && Deliveredcontainer.Userlist[0].Litres == ShipOrder.Litres) {
		
		fmt.Println("Thanks, I got  the right product, transferring amount to Supplier/Manufacturer")
		stub.PutState("Market Response",[]byte("Product received"))
		var b [3]string
		b[0]= "50"
		b[1] = "Market"
		b[2] = "Supplier"
		
		err = transfer(stub,b)
		if err!=nil{
			return err
		}
	        b[0]= "25"
		b[1] = "Supplier"
		b[2] = "Logistics"
		err = transfer(stub,b)
		if err!=nil{
			return err
		}
		return nil
       }else{
                stub.PutState("checktheproduct",[]byte("failure"))
		fmt.Println("I didn't get the right product")
              
                return nil
        }

	
return nil


}



func transfer( stub shim.ChaincodeStubInterface, args [3]string) ( error) {
	
//args[0]             args[1]         args[2]
//No of supplycoin      Sender         Reciever   
	//lets keep it simple for now, just fetch the coin from ledger, change username to Supplier and End of Story
	transferamount,_ := strconv.Atoi(args[0])
	sender := args[1]                               // this thing should be given by us in UI background
	receiver := args[2]                            // this will be given by the user on web page
	
	fmt.Println( sender, "transferring", transferamount, "coins to", receiver)
	
        senderAssets := sender +"Assets"
        senderassetAsBytes,_ := stub.GetState(senderAssets)        // The same key which we used in Init function 
	senderasset := Asset{}
	json.Unmarshal( senderassetAsBytes, &senderasset)
	
	
	receiverAssets := receiver+"Assets"
        receiverassetAsBytes,_ := stub.GetState(receiverAssets)        // The same key which we used in Init function 
	receiverasset := Asset{}
	json.Unmarshal( receiverassetAsBytes, &receiverasset)
	
	if ( senderasset.Supplycoins >= transferamount){
		
	senderasset.Supplycoins -= transferamount
	receiverasset.Supplycoins += transferamount
	
        senderassetAsBytes,_=  json.Marshal(senderasset)
	stub.PutState(senderAssets,  senderassetAsBytes)
	fmt.Println("Balance of " , sender)
       fmt.Printf("%+v\n", senderasset)
		
	receiverassetAsBytes,_=  json.Marshal(receiverasset)
	stub.PutState( receiverAssets,receiverassetAsBytes)
	fmt.Println("Balance of " , receiver)
        fmt.Printf("%+v\n", receiverasset)
		return  nil
	}else {
		str := "Failed to transfer amount from" + sender + "to" + receiver
		return  errors.New(str)
	}
	

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
