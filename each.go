package eacht

type each struct {
	errs         []error
	bufferLength int
	callbackFunc f
	channel      chan []interface{}
	channelClose chan struct{}
	buffer       []interface{}
}

type Iterator interface {
	run()
	HasError() bool
	Add(line interface{}) bool
	Close()
	GetErrors() *[]error
}

type f func([]interface{}, bool) error

func NewEach(bufferLength int, callbackFunc f) Iterator {
	var e Iterator

	e = &each{
		channel:      make(chan []interface{}),
		buffer:       make([]interface{}, 0, bufferLength),
		bufferLength: bufferLength,
		callbackFunc: callbackFunc,
		channelClose: make(chan struct{}),
	}

	go e.run()

	return e
}

func (e *each) HasError() bool {
	return len(e.errs) > 0
}
func (e *each) run() {
	for c := range e.channel {
		err := e.callbackFunc(c, e.HasError())
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	close(e.channelClose)
}

func (e *each) Add(line interface{}) bool {
	e.buffer = append(e.buffer, line)
	if len(e.buffer) >= e.bufferLength {
		e.flush()
	}

	return e.HasError()
}

func (e *each) Close() {
	if len(e.buffer) > 0 {
		e.flush()
	}

	close(e.channel)
	<-e.channelClose
}

func (e *each) flush() {
	e.channel <- e.buffer
	e.buffer = make([]interface{}, 0, e.bufferLength)
}

func (e *each) GetErrors() *[]error {
	return &e.errs
}
