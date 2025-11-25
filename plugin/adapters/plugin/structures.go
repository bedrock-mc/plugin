package plugin

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/world"
	pb "github.com/secmc/plugin/proto/generated/go"
)

// protoStructure implements world.Structure for sparse voxel data sent over protobuf.
type structureVoxel struct {
	block  world.Block
	liquid world.Liquid
}

type protoStructure struct {
	w, h, l int
	vox     map[[3]int]structureVoxel
}

func (s *protoStructure) Dimensions() [3]int {
	return [3]int{s.w, s.h, s.l}
}

func (s *protoStructure) At(x, y, z int, _ func(x, y, z int) world.Block) (world.Block, world.Liquid) {
	if v, ok := s.vox[[3]int{x, y, z}]; ok {
		return v.block, v.liquid
	}
	// nils mean: do nothing at this coordinate
	return nil, nil
}

// buildProtoStructure converts a proto StructureDef into a protoStructure wrapper.
// It validates dimensions and voxels and returns an error on invalid definitions.
func buildProtoStructure(def *pb.StructureDef) (*protoStructure, error) {
	if def == nil {
		return nil, fmt.Errorf("missing structure")
	}
	if def.Width <= 0 || def.Height <= 0 || def.Length <= 0 {
		return nil, fmt.Errorf("structure dimensions must be positive")
	}
	ps := &protoStructure{
		w:   int(def.Width),
		h:   int(def.Height),
		l:   int(def.Length),
		vox: make(map[[3]int]structureVoxel, len(def.Voxels)),
	}
	for _, v := range def.Voxels {
		if v == nil || v.Block == nil {
			continue
		}
		x, y, z := int(v.X), int(v.Y), int(v.Z)
		// Drop out-of-bounds voxels early.
		if x < 0 || y < 0 || z < 0 || x >= ps.w || y >= ps.h || z >= ps.l {
			continue
		}
		blk, ok := blockFromProto(v.Block)
		if !ok {
			return nil, fmt.Errorf("unknown block in voxel at (%d,%d,%d)", x, y, z)
		}
		var liq world.Liquid
		if v.Liquid != nil && v.Liquid.Block != nil {
			if lb, ok := blockFromProto(v.Liquid.Block); ok {
				if lq, ok := lb.(world.Liquid); ok {
					liq = lq
				}
			}
		}
		ps.vox[[3]int{x, y, z}] = structureVoxel{block: blk, liquid: liq}
	}
	return ps, nil
}
