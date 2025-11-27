package dispute

type Dispute struct {
	Id    string
	State string
}

func New(id string, state string) *Dispute {
	return &Dispute{Id: id, State: state}
}
