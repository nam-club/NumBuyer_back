package consts

import "nam-club/NumBuyer_back/models/orgerrors"

type AbilityStatus string
type AbilityTrigger string
type AbilityTiming string
type AbilityType string
type Ability struct {
	ID            string
	Trigger       AbilityTrigger
	Timing        AbilityTiming
	Type          AbilityType
	UsableNum     int // -1なら無制限
	InitialStatus AbilityStatus
}

const (
	// 実行のされ方
	AbilityTriggerActive  AbilityTrigger = "active"
	AbilityTriggerPassive AbilityTrigger = "passive"
	// 実行タイミング
	AbilityTimingSoon AbilityTiming = "soon"
	AbilityTimingWait AbilityTiming = "wait"
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
	AbilityIdReload      = "rcv_tmp_001"
	AbilityIdShutdown    = "jam_prm_001"
	AbilityIdCatastrophe = "cnf_tmp_001"
)

var (
	// keyにID, valueにアビリティ情報
	abilities = map[string]Ability{
		AbilityIdFiBoost: {
			ID:            AbilityIdFiBoost,
			Trigger:       AbilityTriggerPassive,
			Timing:        AbilityTimingWait,
			Type:          AbilityTypeBoost,
			UsableNum:     -1,
			InitialStatus: AbilityStatusReady},
		AbilityIdNumViolence: {
			ID:            AbilityIdNumViolence,
			Trigger:       AbilityTriggerPassive,
			Timing:        AbilityTimingWait,
			Type:          AbilityTypeAttack,
			UsableNum:     -1,
			InitialStatus: AbilityStatusReady},
		AbilityIdReload: {
			ID:            AbilityIdReload,
			Trigger:       AbilityTriggerActive,
			Timing:        AbilityTimingSoon,
			Type:          AbilityTypeRecover,
			UsableNum:     -1,
			InitialStatus: AbilityStatusUnused},
		AbilityIdShutdown: {
			ID:            AbilityIdShutdown,
			Trigger:       AbilityTriggerPassive,
			Timing:        AbilityTimingSoon,
			Type:          AbilityTypeJam,
			UsableNum:     -1,
			InitialStatus: AbilityStatusReady},
		AbilityIdCatastrophe: {
			ID:            AbilityIdCatastrophe,
			Trigger:       AbilityTriggerActive,
			Timing:        AbilityTimingWait,
			Type:          AbilityTypeConfuse,
			UsableNum:     2,
			InitialStatus: AbilityStatusUnused},
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
		return Ability{}, orgerrors.NewValidationError("ability.parseError", "ability parse error", map[string]string{"abilidyId": s})
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
