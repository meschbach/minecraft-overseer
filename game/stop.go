package game

type stopCommand struct {
}

func (*stopCommand) run(state *internalState, game *Game) error {
	if state.state != runningState {
		return nil
	}
	state.state = stoppingState

	state.commands <- "/stop"
	_, err := state.serviceProcess.Process.Wait()
	if err != nil {
		return err
	}
	state.state = idleState
	return nil
}
