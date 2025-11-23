package plugin

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/df-mc/dragonfly/server/item/category"
	"github.com/df-mc/dragonfly/server/world"
	pb "github.com/secmc/plugin/proto/generated/go"
)

// customItem implements world.CustomItem
type customItem struct {
	id           string
	displayName  string
	texture      image.Image
	itemCategory category.Category
	meta         int16
}

func (c *customItem) EncodeItem() (name string, meta int16) {
	return c.id, c.meta
}

func (c *customItem) Name() string {
	return c.displayName
}

func (c *customItem) Texture() image.Image {
	return c.texture
}

func (c *customItem) Category() category.Category {
	return c.itemCategory
}

// registerCustomItems registers custom items declared in PluginHello
func (m *Manager) registerCustomItems(p *pluginProcess, defs []*pb.CustomItemDefinition) {
	if len(defs) == 0 {
		return
	}
	pluginName := p.id
	if hello := p.helloInfo(); hello != nil && hello.Name != "" {
		pluginName = hello.Name
	}
	for _, def := range defs {
		if def == nil {
			continue
		}
		if err := m.registerSingleCustomItem(def); err != nil {
			p.log.Error("failed to register custom item", "id", def.Id, "error", err)
			continue
		}
		m.log.Info("registered custom item", "plugin", pluginName, "id", def.Id, "name", def.DisplayName)
	}
}

func (m *Manager) registerSingleCustomItem(def *pb.CustomItemDefinition) error {
	if def.Id == "" {
		return fmt.Errorf("custom item ID cannot be empty")
	}

	if def.DisplayName == "" {
		return fmt.Errorf("custom item display name cannot be empty")
	}

	if len(def.TextureData) == 0 {
		return fmt.Errorf("custom item texture data cannot be empty")
	}

	// Decode texture from PNG bytes
	img, err := png.Decode(bytes.NewReader(def.TextureData))
	if err != nil {
		return fmt.Errorf("decode texture: %w", err)
	}

	// Convert proto category to Dragonfly category
	cat := convertProtoCategory(def.Category)

	// Apply optional group if provided
	if def.Group != nil && *def.Group != "" {
		cat = cat.WithGroup(*def.Group)
	}

	// Register with Dragonfly
	// Note: This will panic if the item is already registered
	world.RegisterItem(&customItem{
		id:           def.Id,
		displayName:  def.DisplayName,
		texture:      img,
		itemCategory: cat,
		meta:         int16(def.Meta),
	})

	return nil
}

func convertProtoCategory(cat pb.ItemCategory) category.Category {
	switch cat {
	case pb.ItemCategory_ITEM_CATEGORY_CONSTRUCTION:
		return category.Construction()
	case pb.ItemCategory_ITEM_CATEGORY_NATURE:
		return category.Nature()
	case pb.ItemCategory_ITEM_CATEGORY_EQUIPMENT:
		return category.Equipment()
	case pb.ItemCategory_ITEM_CATEGORY_ITEMS:
		return category.Items()
	default:
		return category.Items()
	}
}
