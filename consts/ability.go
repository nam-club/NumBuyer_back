package consts

type AbilityStatus string
type AbilityTrigger string
type AbilityType string
type Ability struct {
	ID            string
	Trigger       AbilityTrigger
	Type          AbilityType
	UsableNum     int // -1なら無制限
	InitialStatus AbilityStatus
}

const (
	// 実行のされ方
	AbilityTriggerActive  AbilityTrigger = "active"
	AbilityTriggerPassive AbilityTrigger = "passive"
	// 区分
	AbilityTypeBoost   AbilityType = "boost"
	AbilityTypeAttack  AbilityType = "attack"
	AbilityTypeDefense AbilityType = "defense"
	AbilityTypeJam     AbilityType = "jam"
	AbilityTypeConfuse AbilityType = "confuse"
	// 実行状態
	AbilityStatusUnused AbilityStatus = "unused"
	AbilityStatusReady  AbilityStatus = "ready"
	AbilityStatusActive AbilityStatus = "active"
	AbilityStatusUsed   AbilityStatus = "used"
	// ID
	AbilityIdFiBoost       = "boost_prm_001"
	AbilityIdNumViolence   = "atk_prm_001"
	AbilityIdBringYourself = "def_tmp_001"
	AbilityIdShutdown      = "jam_prm_001"
	AbilityIdShakeShake    = "cnf_tmp_001"
)

var (
	// keyにID, valueにアビリティ情報
	abilities = map[string]Ability{
		AbilityIdFiBoost:       {AbilityIdFiBoost, AbilityTriggerPassive, AbilityTypeBoost, -1, AbilityStatusUnused},
		AbilityIdNumViolence:   {AbilityIdNumViolence, AbilityTriggerPassive, AbilityTypeAttack, -1, AbilityStatusUnused},
		AbilityIdBringYourself: {AbilityIdBringYourself, AbilityTriggerActive, AbilityTypeDefense, 5, AbilityStatusActive},
		AbilityIdShutdown:      {AbilityIdShutdown, AbilityTriggerPassive, AbilityTypeJam, -1, AbilityStatusUnused},
		AbilityIdShakeShake:    {AbilityIdShakeShake, AbilityTriggerActive, AbilityTypeConfuse, 1, AbilityStatusActive},
	}
)

func GetAbilities() []Ability {
	ret := []Ability{}
	for _, v := range abilities {
		ret = append(ret, v)
	}

	return ret
}

func ParseAbilities(s []string) []Ability {
	ret := []Ability{}
	for _, v := range s {
		if val, ok := abilities[v]; ok {
			ret = append(ret, val)
		}
	}
	return ret
}
