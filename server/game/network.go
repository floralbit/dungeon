package game

func (z *zone) send(event serverEvent) {
	for _, e := range z.Entities {
		if e.Type != entityTypePlayer {
			continue
		}

		e.client.In <- event
	}
}

func (e *entity) send(event serverEvent) {
	if e.Type != entityTypePlayer {
		return
	}

	e.client.In <- event
}