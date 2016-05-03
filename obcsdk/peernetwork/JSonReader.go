package peernetwork

import (
	"encoding/json"
	"fmt"
	"io"
	//"os"
	//"log"
)

type userData struct {
	USER   string `json:"username"`
	SECRET string `json:"secret"`
}

type peerData struct {
	IP   string `json:"api-host"`
	PORT string `json:"api-port"`
}

type networkCredentials struct {
	PEERLIST []peerData `json:"PeerData"`
	USERLIST []userData `json:"UserData"`
	NAME     string     `json:"Name"`
}

type chainCodeData struct {
	NAME           string `json:"name"`
	TYPE           string `json:"type"`
	PATH           string `json:"path"`
	//DEP_TXID       string `json:"dep_txid"`
	//DEPLOYED       string `json:"deployed"`
}


/*
  converts input stream to NetworkCredentials
	reader is an open file
 */
func unmarshalNetworkCredentials(reader io.Reader) (networkCredentials, error) {

	decoder := json.NewDecoder(reader)
	//fmt.Println("Inside Unmarshal Network Credentials JSONREADER")
	var NC networkCredentials
	err := decoder.Decode(&NC)
	if err != nil {
		fmt.Println("Error in decoding ")
	}
	return NC, err
}

/*
  converts input stream to chaincodes.
 */
func unmarshalChainCodes(reader io.Reader) ([]*chainCodeData, error) {

	decoder := json.NewDecoder(reader)
	//fmt.Println("Inside Unmarshal ChainCode JSONREADER")
	var ChainCodeCollection []*chainCodeData
	err := decoder.Decode(&ChainCodeCollection)
	if err != nil {
		fmt.Println("Error in decoding ")
	}
	return ChainCodeCollection, err
}
