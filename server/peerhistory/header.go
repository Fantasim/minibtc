package peerhistory

import (
	"time"
	"tway/serverutil"
)

type ListHeadersHistory []GetHeadersHistory

type GetHeadersHistory struct {
	Message *serverutil.MsgAskHeaders
	Date    time.Time
	Sent    bool
}

//Créer un historique de requete getheaders
func (hm *HistoryManager) NewGetHeadersHistory(msg *serverutil.MsgAskHeaders, sent bool) {
	hm.muGetHeader.Lock()
	defer hm.muGetHeader.Unlock()
	var addr string
	if sent == true {
		addr = msg.AddrReceiver.String()
	} else {
		addr = msg.AddrSender.String()
	}
	hm.GetHeader[addr] = append(hm.GetHeader[addr], GetHeadersHistory{msg, time.Now(), sent})
}
