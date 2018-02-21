package server

import (
	"time"
)

//return average time client do to download one block
func GetAverageTimeToDownloadABlock(dis map[string]*DownloadInformations) int64 {
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	var cpt int64 = 0
	var total_time int64 = 0 
	for _, di := range dis {
		if di.start > oneHourAgo.UnixNano() && di.receivedAt != 0{
			cpt++
			total_time += di.receivedAt - di.start
		}
	}
	if cpt == 0 {
		return 0
	}
	return total_time / cpt
}




