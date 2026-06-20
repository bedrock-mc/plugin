package dragonflyhost

import (
	"fmt"

	dfblock "github.com/df-mc/dragonfly/server/block"
	dfentity "github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func DamageSourceSnapshot(src world.DamageSource) *DamageSource {
	if src == nil {
		return nil
	}
	desc := fmt.Sprint(src)
	info := &DamageSource{
		Type:        fmt.Sprintf("%T", src),
		Description: desc,
		Kind:        DamageKindCustom,
		Fire:        src.Fire(),
	}
	switch s := src.(type) {
	case dfentity.AttackDamageSource:
		info.Kind = DamageKindAttack
		if p, ok := s.Attacker.(*player.Player); ok {
			info.DamagerUUID = p.UUID().String()
			info.DamagerName = p.Name()
		}
	case dfentity.ProjectileDamageSource:
		info.Kind = DamageKindProjectile
		if p, ok := s.Owner.(*player.Player); ok {
			info.DamagerUUID = p.UUID().String()
			info.DamagerName = p.Name()
		}
	case dfentity.FallDamageSource, dfentity.GlideDamageSource:
		info.Kind = DamageKindFall
	case dfentity.VoidDamageSource:
		info.Kind = DamageKindVoid
	case dfentity.DrowningDamageSource:
		info.Kind = DamageKindDrowning
	case dfentity.SuffocationDamageSource:
		info.Kind = DamageKindSuffocate
	case dfblock.LavaDamageSource:
		info.Kind = DamageKindLava
	case dfblock.FireDamageSource:
		info.Kind = DamageKindFire
	default:
		if info.Fire {
			info.Kind = DamageKindFire
		}
	}
	return info
}
