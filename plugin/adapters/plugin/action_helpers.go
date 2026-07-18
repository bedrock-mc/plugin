package plugin

import (
	"time"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func soundFromProto(s pb.Sound) world.Sound {
	switch s {
	case pb.Sound_ATTACK:
		return sound.Attack{Damage: true}
	case pb.Sound_DROWNING:
		return sound.Drowning{}
	case pb.Sound_BURNING:
		return sound.Burning{}
	case pb.Sound_FALL:
		return sound.Fall{Distance: 4}
	case pb.Sound_BURP:
		return sound.Burp{}
	case pb.Sound_POP:
		return sound.Pop{}
	case pb.Sound_EXPLOSION:
		return sound.Explosion{}
	case pb.Sound_THUNDER:
		return sound.Thunder{}
	case pb.Sound_LEVEL_UP:
		return sound.LevelUp{}
	case pb.Sound_EXPERIENCE:
		return sound.Experience{}
	case pb.Sound_FIREWORK_LAUNCH:
		return sound.FireworkLaunch{}
	case pb.Sound_FIREWORK_HUGE_BLAST:
		return sound.FireworkHugeBlast{}
	case pb.Sound_FIREWORK_BLAST:
		return sound.FireworkBlast{}
	case pb.Sound_FIREWORK_TWINKLE:
		return sound.FireworkTwinkle{}
	case pb.Sound_TELEPORT:
		return sound.Teleport{}
	case pb.Sound_ARROW_HIT:
		return sound.ArrowHit{}
	case pb.Sound_ITEM_BREAK:
		return sound.ItemBreak{}
	case pb.Sound_ITEM_THROW:
		return sound.ItemThrow{}
	case pb.Sound_TOTEM:
		return sound.Totem{}
	case pb.Sound_FIRE_EXTINGUISH:
		return sound.FireExtinguish{}
	default:
		return sound.Pop{}
	}
}

func particleFromProto(act *pb.WorldAddParticleAction) (world.Particle, bool) {
	return particleFromType(act.GetParticle(), act.Block, act.Face)
}

// particleFromType maps a particle enum plus optional block/face into a world.Particle.
func particleFromType(pt pb.ParticleType, blk *pb.BlockState, f *int32) (world.Particle, bool) {
	switch pt {
	case pb.ParticleType_PARTICLE_HUGE_EXPLOSION:
		return particle.HugeExplosion{}, true
	case pb.ParticleType_PARTICLE_ENDERMAN_TELEPORT:
		return particle.EndermanTeleport{}, true
	case pb.ParticleType_PARTICLE_SNOWBALL_POOF:
		return particle.SnowballPoof{}, true
	case pb.ParticleType_PARTICLE_EGG_SMASH:
		return particle.EggSmash{}, true
	case pb.ParticleType_PARTICLE_SPLASH:
		return particle.Splash{}, true
	case pb.ParticleType_PARTICLE_EFFECT:
		return particle.Effect{}, true
	case pb.ParticleType_PARTICLE_ENTITY_FLAME:
		return particle.EntityFlame{}, true
	case pb.ParticleType_PARTICLE_FLAME:
		return particle.Flame{}, true
	case pb.ParticleType_PARTICLE_DUST:
		return particle.Dust{}, true
	case pb.ParticleType_PARTICLE_BLOCK_FORCE_FIELD:
		return particle.BlockForceField{}, true
	case pb.ParticleType_PARTICLE_BONE_MEAL:
		return particle.BoneMeal{}, true
	case pb.ParticleType_PARTICLE_EVAPORATE:
		return particle.Evaporate{}, true
	case pb.ParticleType_PARTICLE_WATER_DRIP:
		return particle.WaterDrip{}, true
	case pb.ParticleType_PARTICLE_LAVA_DRIP:
		return particle.LavaDrip{}, true
	case pb.ParticleType_PARTICLE_LAVA:
		return particle.Lava{}, true
	case pb.ParticleType_PARTICLE_DUST_PLUME:
		return particle.DustPlume{}, true
	case pb.ParticleType_PARTICLE_BLOCK_BREAK:
		if b, ok := blockFromProto(blk); ok {
			return particle.BlockBreak{Block: b}, true
		}
	case pb.ParticleType_PARTICLE_PUNCH_BLOCK:
		if b, ok := blockFromProto(blk); ok {
			face := cube.FaceUp
			if f != nil {
				face = cube.Face(*f)
			}
			return particle.PunchBlock{Block: b, Face: face}, true
		}
	}
	return nil, false
}

func particleFromPlayerAction(act *pb.PlayerShowParticleAction) (world.Particle, bool) {
	return particleFromType(act.GetParticle(), act.Block, act.Face)
}

func playerTitleFromAction(act *pb.SendTitleAction) title.Title {
	t := title.New(act.Title)
	if act.Subtitle != nil && *act.Subtitle != "" {
		t = t.WithSubtitle(*act.Subtitle)
	}
	if act.FadeInMs != nil {
		t = t.WithFadeInDuration(time.Duration(*act.FadeInMs) * time.Millisecond)
	}
	if act.DurationMs != nil {
		t = t.WithDuration(time.Duration(*act.DurationMs) * time.Millisecond)
	}
	if act.FadeOutMs != nil {
		t = t.WithFadeOutDuration(time.Duration(*act.FadeOutMs) * time.Millisecond)
	}
	return t
}

// World query handlers
