package main

import ( 
)

type Peer struct {
	addr		 string
	LastPingSentTime int64
	LastPingReceivedTime int64

	VersionSent bool
	VersionReceived bool
	VerackSent bool
	VerackReceived bool
}

func tt(){

	//var nodeList = []*Peer{&Peer{addr:"NODE_A"}}

	var addrList []string
	//addrList = askAddrFrom("NODE_A")
	//var tmpAddrList []string
	
	for i := 0; i < len(addrList); i++ {
		/*if addrList[i].IsContainedIn(nodeList) == false {
			tmpAddrList = append(tmpAddrList, addrList[i])
		}*/



	}


}