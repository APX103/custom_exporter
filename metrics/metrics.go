package metrics

type Metrics interface {
	Update() error
}
