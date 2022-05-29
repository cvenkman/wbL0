package stan

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"
)


type Stan struct {
	Channel		string
	ClientID	string
	ClusterID	string
}

func New(Channel, ClientID, ClusterID string) *Stan {
	return &Stan{
		Channel: Channel,
		ClientID: ClientID,
		ClusterID: ClientID,
	}
}

// connect and subscribe o stan
func (st *Stan) InitStan(msgHandler stan.MsgHandler) (sub stan.Subscription, sc stan.Conn) {
	sc, err := stan.Connect(st.ClusterID, st.ClientID)
	if err != nil {
		log.Fatalf("Can't connect to stan %s", err.Error())
	}

	sub, err = sc.Subscribe(st.Channel, msgHandler,
				stan.SetManualAckMode(), stan.AckWait(6 * time.Second),
				stan.DurableName(st.ClientID), stan.StartWithLastReceived())
	if err != nil {
		sc.Close()
		log.Fatalf("Can't subscribe to %s channel: %s", st.Channel, err.Error())
	}
	log.Println("Subscriber listening on channel", st.Channel)

	return
}