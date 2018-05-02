package server

import (
	"log"
	"tway/serverutil"
)

func (s *Server) NewMsgHeaders(addrTo *serverutil.NetAddress, list []serverutil.Header, getHeadersRequest *serverutil.MsgAskHeaders) *serverutil.MsgHeaders {
	return &serverutil.MsgHeaders{s.ipStatus, addrTo, s.version, list, getHeadersRequest}
}

func (s *Server) sendHeaders(addrTo *serverutil.NetAddress, list []serverutil.Header, getHeadersRequest *serverutil.MsgAskHeaders) ([]byte, error) {
	s.Log(true, "headers sent to:", addrTo.String())
	headers := s.NewMsgHeaders(addrTo, list, getHeadersRequest)
	//assigne en []byte la structure getblocks
	payload := gobEncode(*headers)
	//on append la commande et le payload
	request := append(commandToBytes("headers"), payload...)
	err := s.sendData(addrTo.String(), request)
	return request, err
}

//Receptionne une liste header de block
//voir structure MsgHeaders
func (s *Server) handleHeaders(request []byte) {
	var payload serverutil.MsgHeaders
	if err := getPayload(request, &payload); err != nil {
		log.Panic(err)
	}
	addr := payload.AddrSender.String()
	s.Log(true, "headers received from:", addr)

	go s.HistoryManager.NewHeadersHistory(&payload)

	p, _ := s.GetPeer(addr)
	p.IncreaseBytesReceived(uint64(len(request)))
	s.AddPeer(p)

}
