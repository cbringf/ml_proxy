package dom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type SysReponseStats struct {
	StatusCode int `json:"status_code"`
	Count      int `json:"count"`
}

type SysRequestStats struct {
	Date                    time.Time         `json:"date"`
	AvgResponseTime         int               `json:"avg_response_time"`
	AvgResponseTimeAPICalls int               `json:"avg_response_time_api_calls"`
	TotalRequests           int               `json:"total_requests"`
	TotalCountAPICalls      int               `json:"total_count_api_calls"`
	InfoRequests            []SysReponseStats `json:"info_requests"`
}

type SysRequestSnapshot struct {
	SysRequestList *[]SysRequestStats
}

func (sysReq SysRequestSnapshot) HandleRequest(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(sysReq.SysRequestList)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func BuildSnapshot(reqInfoList []RequestInfo) *[]SysRequestStats {
	if len(reqInfoList) <= 0 {
		return nil
	}

	var snapshot = make([]SysRequestStats, 0)

	sort.Slice(reqInfoList, func(i, j int) bool {
		return reqInfoList[i].RequestDate.Before(reqInfoList[i].RequestDate)
	})
	_, minute, _ := reqInfoList[0].RequestDate.Clock()
	fmt.Println(reqInfoList[0].RequestDate.Clock())
	perMinuteStats := SysRequestStats{
		Date:         reqInfoList[0].RequestDate,
		InfoRequests: make([]SysReponseStats, 0),
	}
	resCodes := make(map[int]int)

	for i := 0; i < len(reqInfoList); i++ {
		temp := reqInfoList[i]

		if _, m, _ := temp.RequestDate.Clock(); minute == m {
			perMinuteStats.AvgResponseTime += temp.ResponseTime
			perMinuteStats.AvgResponseTimeAPICalls += temp.RemoteResponseTime
			perMinuteStats.TotalRequests++

			if temp.Remote {
				perMinuteStats.TotalCountAPICalls++
				resCodes[temp.RemoteResponseStatus]++
			}
		} else {
			for code, count := range resCodes {
				perMinuteStats.InfoRequests = append(perMinuteStats.InfoRequests, SysReponseStats{
					Count:      count,
					StatusCode: code,
				})
			}
			perMinuteStats.AvgResponseTime = perMinuteStats.AvgResponseTime / perMinuteStats.TotalRequests
			perMinuteStats.AvgResponseTimeAPICalls = perMinuteStats.AvgResponseTimeAPICalls / perMinuteStats.TotalCountAPICalls

			snapshot = append(snapshot, perMinuteStats)
			perMinuteStats = SysRequestStats{
				Date:         temp.RequestDate,
				InfoRequests: make([]SysReponseStats, 0),
			}
			_, minute, _ = reqInfoList[i].RequestDate.Clock()
			resCodes = make(map[int]int)
			i--
		}
	}
	return &snapshot
}
