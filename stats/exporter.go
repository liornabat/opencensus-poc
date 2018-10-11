package stats

import (
	"fmt"

	"go.opencensus.io/stats/view"
)

const indent = "  "

type PrintExporter struct{}

// ExportView logs the view data.
func (e *PrintExporter) ExportView(vd *view.Data) {

	for _, row := range vd.Rows {
		fmt.Printf("%v %-45s", vd.End.Format("15:04:05"), vd.View.Name)
		switch v := row.Data.(type) {
		case *view.DistributionData:
			fmt.Printf("distribution: min=%.1f max=%.1f mean=%.1f", v.Min, v.Max, v.Mean)
		case *view.CountData:
			fmt.Printf("count:        value=%v", v.Value)
		case *view.SumData:
			fmt.Printf("sum:          value=%v", v.Value)
		case *view.LastValueData:
			fmt.Printf("last:         value=%v", v.Value)
		}
		fmt.Println()

		for _, tag := range row.Tags {
			fmt.Printf("%v- %v=%v\n", indent, tag.Key.Name(), tag.Value)
		}
	}
}
