package main

import "errors"

func interpretMessage(message string) (error, Message) {
	if( message == "stop") {
		return nil, &StopMessage{}
	} else if( message == "start" ){
		return nil, &StartMessage{}
	} else {
		return errors.New("unknown message"), nil
	}
}
