package executable

type LabelPool struct {
	// use index as label's name, store address's offset in code byte
	labels []int
}

func NewLabelPool() *LabelPool {
	return &LabelPool{
		labels: []int{},
	}
}

func (labelPool *LabelPool) NewLabel() int {
	labelPool.labels = append(labelPool.labels, -1)
	return len(labelPool.labels) - 1
}

func (labelPool *LabelPool) SetLabel(label int, address int) {
	labelPool.labels[label] = address
}