package kafka

import (
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jinzhu/gorm"
	"github.com/keyo-oliveira/codepix/application/factory"
	appmodel "github.com/keyo-oliveira/codepix/application/model"
	"github.com/keyo-oliveira/codepix/application/usecases"
	"github.com/keyo-oliveira/codepix/domain/model"
	"os"
)

type KafKaProcessor struct {
	Database     *gorm.DB
	Producer     *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafKaProcessor {
	return &KafKaProcessor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

func (k *KafKaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := ckafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}

	c.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			k.processMessage(msg)
		}
	}
}

func (k *KafKaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (k *KafKaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
	)

	if err != nil {
		fmt.Println("Error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJson, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)

	if err != nil {
		return err
	}
	return nil
}

func (k *KafKaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)
	if transaction.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}

	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)

		if err != nil {
			return err
		}
	}

	return nil

}

func (k *KafKaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecases.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)

	if err != nil {
		return err
	}

	return nil
}
