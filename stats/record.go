package stats

type Recorder interface {
	Record(item ...Item) error
}
