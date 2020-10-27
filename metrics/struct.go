package metrics

import (
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"strconv"
	"strings"
	"time"

	"encoding/json"
)

type MetricsData struct {
	Dimensions  string `json:"Dimensions"`
	Project     string `json:"Project"`
	Environment string `json:"Environment"`
	//Service          string `json:"Service"`
	MeasureName      string `json:"MeasureName"`
	MeasureValue     string `json:"MeasureValue"`
	MeasureValueType string `json:"MeasureValueType"`
	Time             string `json:"Time"`
}

func (md *MetricsData) FromCostRow(cr costs.CostRow) {
	md.Dimensions = "dimensions"
	md.Project = "CM4" //cr.Account.Name
	//md.Service = cr.Service
	md.Environment = cr.Account.Environment
	md.MeasureName = "cost"
	md.MeasureValue = cr.Cost
	md.MeasureValueType = "DOUBLE"
	mytime, _ := time.Parse(dates.AWSDateFormat(), cr.Date)
	t := mytime.UnixNano() / int64(time.Millisecond)
	md.Time = strconv.FormatInt(t, 10)
}

type MetricsRecord struct {
	Data      string `json:"data"`
	Partition string `json:"partition-key"`
}

type MetricsPutData struct {
	Records []MetricsRecord `json:"records"`
}

func FromCosts(costs []costs.CostRow) ([]byte, error) {
	mpd := MetricsPutData{}

	for _, c := range costs {
		record := MetricsRecord{}
		record.Partition = "some key"
		data := MetricsData{}
		data.FromCostRow(c)
		j, _ := json.Marshal(data)
		d := string(j)
		record.Data = strings.ReplaceAll(d, `"`, `'`)
		mpd.Records = append(mpd.Records, record)

	}

	mpd.Records = mpd.Records[0:1]
	return json.Marshal(mpd)
}
