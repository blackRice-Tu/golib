package xkafka

import (
	"github.com/blackRice-Tu/golib"

	"github.com/Shopify/sarama"
)

var (
	V0_8_2_0  = sarama.V0_8_2_0
	V0_8_2_1  = sarama.V0_8_2_1
	V0_8_2_2  = sarama.V0_8_2_2
	V0_9_0_0  = sarama.V0_9_0_0
	V0_9_0_1  = sarama.V0_9_0_1
	V0_10_0_0 = sarama.V0_10_0_0
	V0_10_0_1 = sarama.V0_10_0_1
	V0_10_1_0 = sarama.V0_10_1_0
	V0_10_1_1 = sarama.V0_10_1_1
	V0_10_2_0 = sarama.V0_10_2_0
	V0_10_2_1 = sarama.V0_10_2_1
	V0_11_0_0 = sarama.V0_11_0_0
	V0_11_0_1 = sarama.V0_11_0_1
	V0_11_0_2 = sarama.V0_11_0_2
	V1_0_0_0  = sarama.V1_0_0_0
	V1_1_0_0  = sarama.V1_1_0_0
	V1_1_1_0  = sarama.V1_1_1_0
	V2_0_0_0  = sarama.V2_0_0_0
	V2_0_1_0  = sarama.V2_0_1_0
	V2_1_0_0  = sarama.V2_1_0_0
	V2_2_0_0  = sarama.V2_2_0_0
	V2_3_0_0  = sarama.V2_3_0_0

	OffsetNewest = sarama.OffsetNewest
	OffsetOldest = sarama.OffsetOldest
)

func init() {
	mode := golib.GetMode()
	switch mode {
	case golib.DebugMode:
		sarama.Logger = golib.NewThirdPartyLogger("sarama")
	case golib.ReleaseMode:
	}
}
