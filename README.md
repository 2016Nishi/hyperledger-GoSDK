# obc-test
This repository contains test-cases.

The test cases will be opened as Issues in this repository.  The title format for the Issue should include the name of the chaincode to be run, application, or automated script details.  Also, please include the key labels to identify the nodes (peers) as well as the configuration options to be set in the obc-peer/openchain.yaml file. Options not specified using a label should be assumed to be default configuration values. If the test fails, the corresponding Issue should **not** be opened in this project but in the project where the failing code exists.  For example, if the fabric hangs, the Issue should be opened against the obc-peer repository and the Issue number placed in the template.  The defect template that should be used is the [Bug Template](BugTemplate.md).  The application/chaincode used for the test should be placed in the testfiles folder and the name of the file(s) include the test case name, for example, "Car application.*"

Once a test is passed, the tester should add the "Passed" label to the Issue.  If the test fails, the tester should add the "Fail" label to the Issue.   The Passed/Fail label should apply to the last person who ran the test.

The following template should be used when a test is run and placed in the Issue:

**Application 
Automated?   
Commit Level   Results   
Issue Number (If applicable)**

**Test Steps:**
{Anything needed to run the tests, payloads (json), scripts and included in the folder you created.
If your test artifacts live in another git project or are hosted and can be downloaded, plesse include the url and build/commit level that was used for this test}

**Expected Results:**
Application and chaincode tests


    a.	Ensure published APIs are all working
        i.	GET/ chain, GET/chain/blocks/, POST/devops/deploy, POST/devops/invoke, POST/devops/query
    b.	 Measure concurrency and scalability limits (not staffed)
        i.	Number of chaincodes 
        ii.	Increasing chaincode calls (loop)
        iii.	Track blockchain and chaincode size 
        iv.	Using Trie and Bucketlist for state
    c.	Long Runs
        i.	Variances in step c. from 1 hour to 1 week
            1.	Monitor block size, memory, cpu, error messages
    d.	Negative Testing
        i.	While running the tests bounce the nodes randomly
            1.	Ensure resiliency - ledger state
    e.	Ensure Bluemix service is used in testing and stable for Interconnect
        i.	Chain-code and apps run in Bluemix
        ii.	General usability
        
 


---------------------------------------------------------------------------------------
So for example,

Issue: Chaincode Example 2 noops single node trie

**Application** Postman   
**Automated?**  No  
**Commit Level**   9e21a1e29fe5cb65715a699d0063cbfacbba6c2d                           
**Results**     Pass
**Issue Number if Failing Test**

**Test Steps:** 

Deploy:
{
  "type": "GOLANG",
  "chaincodeID": {
    "path": "https://hub.jazz.net/git/averyd/cc_ex02/chaincode_example02"
  },
  "ctorMsg": {
    "function": "init",
    "args": [
      "a", "100", "b", "200"
    ]
  }
} 

Invoke & Query:

{
  "chaincodeSpec": {
  "type": "GOLANG",
  "chaincodeID": {
    "name": "09c08b4770c13cbbd877eb96728a4644068f77be1b85489f9b51700c046429af214d43c8d4428422f7a5c785d5afdc3837e9e86abaab3da1f5786d353418d307"
  },
  "ctorMsg": {
    "function": "invoke",
    "args": [
      "a", "b", "20"
    ]
  }
  }
} {
  "chaincodeSpec": {
  "type": "GOLANG",
  "chaincodeID": {
    "name": "09c08b4770c13cbbd877eb96728a4644068f77be1b85489f9b51700c046429af214d43c8d4428422f7a5c785d5afdc3837e9e86abaab3da1f5786d353418d307"
  },
  "ctorMsg": {
    "function": "query",
    "args": [
      "a"
    ]
  }
  }
} 


** Expected Results:**

Application and chaincode tests

    i)	Single and Two Node Setup (validating and non-validating nodes) – using multiple configs (yaml)
        a.	Ensure Deploy, Invoke, and Query work properly  - Worked
         i.	Use data structure for trie and bucketlist for state - Worked with Both
            1.	Reference https://github.com/openblockchain/obc-peer/blob/master/openchain.yaml, state and dataStructure
    b.	Ensure published APIs are all working
        i.	GET/ chain, GET/chain/blocks/, POST/devops/deploy, POST/devops/invoke, POST/devops/query - Worked
    c.	 Measure concurrency and scalability limits (not staffed)
        i.	Number of chaincodes 
        Ran a single chaincode
        ii.	Increasing chaincode calls (loop)
        I had a loop driving the chaincode doing invoke and query from two clients
        iii.	Track blockchain and chaincode size 
        I ran for 15 minutes and the blocks grew to x.
        iv.	Using Trie and Bucketlist for state
        I performed the looping test with both state structures
    d.	Long Runs
        i.	Variances in step c. from 1 hour to 1 week
            1.	Monitor block size, memory, cpu, error messages
            I will kick off a weekend run and monitor the two nodes
        
            
    e.	Negative Testing
        i.	While running the tests bounce the nodes randomly
            1.	Ensure resiliency
            I manually bounced peers during the test and checked state and key/values were correct.
    f.	Ensure Bluemix service is used in testing and stable for Interconnect
        i.	Chain-code and apps run in Bluemix
        N/A I created my own fabric in my lab
        ii.	General usability
        N/A
