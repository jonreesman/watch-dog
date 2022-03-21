package pb

import "google.golang.org/grpc"

type sentimentClientHandler struct {
	sentClient SentimentClient
}

func (h *sentimentClientHandler) establishConnection(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//Handle error
	}
	h.sentClient = NewSentimentClient(conn)
}
