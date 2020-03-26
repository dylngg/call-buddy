package telephono

type SortableVariableResolver []VariableResolver

func (o SortableVariableResolver) get(index int) VariableResolver {
	return o[index]
}

func (o SortableVariableResolver) Len() int {
	return len(o)
}

func (o SortableVariableResolver) Less(i, j int) bool {
	var (
		firstVal  = o.get(i).Order()
		secondVal = o.get(j).Order()
	)

	return firstVal < secondVal
}

func (o SortableVariableResolver) Swap(i, j int) {
	temp := o.get(i)
	o[i] = o[j]
	o[j] = temp
}
