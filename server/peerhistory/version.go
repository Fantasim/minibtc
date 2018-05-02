package peerhistory

import (
	"time"
	"tway/serverutil"
)

type Version map[string]ListVersionHistory
type ListVersionHistory []VersionHistory

type VersionHistory struct {
	Message *serverutil.MsgVersion
	Date    time.Time
	Sent    bool
}

//Créer un historique de requete getheaders
func (hm *HistoryManager) NewVersionHistory(msg *serverutil.MsgVersion, sent bool) {
	hm.muVersion.Lock()
	defer hm.muVersion.Unlock()
	var addr string
	if sent == true {
		addr = msg.AddrReceiver.String()
	} else {
		addr = msg.AddrSender.String()
	}
	vh := VersionHistory{msg, time.Now(), sent}
	hm.Version[addr] = append(hm.Version[addr], vh)
}
