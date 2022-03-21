package consts

type AbilityStatus string
type AbilityTrigger string
type AbilityType string
type Ability struct {
	ID        string
	Trigger   AbilityTrigger
	Type      AbilityType
	UsableNum int // -1なら無制限
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
	AbilityStatusActive AbilityStatus = "active"
	AbilityStatusReady  AbilityStatus = "ready"
	AbilityStatusUnused AbilityStatus = "unused"
	AbilityStatusUsed   AbilityStatus = "used"
)

var (
	// keyにID, valueにアビリティ情報
	abilities = map[string]Ability{
		//FiBoost
		"boost_prm_001": {"boost_prm_001", AbilityTriggerPassive, AbilityTypeBoost, -1},
		// NumViolence
		"atk_prm_001": {"atk_prm_001", AbilityTriggerPassive, AbilityTypeAttack, -1},
		// BringYourself
		"def_tmp_001": {"def_tmp_001", AbilityTriggerActive, AbilityTypeDefense, 5},
		// Shutdown
		"jam_prm_001": {"jam_prm_001", AbilityTriggerPassive, AbilityTypeJam, -1},
		// ShakeShake
		"cnf_tmp_001": {"cnf_tmp_001", AbilityTriggerActive, AbilityTypeConfuse, 1},
	}
)

func ParseAbilities(s []string) []Ability {
	ret := []Ability{}
	for _, v := range s {
		if val, ok := abilities[v]; ok {
			ret = append(ret, val)
		}
	}
	return ret
}
