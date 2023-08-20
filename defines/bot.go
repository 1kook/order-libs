package defines

type (
	BotStatus int32
)

const (
	BotStatusActive        BotStatus = 1
	BotStatusStopped       BotStatus = 2
	BotStatusPause         BotStatus = 3
	BotStatusRunningWithEx BotStatus = 4

	BotTypeBLSH   = "0"
	BotTypeCustom = "1"
	BotTypeMDCA   = "2"

	BotNameBLSH   = "buylowsellhigh"
	BotNameCustom = "custom"
	BotNameMDCA   = "mdca"
)

var (
	MapBotName = map[string]string{
		BotTypeBLSH:   BotNameBLSH,
		BotTypeCustom: BotNameCustom,
		BotTypeMDCA:   BotNameMDCA,
	}
)
