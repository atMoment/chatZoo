package main

const (
	ComponentBase = "ComponentBase"
)

func NewBaseComponent(room IRoomBase) *_BaseComponent {
	BaseComponent := &_BaseComponent{
		IRoomBase: room,
	}
	return BaseComponent
}

type _BaseComponent struct {
	IRoomBase
}

func (r *_BaseComponent) PlayerJoinRoom(member string) error {
	return r.JoinRoom(member)
}
