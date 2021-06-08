package eacht

func Iterate(rows []string, bufferLength int, callbackFunc f) *[]error {
	var errs *[]error
	e := NewEach(bufferLength, callbackFunc)

	for _, row := range rows {
		hasError := e.Add(row)
		if true == hasError {
			break
		}
	}
	e.Close()

	if e.HasError() {
		errs = e.GetErrors()
	}

	return errs
}
