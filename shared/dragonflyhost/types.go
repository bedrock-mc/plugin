package dragonflyhost

type Position struct {
	X     float64
	Y     float64
	Z     float64
	World string
}

type WorldRef struct {
	ID        string
	Name      string
	Dimension string
}

type ItemStack struct {
	TypeID string
	Name   string
	Meta   int16
	Count  int
}

type InventorySlot struct {
	Slot int
	Item ItemStack
}

type PlayerState struct {
	Position   Position
	Health     float64
	MaxHealth  float64
	Gamemode   string
	XPLevel    int
	XPProgress float64
}

type DamageKind string

const (
	DamageKindAttack     DamageKind = "attack"
	DamageKindProjectile DamageKind = "projectile"
	DamageKindFall       DamageKind = "fall"
	DamageKindVoid       DamageKind = "void"
	DamageKindDrowning   DamageKind = "drowning"
	DamageKindSuffocate  DamageKind = "suffocate"
	DamageKindLava       DamageKind = "lava"
	DamageKindFire       DamageKind = "fire"
	DamageKindCustom     DamageKind = "custom"
)

type DamageSource struct {
	Type        string
	Description string
	Kind        DamageKind
	Fire        bool
	DamagerUUID string
	DamagerName string
}
