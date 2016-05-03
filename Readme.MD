Notice regarding the Linux Foundation's Hyperledger project

The openblockchain project is IBM's proposed contribution to the Linux Foundation's Hyperledger project. We have made it available as open source to enable others to explore our architecture and design. IBM's intention is to engage rigorously in the Linux Foundation's Hyperledger project as the community establishes itself, and decides on a code base. Once established, we will transition our development focus to the Hyperledger effort, and this code will be maintained as needed for IBM's use.
While we invite contribution to the openblockchain project, we believe that the broader blockchain community's focus should be the Hyperledger project.

openblockchain/obc-test/obcsdk/
--------------------------------------------

This project contains the code to build applications to test chaincodes from openchain obc-peer fabric.

Building the project

1. Create a fork of the obc-test repository using the GitHub web interface. 

2. Next, clone your fork in the appropriate location.

	`cd $GOPATH/src` 

	`git clone https://github.com/2016Nishi/hyperledger-GoSDK.git
	

3. To build code :

 P.S. Change the path location to json files `"NetworkCredentials.json"` and `"CC_COllection.json"` inside `obcsdk/peernetwork/peerNetworkSetup.go` to your local path. Currently it is set as 
 
 `/home/nishi/go/src/obcsdk/util.` on my unix system
 
 
4. To Run sample test program:

	`cd $GOPATH/src/obcsdk/chcotest`
	
	`go run RestCalls.go Deploy_With_Tagname.go`

5. To see what commands are available, or to generate go doc:

	`cd $GOPATH/src`
	

	`go install obcsdk/chaincode`
	
	           `obcsdk/peernetwork`
	      
	           `obcsdk/peerrest`

	launch git bash command

	`cd $GOPATH/src`

	`godoc server godoc -http=":6060"`

6. To view go doc on browser:
	Open a browser and launch `http://localhost:6060/pkg/obcsdk` 
	above url lists available obcsdk API and it's usage.
