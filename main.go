package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	HOST = "http://88.198.13.202:9051"
)

func main() {
	ergoNodeClient := initializeClient()
	//resp, err := ergoNodeClient.GetBlocks(&GetBlocksRequest{Offset: 10000})
	//resp, err := ergoNodeClient.GetBlock("Vumq5gex8Ty3TuAk8Xxxc9UmRgRd64pnRxvV3PM7Q4Q")
	//if err != nil {
	//	panic(err)
	//}
	//decoded, err := DecodeToBig([]byte(resp.Header.Id))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%b\n", decoded)
	resp, err := ergoNodeClient.GetBlocks(&GetBlocksRequest{Limit: 50, Offset: 0})
	if err != nil {
		panic(err)
	}

	blocks := make([]Block, 0)
	for i := 0; i < len(resp.BlocksIds); i++ {
		block, err := ergoNodeClient.GetBlock(string(resp.BlocksIds[i]))
		if err != nil {
			panic(err)
		}
		blocks = append(blocks, *block)
	}
	proofs, lastBlocks := Prove(&Chain{blocks})
	fmt.Println("Proofs:")
	for _, proof := range proofs {
		fmt.Printf("%+v\n", proof)
	}
	fmt.Println("LastBlocks:")
	for _, lastBlock := range lastBlocks {
		fmt.Printf("%+v\n", lastBlock)
	}

	//verified, proofs := Verify(proofs, lastBlocks)
}

func initializeClient() *ErgoNodeClient {
	httpClient := http.Client{Timeout: time.Second}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar().With(zap.String("logger", "ErgoNodeClient"))

	client := ErgoNodeClient{
		URL:    HOST,
		Logger: logger,
		Client: &httpClient,
	}
	return &client
}
