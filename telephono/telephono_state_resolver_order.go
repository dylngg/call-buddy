package telephono

type OrderedVariableResolver []VariableResolver

func (o OrderedVariableResolver) get(index int) VariableResolver {
	return o[index]
}

func (o OrderedVariableResolver) Len() int {
	return len(o)
}

func (o OrderedVariableResolver) Less(i, j int) bool {
	var (
		firstVal  = o.get(i).Order()
		secondVal = o.get(j).Order()
	)

	return firstVal < secondVal
}

func (o OrderedVariableResolver) Swap(i, j int) {
	temp := o.get(i)
	o[i] = o[j]
	o[j] = temp
}

