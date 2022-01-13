package events

type dispatcherAction = func(l *LogDispatcher)

type LogDispatcher struct {
	input  chan dispatcherAction
	output []chan<- LogEntry
}

func NewLogDispatcher() *LogDispatcher {
	input := make(chan dispatcherAction, 16)
	dispatcher := &LogDispatcher{
		input:  input,
		output: nil,
	}
	go dispatcher.runLoop()
	return dispatcher
}

func (d *LogDispatcher) runLoop() {
	for action := range d.input {
		action(d)
	}
}

func (d *LogDispatcher) Consume(in <-chan LogEntry) {
	for e := range in {
		d.input <- func(l *LogDispatcher) {
			for _, listener := range d.output {
				listener <- e
			}
		}
	}
}

func (d *LogDispatcher) Add(out chan<- LogEntry) func() {
	d.input <- func(l *LogDispatcher) {
		l.output = append(d.output, out)
	}
	return func() {
		d.input <- func(l *LogDispatcher) {
			for i, ch := range l.output {
				if ch == out {
					left := l.output[0:i]
					right := l.output[i+1:]
					l.output = append(left, right...)
					return
				}
			}
			panic("Could not find channel -- may result in later runtime panics")
		}
	}
}
