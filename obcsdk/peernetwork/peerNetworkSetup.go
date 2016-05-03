package peernetwork

import (
	"fmt"
	"os"
	"strconv"
	//"errors"
	"log"
	//"strings"
)



type Peer struct {
	PeerDetails map[string]string
	UserData    map[string]string
}

type PeerNetwork struct {
	Peers []Peer
	Name string
}

type LibChainCodes struct {
	ChainCodes map[string]ChainCode
}

type ChainCode struct {
	Detail map[string]string
	Versions map[string]string
}

/*
  creates network as defined in NetworkCredentials.json, distributing users evenly among the peers of the network.
 */
func LoadNetwork() PeerNetwork {

	p, n := initializePeers()

	anetwork := PeerNetwork{Peers: p, Name: n}
	return anetwork
}

/*
  reads CC_Collection.json and returns a library of chain codes.
 */
func InitializeChainCodes() LibChainCodes {

	file, err := os.Open("/home/nishi/go/src/obcsdk/util/CC_Collection.json")
	if err != nil {
		log.Fatal("Error in opening CC_Collection.json file ")
	}

	poolChainCode, err := unmarshalChainCodes(file)
	if err != nil {
		log.Fatal("Error in unmarshalling")
	}

	//make a map to hold each chaincode detail
	ChCos := make(map[string]ChainCode)
	for i := 0; i < len(poolChainCode); i++ {
		//construct a map for each chaincode detail
		detail := make(map[string]string)
		detail["type"] = poolChainCode[i].TYPE
		detail["path"] = poolChainCode[i].PATH
		//detail["dep_txid"] = poolChainCode[i].DEP_TXID
		//detail["deployed"] = poolChainCode[i].DEPLOYED

		versions := make(map[string]string)
		CC := ChainCode{Detail: detail, Versions: versions}
		//add the structure to map of chaincodes
		ChCos[poolChainCode[i].NAME] = CC
	}
	//finally add this map - collection of chaincode detail to library of chaincodes
	libChainCodes := LibChainCodes{ChainCodes: ChCos}
	return libChainCodes
}

func initializePeers() (peers []Peer, name string) {

	peerDetails, userDetails, Name := initNetworkCredentials()
	//userDetails := initializeUsers()
	numOfPeersOnNetwork := len(peerDetails)
	numOfUsersOnNetwork := len(userDetails)
	fmt.Println("Num of Peers", numOfPeersOnNetwork)
	fmt.Println("Num of Users", numOfUsersOnNetwork)
	fmt.Println("Name of network", Name)

	allPeers := make([]Peer, numOfPeersOnNetwork)

	factor := numOfUsersOnNetwork / numOfPeersOnNetwork
	remainder := numOfUsersOnNetwork % numOfPeersOnNetwork
	i := 0
	k := 0
	//for each peerDetail we construct a new peer evenly distributing the list of users
	for i < numOfPeersOnNetwork {
		aPeer := new(Peer)
		aPeerDetail := make(map[string]string)
		name := "vp" + strconv.Itoa(i)
		aPeerDetail["ip"] = peerDetails[i].IP
		aPeerDetail["port"] = peerDetails[i].PORT
		aPeerDetail["name"] = name

		//fmt.Println(" value in i", i , "k ", k)

		j := 0
		userInfo := make(map[string]string)
		for j < factor {
			for k < numOfUsersOnNetwork {
				//fmt.Println(" **********value in inside i", i , "k ", k, "factor", factor, " j ", j)
				userInfo[userDetails[k].USER] = userDetails[k].SECRET
				j++
				k++
				if j == factor {
					break
				}
			}
		}

		aPeer.PeerDetails = aPeerDetail
		aPeer.UserData = userInfo
		allPeers[i] = *aPeer
		i++
	}
	//do we have any left over users details
	if remainder > 0 {
		for m := 0; m < remainder; m++ {
			allPeers[m].UserData[userDetails[k].USER] = userDetails[k].SECRET
			k++
		}
	}
	return allPeers, Name
}

func initNetworkCredentials() ([]peerData, []userData, string) {

	file, err := os.Open("/home/nishi/go/src/obcsdk/util/NetworkCredentials.json")
	if err != nil {
		fmt.Println("Error in opening NetworkCredentials file ", err)
		log.Fatal("Error in opening Network Credential json File")

	}
	NC, err := unmarshalNetworkCredentials(file)
	if err != nil {
		log.Fatal("Error in unmarshalling")
	}
	//peerdata := make(map[string]string)
	peerData := NC.PEERLIST
	userData := NC.USERLIST
	name := NC.NAME
	return peerData, userData, name
}
