package eacht

type each struct {
	errs         []error
	bufferLength int
	callbackFunc f
	channel      chan []string
	channelClose chan struct{}
	buffer       []string
}

type Iterator interface {
	run()
	HasError() bool
	Add(line string) bool
	Close()
	GetErrors() *[]error
}

type f func([]string, bool) error

func NewEach(bufferLength int, callbackFunc f) Iterator {
	var e Iterator

	e = &each{
		channel:      make(chan []string),
		buffer:       make([]string, 0, bufferLength),
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

func (e *each) Add(line string) bool {
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
	e.buffer = make([]string, 0, e.bufferLength)
}

func (e *each) GetErrors() *[]error {
	return &e.errs
}
