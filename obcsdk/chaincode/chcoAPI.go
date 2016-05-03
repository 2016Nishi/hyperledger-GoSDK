package chaincode

import (
	"fmt"
	"log"
	"obcsdk/peernetwork"
	"strings"
	"errors"
)

var ThisNetwork peernetwork.PeerNetwork
var Peers = ThisNetwork.Peers
var ChainCodeDetails, Versions map[string]string
var LibCC peernetwork.LibChainCodes

const (deployUsage =  ("dAPIArgs0 := []string{\"example02\", \"init\",  \"(optional)<tagName>\"}\n" +
											 "depArgs0 := []string{\"a\", \"20000\", \"b\", \"9000\"}\n"  +
											 "chaincode.Deploy(dAPIArgs0, depArgs0)\n")

				invokeUsage = ("iAPIArgs0 := []string{\"example02\", \"invoke\"}\n" +
												"invArgs0 := []string{\"a\", \"b\", \"500\"}\n" +
												"chaincode.Invoke(iAPIArgs0, invArgs0)}\n")
				invokeOnPeerUsage = ("iAPIArgs0 := []string{\"example02\", \"invoke\", \"<PEER_IP_ADDRESS>\" , \"(optional)<tagName>\"}\n" +
										 					"invArgs0 := []string{\"a\", \"b\", \"500\"} \n" +
										 					"chaincode.Invoke(iAPIArgs0, invArgs0)}\n")
				invokeAsUserUsage = ("\niAPIArgs0 := []string{\"example02\", \"invoke\", \"<Registered_USER_NAME>\" , \"(optional)<tagName>\"}\n" +
									 					"invArgs0 := []string{\"a\", \"b\", \"500\"} \n" +

									 					"chaincode.Invoke(iAPIArgs0, invArgs0)}\n")

				queryUsage = ("qAPIArgs0 := []string{\"example02\", \"query\" , \"(optional)<tagName>\"}\n" +
											"qArgsa := []string{\"a\"}\n" +
										"chaincode.Query(qAPIArgs0, qArgsa)\n")
		)

/*
  initializes users on network using data supplied in NetworkCredentials.json file
 */
func InitNetwork() {
	ThisNetwork = peernetwork.LoadNetwork()
}
/*
   initializes chaincodes on network using information supplied in CC_Collections.json file
 */
func InitChainCodes() {
	LibCC = peernetwork.InitializeChainCodes()
}
/*
  initializes network based on files in directory utils
 */
func Init() {

	InitNetwork()
	InitChainCodes()
}

/*
   Registers each user on the network based on the content of ThisNetwork.Peers.
 */
func RegisterUsers() {
	fmt.Println("\nCalling Register ")

	//testuser := peernetwork.AUser(ThisNetwork)
	Peers := ThisNetwork.Peers
	i := 0
	for i < len(Peers) {

		userList := ThisNetwork.Peers[i].UserData
		for user, secret := range userList {
			url := "http://" + Peers[i].PeerDetails["ip"] + ":" + Peers[i].PeerDetails["port"]
			msgStr := fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
			fmt.Println(msgStr)
			register(url, user, secret)
		}
		fmt.Println("Done Registering ", len(userList), "users on ", Peers[i].PeerDetails["name"])
		i++
	}
}

/*
   deploys a chaincode in the fabric to later execute functions on this deployed chaincode
   Takes two arguments
 	 A. args []string
	   	1.ccName (string)			- name of the chaincode as specified in CC_Collections.json file
		2.funcName (string)			- name of the function to call from chaincode specification
									"init" for chaincodeexample02
		3.tagName(string)(optional)		- tag a deployment to support something like versioning

 	B. depargs []string				- actual arguments passed to initialize chaincode inside the fabric.

		Sample Code:
		dAPIArgs0 := []string{"example02", "init"}
		depArgs0 := []string{"a", "20000", "b", "9000"}

		var depRes string
		var err error
		depRes, err := chaincode.Deploy(dAPIArgs0, depArgs0)
 */
