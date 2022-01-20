package events

import "fmt"

type dispatcherAction = func(l *LogDispatcher)

type logConsumer struct {
	name string
	sink chan<- LogEntry
}
type LogDispatcher struct {
	input  chan dispatcherAction
	output []logConsumer
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
				select {
				case listener.sink <- e:
				default:
					fmt.Printf("WARNING: Channel %q would block, dropping message\n", listener.name)
				}
			}
		}
	}
}

func (d *LogDispatcher) Add(name string, out chan<- LogEntry) func() {
	d.input <- func(l *LogDispatcher) {
		l.output = append(d.output, logConsumer{
			name: name,
			sink: out,
		})
	}
	return func() {
		d.input <- func(l *LogDispatcher) {
			for i, consumer := range l.output {
				if consumer.name == name && consumer.sink == out {
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
