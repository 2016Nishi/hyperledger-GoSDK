package peernetwork

import (
	"fmt"
	//"strconv"
	"errors"
	//"log"
	"strings"
)

/*
  prints the content of the network: peers, users, and chaincodes.
 */
func PrintNetworkDetails() {

	ThisNetwork := LoadNetwork()
	Peers := ThisNetwork.Peers
	i := 0
	for i < len(Peers) {

		msgStr := fmt.Sprintf("ip: %s port: %s name %s ", Peers[i].PeerDetails["ip"], Peers[i].PeerDetails["port"], Peers[i].PeerDetails["name"])
		fmt.Println(msgStr)
		userList := ThisNetwork.Peers[i].UserData
		fmt.Println("Users:")
		for user, secret := range userList {

			fmt.Println(user, secret)
		}
		i++
	}
	fmt.Println("Available Chaincodes :")
	libChainCodes := InitializeChainCodes()
	for k, v := range libChainCodes.ChainCodes {
		fmt.Println("\nChaincode :", k)
	  fmt.Println("\nDetail :\n")
		for i, j := range v.Detail {
			msgStr := fmt.Sprintf("user: %s secret: %s", i, j)
			fmt.Println(msgStr)
		}
		fmt.Println("\n")
	}
}

/*
 Gets ChainCode detail for a given chaincode name
  Takes two arguments
	1. name (string)			- name of the chaincode as specified in CC_Collections.json file
	2. lcc (LibChainCodes)		- LibChainCodes struct having current collection of all chaincodes loaded in the network.
  Returns:
 	1. ccDetail map[string]string  	- chaincode details of the chaincode requested as a map of key/value pairs.
	2. Versions map[string]string   - versioning or tagging details on the chaincode requested as a map of key/value pairs
 */
func GetCCDetailByName(name string, lcc LibChainCodes) (ccDetail map[string]string, versions map[string]string, err error) {
	var errStr string
	var err1 error
	for k, v := range lcc.ChainCodes {
		if strings.Contains(k, name) {
			return v.Detail, v.Versions, err1
		}
	}
	//no more chaincodes construct error string and empty maps
	errStr = fmt.Sprintf("chaincode %s does not exist on the network", name)
	//need to check for this
	j := make(map[string]string)
	return j, j, errors.New(errStr)
}


/** utility functions to aid users in getting to a valid URL on network
 ** to post chaincode rest API
 **/

/*
  gets any one peer from 'thisNetwork' as set by chaincode.Init()
 */
func APeer(thisNetwork PeerNetwork) *Peer {
	//thisNetwork := LoadNetwork()
	Peers := thisNetwork.Peers
	var aPeer *Peer
	//get any peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			aPeer = &Peers[peerIter]
		}
	}
	//fmt.Println(" * aPeer ", *aPeer)
	//fmt.Println(" ip ", aPeer.PeerDetails["ip"])
	return (aPeer)
}

/*
  gets any one user from any Peer on the entire network.
 */
func AUserFromNetwork(thisNetwork PeerNetwork) (ip string, port string, user string) {

	//fmt.Println("Values inside AUserFromNetwork ", ip, port, user)
	var u string
	aPeer := APeer(thisNetwork)
	users := aPeer.UserData

	//fmt.Println(" ip ", aPeer.PeerDetails["ip"])
	for u, _ = range users {
		break
	}
	//fmt.Println(" ip ", aPeer.UserData["ip"])
	//fmt.Println(" ip ", user)
	return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], u
}

/*
  finds any one user associated with the given peer
*/
func AUserFromAPeer(thisPeer Peer) (ip string, port string, user string) {

	//var aPeer *Peer
	aPeer := thisPeer
	var curUser string
	userList := aPeer.UserData
	for curUser, _ = range userList {
		break
	}
	//fmt.Println(" ip ", aPeer.UserData["ip"])
	//fmt.Println(" ip ", user)
	return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], curUser
}

/*
 gets a particular user from a given Peer on the PeerNetwork
 */
func AUserFromThisPeer(thisNetwork PeerNetwork, host string) (ip string, port string, user string, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var u string
	var errStr string
	var err1 error
	//get a random peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if strings.Contains(Peers[peerIter].PeerDetails["ip"], host) {
				aPeer = &Peers[peerIter]
			}
		}
	}
	//fmt.Println(" * aPeer ", *aPeer)
	if aPeer != nil {
		users := aPeer.UserData
		for u, _ = range users {
			break
		}
		return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], u, err1
	} else {
		errStr= fmt.Sprintf("%s, Not found on network", host)
		return "", "", "", errors.New(errStr)
	}
}

/*
  finds the peer address corresponding to a given user
    thisNetwork as set by chaincode.init
	ip, port are the address of the peer.
	user is the user details: name, credential.
	err	is an error message, or nil if no error occurred.
 */
func PeerOfThisUser(thisNetwork PeerNetwork, username string) (ip string, port string, user string, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	var err1 error
	//fmt.Println("Inside function")
	//get a random peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if _, ok := Peers[peerIter].UserData[username]; ok {
				fmt.Println("Found %s in network", username)
				aPeer = &Peers[peerIter]
			}
		}
	}

	if aPeer == nil {
		errStr = fmt.Sprintf("%s, Not found on network", username)
		return "", "", "", errors.New(errStr)
	} else {
		return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], username, err1
	}
}

/********************
type PeerNetworks struct {
	PNetworks      []PeerNetwork
}


func AddAPeerNetwork() {

}

func DeleteAPeerNetwork() {

}

func AddUserOnAPeer(){

}

func RemoveUserOnAPeer(){

}


func LoadNetworkByName(name string) PeerNetwork {

  networks := LoadPeerNetworks()
	pnetworks := networks.PNetworks
	for peerIter := range pnetworks {
		//fmt.Println(pnetworks[peerIter].Name)
		if strings.Contains(pnetworks[peerIter].Name, name) {
			return pnetworks[peerIter]
		}
	}
	//return *new(PeerNetwork)
}
*********************/
