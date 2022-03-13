package consts

type AbilityStatus string
type AbilityType string
type AbilityAttribute string
type Ability struct {
	ID        string
	Type      AbilityType
	Attribute AbilityAttribute
	UsableNum int // -1なら無制限
}

const (
	// 実行のされ方
	AbilityTypeActive  AbilityType = "active"
	AbilityTypePassive AbilityType = "passive"
	// 区分
	AbilityAttributeBoost   AbilityAttribute = "boost"
	AbilityAttributeAttack  AbilityAttribute = "attack"
	AbilityAttributeDefense AbilityAttribute = "defense"
	AbilityAttributeJam     AbilityAttribute = "jam"
	AbilityAttributeConfuse AbilityAttribute = "confuse"
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
		"boost_prm_001": {"boost_prm_001", AbilityTypePassive, AbilityAttributeBoost, -1},
		// NumViolence
		"atk_prm_001": {"atk_prm_001", AbilityTypePassive, AbilityAttributeAttack, -1},
		// BringYourself
		"def_tmp_001": {"def_tmp_001", AbilityTypeActive, AbilityAttributeDefense, 5},
		// Shutdown
		"jam_prm_001": {"jam_prm_001", AbilityTypePassive, AbilityAttributeJam, -1},
		// ShakeShake
		"cnf_tmp_001": {"cnf_tmp_001", AbilityTypeActive, AbilityAttributeConfuse, 1},
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
