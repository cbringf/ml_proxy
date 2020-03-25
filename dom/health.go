package dom

import (
	"encoding/json"
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
	var snapshot = make([]SysRequestStats, 0)

	if len(reqInfoList) <= 0 {
		return &snapshot
	}

	var slices = make([][2]int, 0)

	var currentPair = [2]int{0, 1}

	y, m, d := reqInfoList[0].RequestDate.Date()
	_, minute, _ := reqInfoList[0].RequestDate.Clock()

	sort.Slice(reqInfoList, func(i, j int) bool {
		return reqInfoList[i].RequestDate.Before(reqInfoList[i].RequestDate)
	})

	for i := 1; i < len(reqInfoList); i++ {
		y1, m1, d1 := reqInfoList[i].RequestDate.Date()
		_, minute1, _ := reqInfoList[i].RequestDate.Clock()

		if y != y1 || m != m1 || d != d1 || minute != minute1 {
			currentPair[1] = i + 1
			slices = append(slices, currentPair)
			currentPair = [2]int{i + 1, i + 1}
			y, m, d, minute = y1, m1, d1, minute1
		}
	}
	if currentPair[1] < len(reqInfoList)-1 {
		currentPair[1] = len(reqInfoList)
		slices = append(slices, currentPair)
	}

	for _, s := range slices {
		snapshot = append(snapshot, buildSnapshotFromSlice(reqInfoList[s[0]:s[1]]))
	}

	return &snapshot
}

func buildSnapshotFromSlice(reqInfoList []RequestInfo) SysRequestStats {
	stats := SysRequestStats{
		Date:         reqInfoList[0].RequestDate,
		InfoRequests: make([]SysReponseStats, 0),
	}
	resCodes := make(map[int]int)

	for i := 0; i < len(reqInfoList); i++ {
		temp := reqInfoList[i]

		stats.AvgResponseTime += temp.ResponseTime
		stats.TotalRequests++

		if temp.Remote {
			stats.AvgResponseTimeAPICalls += temp.RemoteResponseTime
			stats.TotalCountAPICalls++
			resCodes[temp.RemoteResponseStatus]++
		}
	}
	if stats.TotalRequests > 0 {
		stats.AvgResponseTime = stats.AvgResponseTime / stats.TotalRequests
	}
	if stats.TotalCountAPICalls > 0 {
		stats.AvgResponseTimeAPICalls = stats.AvgResponseTimeAPICalls / stats.TotalCountAPICalls
	}
	for code, count := range resCodes {
		stats.InfoRequests = append(stats.InfoRequests, SysReponseStats{
			Count:      count,
			StatusCode: code,
		})
	}

	return stats
}
