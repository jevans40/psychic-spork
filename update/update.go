package update

//This is the real project here

type Entity interface {
	//Main update tick
	Update() []byte
	GetReply(reply []byte)
	ResolveDeltas() ([28]float64, bool)
	SendFullVertex() [28]float64
}

func Update(Render chan string) {
	//Make agents here
}
