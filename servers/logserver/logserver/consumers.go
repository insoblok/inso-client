package logserver

import (
	contract "eth-toy-client/core/contracts"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/logbus"
	"log"
)

type ConsoleConsumer struct {
	Name             string
	ContractRegistry *contract.Registry
	Events           chan logbus.LogEvent
}

func (consumer *ConsoleConsumer) Consume() {
	for event := range consumer.Events {
		contractAddress := toytypes.ContractAddress{
			Address: event.Contract,
		}
		log.Printf(
			"🚀 %s received event: %v, address:%s, TxHash:%s and Args:%v",
			consumer.Name,
			event.LogType,
			event.Contract,
			event.TxHash,
			event.Args)
		meta, ok := consumer.ContractRegistry.Get(contractAddress)
		if !ok {
			log.Println("It seems i do not have any info about this contract, %v", contractAddress)
		} else {
			log.Println("I have info about this contract, %v", len(meta.ABI))
		}

	}
}
