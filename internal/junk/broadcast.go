package junk

//StringBroadcast copies an input string from a single channel to multiple target channels.  Mutation is not suitable
//in a multi-threaded context.
type StringBroadcast struct {
	Input <-chan string
	Out   []chan<- string
}

func (s *StringBroadcast) RunLoop() {
	for str := range s.Input {
		for _, target := range s.Out {
			target <- str
		}
	}
}
