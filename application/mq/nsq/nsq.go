package nsq

import (
	"encoding/json"
	"go-nsq/application/mq"
	"log"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type INSQClient interface {
	Publish(string, []byte) error
	Subscribe(string) error
}

type Message struct {
	Timestamp    string
	FileObjectID string
	FileName     string
}

type NSQMessageHandler struct{}

// TODO : Apply message format from MQ to update the specified document at mongodb
func processMessage(body []byte) error {
	log.Printf("Receiving message from NSQ with payload : %v ", string(body))
	return nil
}

func (h *NSQMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	var response mq.Message
	// do whatever actual message processing is desired
	err := json.Unmarshal(m.Body, &response)
	if err != nil {
		log.Printf("Error when unmarshalling json at NSQMessagehandler with error : %v", err)
		return err
	}

	log.Println("Logging message from NSQMessageHandler")
	log.Println(response.FileName)
	log.Println(response.FileObjectID)
	log.Println(response.Timestamp)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

type NSQClient struct {
	config nsq.Config
	msgCounter prometheus.Counter
}

func NewNSQClient() INSQClient {
	config := nsq.NewConfig()
	// after adding config.DialTimeout, NSQ will not throw i/o timeout anymore
	config.DialTimeout = 3 * time.Second
	reg := prometheus.NewRegistry()
	msgCounter := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name:      "message_emmited_count",
		Help:      "Number of message pumped",
	})
	// Register msgCounter metric
	prometheus.Register(msgCounter)
	return &NSQClient{
		config: *config,
		msgCounter: msgCounter,
	}
}

func (n NSQClient) Publish(topic string, message []byte) error {
	publisher, err := nsq.NewProducer("127.0.0.1:4150", &n.config)
	if err != nil {
		return err
	}

	err = publisher.Publish(topic, message)
	if err != nil {
		return err
	}
	n.msgCounter.Inc()
	return nil
}

func (n NSQClient) Subscribe(topic string) error {
	nsqSubscriber, err := nsq.NewConsumer(topic, "channel", &n.config)
	if err != nil {
		return err
	}
	nsqSubscriber.AddHandler(&NSQMessageHandler{})

	// either localhost or 127.0.0.1 as address are acceptable, but prefere 127.0.0.1 for consistency
	nsqSubscriber.ConnectToNSQLookupd("127.0.0.1:4161")

	return nil
}
