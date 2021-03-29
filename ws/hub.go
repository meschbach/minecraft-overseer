package ws

type Hub struct {
	ingest chan Message
	register chan *Spoke
	unregister chan *Spoke
	clients map[*Spoke]bool
	overseer *StateMachine
	broadcast chan string
}

func NewHub() *Hub  {
	broadcast := make(chan string)
	return &Hub{
		ingest:  make(chan Message),
		register:   make(chan *Spoke),
		unregister: make(chan *Spoke),
		clients:    make(map[*Spoke]bool),
		broadcast: broadcast,
		overseer: newStateMachine(broadcast),
	}
}

func (h *Hub) Run() {
	go h.overseer.run()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.message)
			}
		case message := <-h.ingest:
			result := message.apply(h)
			for client := range h.clients {
				select {
				case client.message <- result:
				default:
					go func() { h.unregister <- client }()
				}
			}
		case stringMessage := <-h.broadcast:
			outputMessage := StringOutput{message: stringMessage}
			for client := range h.clients {
				client.message <- outputMessage
			}
		}
	}
}
