package stats

import (
	"fmt"

	"go.opencensus.io/stats/view"
)

type exporter struct {
	aggMap *aggMap
}

func NewExporter() *exporter {
	return &exporter{
		aggMap: newAggMap(),
	}
}

func (e *exporter) ExportView(vd *view.Data) {
	for _, row := range vd.Rows {
		key := makeKeyFromTags(row.Tags)
		index := fmt.Sprintf("%s@@%s", key.String(), vd.View.Name)
		switch v := row.Data.(type) {
		case *view.DistributionData:
			e.aggMap.insert(index, v.Count, v.Sum())
		case *view.CountData:
			e.aggMap.insert(index, v.Value)
		case *view.SumData:
			e.aggMap.insert(index, v.Value)
		case *view.LastValueData:
			e.aggMap.insert(index, v.Value)
		}

	}
}
