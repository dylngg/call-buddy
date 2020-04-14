package telephono

type ContextContributor interface {
	/*
		Contribute returns the name of the contribution and the object (Map or struct)

		For example, the environment variable contributor should return "Env", a map of key values, nil
	*/
	Contribute() (string, interface{}, error)
}
