package types

// Query aiDataSources supported by the provider querier. Eg: custom provider query oScript
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryReward = "reward"
)

// QueryResReward Queries a complete block reward
type QueryResReward struct {
	Reward `json:"block_reward"`
}

// NewQueryResReward is the constructor for the block reward
func NewQueryResReward(reward Reward) QueryResReward {
	return QueryResReward{
		Reward: reward,
	}
}
