package consts

import "nam-club/NumBuyer_back/models/orgerrors"

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
	AbilityTypeRecover AbilityType = "recover"
	AbilityTypeJam     AbilityType = "jam"
	AbilityTypeConfuse AbilityType = "confuse"
	// 実行状態
	AbilityStatusUnused AbilityStatus = "unused"
	AbilityStatusReady  AbilityStatus = "ready"
	AbilityStatusActive AbilityStatus = "active"
	AbilityStatusUsed   AbilityStatus = "used"
	// ID
	AbilityIdFiBoost     = "bst_prm_001"
	AbilityIdNumViolence = "atk_prm_001"
	AbilityIdReboot      = "rcv_tmp_001"
	AbilityIdShutdown    = "jam_prm_001"
	AbilityIdCatastrophe = "cnf_tmp_001"
)

var (
	// keyにID, valueにアビリティ情報
	abilities = map[string]Ability{
		AbilityIdFiBoost:     {AbilityIdFiBoost, AbilityTriggerPassive, AbilityTypeBoost, -1, AbilityStatusReady},
		AbilityIdNumViolence: {AbilityIdNumViolence, AbilityTriggerPassive, AbilityTypeAttack, -1, AbilityStatusReady},
		AbilityIdReboot:      {AbilityIdReboot, AbilityTriggerActive, AbilityTypeRecover, -1, AbilityStatusUnused},
		AbilityIdShutdown:    {AbilityIdShutdown, AbilityTriggerPassive, AbilityTypeJam, -1, AbilityStatusReady},
		AbilityIdCatastrophe: {AbilityIdCatastrophe, AbilityTriggerActive, AbilityTypeConfuse, 1, AbilityStatusUnused},
	}
)

func GetAbilities() []Ability {
	ret := []Ability{}
	for _, v := range abilities {
		ret = append(ret, v)
	}

	return ret
}

func ParseAbility(s string) (Ability, error) {
	if val, ok := abilities[s]; ok {
		return val, nil
	} else {
		return Ability{}, orgerrors.NewValidationError("ability parse error. " + s)
	}
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
