package server

import (
	"time"
)

//Cette fonction retourne le temps moyen que le noeud courant met pour télécharger
//un block en nanosec
func GetAverageTimeToDownloadABlock(dis map[string]*DownloadInformations) int64 {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	var cpt int64 = 0
	var total_time int64 = 0
	for _, di := range dis {
		if di.start > oneHourAgo.UnixNano() && di.receivedAt != 0 {
			cpt++
			total_time += di.receivedAt - di.start
		}
	}
	if cpt == 0 {
		return 0
	}
	return total_time / cpt
}

//Cette fonction retourne le temps moyen que le noeud courant met pour télécharger
//n header en nanosec
func (s *Server) GetAverageTimeToGetNHeaders(count int) int64 {
	//map d'historique de requete getheaders
	getHeadersHistory := s.HistoryManager.GetHeader
	var cpt int64 = 0
	var total_time int64 = 0

	//pour chaque pair ayant recu ou envoyé une requete getheaders avec le noeud courant
	for peerAddr, listGetHeadersReq := range getHeadersHistory {
		//on récupère l'historique de la derniere requete headers recu
		lastHeadersReq := s.HistoryManager.Headers.GetListHeadersHistoryByAddr(peerAddr).SelectByCount(count).SortByDate(true).First()
		if lastHeadersReq != nil {
			//on recupere la liste d'historique de requete getheaders envoyé, triée par date
			listGetHeadersReq = listGetHeadersReq.SelectBySent(true).SortByDate(true).SelectByCount(count)

			//pour chaque requete getheaders
			for _, getHeadersReq := range listGetHeadersReq {
				//si la requete getheaders a été effectué avant la reception de la derniere requete headers
				if lastHeadersReq.Date.After(getHeadersReq.Date) {
					total_time += lastHeadersReq.Date.UnixNano() - getHeadersReq.Date.UnixNano()
					cpt++
					break
				}
			}

		}
	}
	if cpt == 0 {
		return 0
	}
	return total_time / cpt

}