func Deploy(args []string, depargs []string) error {

			 if (len(args) < 2)  || (len(args) > 3) {
				 			 fmt.Println(deployUsage)
							 return errors.New("Deploy : Incorrect number of arguments. Expecting 2 or 3")
			 }
			 ccName := args[0]
 			 funcName := args[1]
			 var tagName string
 			 if (len(args) == 2) {
 					tagName = ""
 				} else if  len(args) == 3 {
				 	tagName = args[2]
 				}
			 dargs := depargs
			 var err error

                         ChainCodeDetails, Versions, err = peernetwork.GetCCDetailByName(ccName, LibCC)
			 if err != nil {
				 fmt.Println("Inside deploy: ", err)
				 //log.Fatal("No Chain Code Details we cannot proceed")
				 return errors.New("No Chain Code Details we cannot proceed")
				 //exit(1)
			 }
			 if strings.Contains(ChainCodeDetails["deployed"], "true") {
				 fmt.Println("\n\n ** Already deployed ..")
				 fmt.Println(" skipping deploy...")
			 } else {
				 msgStr := fmt.Sprintf("\n** Initializing and deploying chaincode %s on network with args %s\n",ChainCodeDetails["path"], dargs)
				 fmt.Println(msgStr)
				 restCallName := "deploy"
				 ip, port, auser := peernetwork.AUserFromNetwork(ThisNetwork)
				 url := "http://" + ip + ":" + port
				 txId := changeState(url, ChainCodeDetails["path"], restCallName, dargs, auser, funcName)
				 //storing the value of most recently deployed chaincode inside chaincode details if no tagname or versioning
				 ChainCodeDetails["dep_txid"] = txId
				 if len(tagName) != 0 {
				 		Versions[tagName] = txId
					}
				 //fmt.Println(ChainCodeDetails["dep_txid"])
				 //ChainCodeDetails["deployed"] = "true"
			 }
			 return err
}
/*
 changes state of a chaincode by passing arguments to BlockChain REST API invoke.
 Takes two arguments
 	 A. args []string
	    1.ccName (string)			- name of the chaincode as specified in CC_Collections.json file
		2.funcName (string)		- name of the function to call from chaincode specification
								"invoke" for chaincodeexample02
		3.tagName(string)(optional)	- tag a deployment to support something like versioning

	B. invargs []string			- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func Invoke(args []string,  invokeargs []string) (id string, err error) {

	if (len(args) < 2)  || (len(args) > 3) {
					fmt.Println("Invoke : Incorrect number of arguments. Expecting 2 in invokeAPI arguments")
					fmt.Println(invokeUsage)
					return "", errors.New("Invoke : Incorrect number of arguments. Expecting 2 in invokeAPI arguments")
	}
	ccName := args[0]
	funcName := args[1]
	var tagName string
	if (len(args) == 2) {
		 tagName = ""
	 } else if  len(args) == 3 {
		 tagName = args[2]
	 }
	invargs := invokeargs
	//fmt.Println("Inside invoke .....")
	var err1 error
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside invoke: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
  restCallName := "invoke"
	aPeer := peernetwork.APeer(ThisNetwork)
	fmt.Println(aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"])
	ip, port, auser := peernetwork.AUserFromAPeer(*aPeer)
	url := "http://" + ip + ":" + port
	msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, invargs, url, auser)
	fmt.Println(msgStr0)
	var txId string
	if len(tagName) != 0 {

		txId = changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
	}else {
		txId = changeState(url, (ChainCodeDetails["dep_txid"]), restCallName, invargs, auser, funcName)
	}
	//fmt.Println("\n\n\n*** END Invoking as  ***\n\n", auser, " on a single peer\n\n")
	return txId, errors.New("")
}

/*
 changes state of a chaincode on a specific peer by passing arguments to REST API call
 Takes two arguments
	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"invoke" for chaincodeexample02
		3. host (string)				- hostname or ipaddress to call invoke from
		4. tagName(string)(optional)			- tag the invocation to support something like versioning

	B. invargs []string					- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke", "127.0.0.1"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func InvokeOnPeer(args []string, invokeargs []string) (id string, err error) {

	//fmt.Println("Inside InvokeOnPeer .....")
	if (len(args) < 3)  || (len(args) > 4) {
					fmt.Println("InvokeOnPeer : Incorrect number of arguments. Expecting 3 or 4 in invokeAPI arguments")
					fmt.Println(invokeOnPeerUsage)
					return "", errors.New("InvokeOnPeer : Incorrect number of arguments. Expecting 3 or 4 number in invokeAPI arguments")
	}
	ccName := args[0]
	funcName := args[1]
	host := args[2]
	var tagName string
	if (len(args) == 3) {
		 tagName = ""
	 } else if  len(args) == 4 {
		 tagName = args[3]
	 }
	invargs := invokeargs
	var err1 error
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside InvokeOnPeer: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}

  restCallName := "invoke"
	ip, port, auser, err2 := peernetwork.AUserFromThisPeer(ThisNetwork, host)
	if err2 != nil {
  	fmt.Println("Inside invoke3: ", err2)
		return "", err2
	} else {
		url := "http://" + ip + ":" + port
		msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s on %s\n", funcName, ccName, invargs, url, auser, host)
		fmt.Println(msgStr0)
		txId := changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
		return txId, errors.New("")
	}
}

/*
 changes state of a chaincode using a specific user credential
  Takes two arguments
 	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"invoke" for chaincodeexample02
		3. user (string)				- login name of a registered user
		4. tagName(string)(optional)			- tag the invocation to support something like versioning

	B. invargs []string					- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke", "jim"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func InvokeAsUser(args []string, invokeargs []string) (id string, err error) {
	if (len(args) < 3)  || (len(args) > 4) {
					fmt.Println("InvokeAsUser : Incorrect number of arguments. Expecting 3 or 4 in invokeAPI arguments")
					fmt.Println(invokeAsUserUsage)
					return "", errors.New("InvokeAsUser : Incorrect number of arguments. Expecting 3 or 4 number in invokeAPI arguments")
	}
	ccName := args[0]
	funcName := args[1]
	userName := args[2]
	var tagName string
	if (len(args) == 3) {
		 tagName = ""
	 } else if  len(args) == 4 {
		 tagName = args[3]
	 }
	invargs := invokeargs
	var err1 error
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside InvokeAsUser: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
  restCallName := "invoke"
	ip, port, auser, err2 := peernetwork.PeerOfThisUser(ThisNetwork, userName)
	if err2 != nil {
    fmt.Println("Inside InvokeAsUser: ", err2)
		return "", err2
	} else {
		url := "http://" + ip + ":" + port
		msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, invargs, url, auser)
		fmt.Println(msgStr0)
		txId := changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
		return txId, errors.New("")
	}
}

/*
  Query fetches the value of the arguments supplied to query function from the fabric.
  Takes two arguments
 	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"query" for chaincodeexample02
		3. tagName(string)(optional)	- tag the invocation to support something like versioning

	B. qargs []string					- actual arguments passed to get the values as stored inside fabric.

		Sample Code:
		qAPIArgs0 := []string{"example02", "query"}
		qArgsa := []string{"a"}

		var queryRes string
		var err error
		queryRes,err := chaincode.Query(qAPIArgs0, qArgsa)
*/
func Query(args []string, queryArgs []string) (id string, err error) {

	if (len(args) < 2)  || (len(args) > 3) {
					fmt.Println(queryUsage)
					return "", errors.New("Incorrect number of arguments. Expecting 2 in queryAPI arguments")
	}
	ccName := args[0]
	funcName := args[1]
	var tagName string
	if (len(args) == 2) {
		 tagName = ""
	 } else if  len(args) == 3 {
		 tagName = args[2]
	 }
	qargs := queryArgs
	var err1 error

	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
  if err1 != nil {
		fmt.Println("Inside Query: ", err1)
		fmt.Println("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
  restCallName := "query"
  ip, port, auser := peernetwork.AUserFromNetwork(ThisNetwork)
	url := "http://" + ip + ":" + port
	var txId string
	msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, queryArgs, url, auser)
	fmt.Println(msgStr0)

	if len(tagName) != 0 {
		txId = readState(url, Versions[tagName], restCallName, qargs, auser, funcName)
	}else {
		txId = readState(url, (ChainCodeDetails["dep_txid"]), restCallName, qargs, auser, funcName)
	}

	return txId, errors.New("")
} /* Query() */
