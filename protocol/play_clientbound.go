// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate protocol_builder $GOFILE Play clientbound

package protocol

import (
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/encoding/nbt"
)

// KeepAliveClientbound is sent by a server to check if the
// client is still responding and keep the connection open.
// The client should reply with the KeepAliveServerbound
// packet setting ID to the same as this one.
//
// Currently the packet id is: 0x00
type KeepAliveClientbound struct {
	ID VarInt
}

// JoinGame is sent after completing the login process. This
// sets the initial state for the client.
//
// Currently the packet id is: 0x01
type JoinGame struct {
	// The entity id the client will be referenced by
	EntityID int32
	// The starting gamemode of the client
	Gamemode byte
	// The dimension the client is starting in
	Dimension int8
	// The difficuilty setting for the server
	Difficulty byte
	// The max number of players on the server
	MaxPlayers byte
	// The level type of the server
	LevelType string
	// Whether the client should reduce the amount of debug
	// information it displays in F3 mode
	ReducedDebugInfo bool
}

// ServerMessage is a message sent by the server. It could be from a player
// or just a system message. The Type field controls the location the
// message is displayed at and when the message is displayed.
//
// Currently the packet id is: 0x02
type ServerMessage struct {
	Message chat.AnyComponent `as:"json"`
	// 0 - Chat message, 1 - System message, 2 - Action bar message
	Type byte
}

// TimeUpdate is sent to sync the world's time to the client, the client
// will manually tick the time itself so this doesn't need to sent repeatedly
// but if the server or client has issues keeping up this can fall out of sync
// so it is a good idea to sent this now and again
//
// Currently the packet id is: 0x03
type TimeUpdate struct {
	WorldAge  int64
	TimeOfDay int64
}

// EntityEquipment is sent to display an item on an entity, like a sword
// or armor. Slot 0 is the held item and slots 1 to 4 are boots, leggings
// chestplate and helmet respectively.
//
// Currently the packet id is: 0x04
type EntityEquipment struct {
	EntityID VarInt
	Slot     int16
	Item     ItemStack `as:"raw"`
}

// SpawnPosition is sent to change the player's current spawn point. Currently
// only used by the client for the compass.
//
// Currently the packet id is: 0x05
type SpawnPosition struct {
	Location Position
}

// UpdateHealth is sent by the server to update the player's health and food.
//
// Currently the packet id is: 0x06
type UpdateHealth struct {
	Health         float32
	Food           VarInt
	FoodSaturation float32
}

// Respawn is sent to respawn the player after death or when they move worlds.
//
// Currently the packet id is: 0x07
type Respawn struct {
	Dimension  int32
	Difficulty byte
	Gamemode   byte
	LevelType  string
}

// TeleportPlayer is sent to change the player's position. The client is expected
// to reply to the server with the same positions as contained in this packet
// otherwise will reject future packets.
//
// Currently the packet id is: 0x08
type TeleportPlayer struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      byte
}

// SetCurrentHotbarSlot changes the player's currently selected hotbar item.
//
// Currently the packet id is: 0x09
type SetCurrentHotbarSlot struct {
	Slot byte
}

// EntityUsedBed is sent by the server when a player goes to bed.
//
// Currently the packet id is: 0x0A
type EntityUsedBed struct {
	EntityID VarInt
	Location Position
}

// Animation is sent by the server to play an animation on a specific entity.
//
// Currently the packet id is: 0x0B
type Animation struct {
	EntityID    VarInt
	AnimationID byte
}

// SpawnPlayer is used to spawn a player when they are in range of the client.
// This packet alone isn't enough to display the player as the skin and username
// information is in the player information packet.
//
// Currently the packet id is: 0x0C
type SpawnPlayer struct {
	EntityID    VarInt
	UUID        UUID `as:"raw"`
	X, Y, Z     int32
	Yaw, Pitch  int8
	CurrentItem int16
	Metadata    Metadata
}

// CollectItem causes the collected item to fly towards the collector. This
// does not destroy the entity.
//
// Currently the packet id is: 0x0D
type CollectItem struct {
	CollectedEntityID VarInt
	CollectorEntityID VarInt
}

// SpawnObject is used to spawn an object or vehicle into the world when it
// is in range of the client. Velocity is only sent if the Data field is
// non-zero.
//
// Currently the packet id is: 0x0E
type SpawnObject struct {
	EntityID                        VarInt
	Type                            byte
	X, Y, Z                         int32
	Pitch, Yaw                      int8
	Data                            int32
	VelocityX, VelocityY, VelocityZ int16 `if:".Data != 0"`
}

// SpawnMob is used to spawn a living entity into the world when it is in
// range of the client.
//
// Currently the packet id is: 0x0F
type SpawnMob struct {
	EntityID                        VarInt
	Type                            byte
	X, Y, Z                         int32
	Yaw, Pitch                      int8
	HeadPitch                       int8
	VelocityX, VelocityY, VelocityZ int16
	Metadata                        Metadata
}

// SpawnPainting spawns a painting into the world when it is in range of
// the client. The title effects the size and the texture of the painting.
//
// Currently the packet id is: 0x10
type SpawnPainting struct {
	EntityID  VarInt
	Title     string
	Location  Position
	Direction byte
}

// SpawnExperienceOrb spawns a single experience orb into the world when
// it is in range of the client. The count controls the amount of experience
// gained when collected.
//
// Currently the packet id is: 0x11
type SpawnExperienceOrb struct {
	EntityID VarInt
	X, Y, Z  int32
	Count    int16
}

// EntityVelocity sets the velocity of an entity in 1/8000 of a block
// per a tick.
//
// Currently the packet id is: 0x12
type EntityVelocity struct {
	EntityID                        VarInt
	VelocityX, VelocityY, VelocityZ int16
}

// EntityDestroy destroys the entities with the ids in the provided slice.
//
// Currently the packet id is: 0x13
type EntityDestroy struct {
	EntityIDs []VarInt `length:"VarInt"`
}

// Entity does nothing. It is a result of subclassing used in Minecraft.
//
// Currently the packet id is: 0x14
type Entity struct {
	EntityID VarInt
}

// EntityMove moves the entity with the id by the offsets provided.
//
// Currently the packet id is: 0x15
type EntityMove struct {
	EntityID               VarInt
	DeltaX, DeltaY, DeltaZ int8
	OnGround               bool
}

// EntityLook rotates the entity to the new angles provided.
//
// Currently the packet id is: 0x16
type EntityLook struct {
	EntityID   VarInt
	Yaw, Pitch int8
	OnGround   bool
}

// EntityLookAndMove is a combination of EntityMove and EntityLook.
//
// Currently the packet id is: 0x17
type EntityLookAndMove struct {
	EntityID               VarInt
	DeltaX, DeltaY, DeltaZ int8
	Yaw, Pitch             int8
	OnGround               bool
}

// EntityTeleport teleports the entity to the target location. This is
// sent if the entity moves further than EntityMove allows.
//
// Currently the packet id is: 0x18
type EntityTeleport struct {
	EntityID   VarInt
	X, Y, Z    int32
	Yaw, Pitch int8
	OnGround   bool
}

// EntityHeadLook rotates an entity's head to the new angle.
//
// Currently the packet id is: 0x19
type EntityHeadLook struct {
	EntityID VarInt
	HeadYaw  int8
}

// EntityAction causes an entity to preform an action based on the passed
// id.
//
// Currently the packet id is: 0x1A
type EntityAction struct {
	EntityID int32
	ActionID byte
}

// EntityAttach attaches to entities together, either by mounting or leashing.
// -1 can be used at the EntityID to deattach.
//
// Currently the packet id is: 0x1B
type EntityAttach struct {
	EntityID int32
	Vehicle  int32
	Leash    bool
}

// EntityMetadata updates the metadata for an entity.
//
// Currently the packet id is: 0x1C
type EntityMetadata struct {
	EntityID VarInt
	Metadata Metadata
}

// EntityEffect applies a status effect to an entity for a given duration.
//
// Currently the packet id is: 0x1D
type EntityEffect struct {
	EntityID      VarInt
	EffectID      int8
	Amplifier     int8
	Duration      VarInt
	HideParticles bool
}

// EntityRemoveEffect removes an effect from an entity.
//
// Currently the packet id is: 0x1E
type EntityRemoveEffect struct {
	EntityID VarInt
	EffectID int8
}

// SetExperience updates the experience bar on the client.
//
// Currently the packet id is: 0x1F
type SetExperience struct {
	ExperienceBar   float32
	Level           VarInt
	TotalExperience VarInt
}

// EntityProperties updates the properties for an entity.
//
// Currently the packet id is: 0x20
type EntityProperties struct {
	EntityID   VarInt
	Properties []EntityProperty `length:"int32"`
}

// EntityProperty is a key/value pair with optional modifiers.
// Used by EntityProperties.
type EntityProperty struct {
	Key       string
	Value     float64
	Modifiers []PropertyModifier `length:"VarInt"`
}

// PropertyModifier is a modifier on a property.
// Used by EntityProperty.
type PropertyModifier struct {
	UUID      UUID `as:"raw"`
	Amount    float64
	Operation int8
}

// ChunkData sends or updates a single chunk on the client. If New is set
// then biome data should be sent too.
//
// Currently the packet id is: 0x21
type ChunkData struct {
	ChunkX, ChunkZ int32
	New            bool
	BitMask        uint16
	Data           []byte `length:"VarInt" nolimit:"true"`
}

// MultiBlockChange is used to update a batch of blocks in a single packet.
//
// Currently the packet id is: 0x22
type MultiBlockChange struct {
	ChunkX, ChunkZ int32
	Records        []BlockChangeRecord `length:"VarInt"`
}

// BlockChangeRecord is a location/id record of a block to be updated
type BlockChangeRecord struct {
	XZ      byte
	Y       byte
	BlockID VarInt
}

// BlockChange is used to update a single block on the client.
//
// Currently the packet id is: 0x23
type BlockChange struct {
	Location Position
	BlockID  VarInt
}

// BlockAction triggers different actions depending on the target block.
//
// Currently the packet id is: 0x24
type BlockAction struct {
	Location  Position
	Byte1     byte
	Byte2     byte
	BlockType VarInt
}

// BlockBreakAnimation is used to create and update the block breaking
// animation played when a player starts digging a block.
//
// Currently the packet id is: 0x25
type BlockBreakAnimation struct {
	EntityID VarInt
	Location Position
	Stage    int8
}

// ChunkDataBulk is like the ChunkData packet but allows for multiple chunks
// at once.
//
// Currently the packet id is: 0x26
type ChunkDataBulk struct {
	SkyLight bool
	Meta     []ChunkMeta `length:"VarInt"`
	Data     []byte      `length:"remaining"`
}

// ChunkMeta is used in ChunkDataBulk. See ChunkData for details.
type ChunkMeta struct {
	ChunkX, ChunkZ int32
	BitMask        uint16
}

// Explosion is sent when an explosion is triggered (tnt, creeper etc).
// This plays the effect and removes the effected blocks.
//
// Currently the packet id is: 0x27
type Explosion struct {
	X, Y, Z                         float32
	Radius                          float32
	Records                         []ExplosionRecord `length:"int32"`
	VelocityX, VelocityY, VelocityZ float32
}

// ExplosionRecord is used by explosion to mark an affected block.
type ExplosionRecord struct {
	X, Y, Z int8
}

// Effect plays a sound effect or particle at the target location with the
// volume (of sounds) being relative to the player's position unless
// DisableRelative is set to true.
//
// Currently the packet id is: 0x28
type Effect struct {
	EffectID        int32
	Location        Position
	Data            int32
	DisableRelative bool
}

// SoundEffect plays the named sound at the target location.
//
// Currently the packet id is: 0x29
type SoundEffect struct {
	Name    string
	X, Y, Z int32
	Volume  float32
	Pitch   byte
}

// Particle spawns particles at the target location with the various
// modifiers. Data's length depends on the particle ID.
//
// Currently the packet id is: 0x2A
type Particle struct {
	ParticleID                int32
	LongDistance              bool
	X, Y, Z                   float32
	OffsetX, OffsetY, OffsetZ float32
	Speed                     float32
	Count                     int32
	Data                      []VarInt `length:"@particleDataLength"`
}

func particleDataLength(p *Particle) int {
	switch p.ParticleID {
	case 36:
		return 2
	case 37, 38:
		return 1
	}
	return 0
}

// ChangeGameState is used to modify the game's state like gamemode or
// weather.
//
// Currently the packet id is: 0x2B
type ChangeGameState struct {
	Reason byte
	Value  float32
}

// SpawnGlobalEntity spawns an entity which is visible from anywhere in the
// world. Currently only used for lightning.
//
// Currently the packet id is: 0x2C
type SpawnGlobalEntity struct {
	EntityID VarInt
	Type     byte
	X, Y, Z  int32
}

// WindowOpen tells the client to open the inventory window of the given
// type. The ID is used to reference the instance of the window in
// other packets.
//
// Currently the packet id is: 0x2D
type WindowOpen struct {
	ID        byte
	Type      string
	Title     chat.AnyComponent `as:"json"`
	SlotCount byte
	EntityID  int32 `if:".Type == \"EntityHorse\""`
}

// WindowClose forces the client to close the window with the given id,
// e.g. a chest getting destroyed.
//
// Currently the packet id is: 0x2E
type WindowClose struct {
	ID byte
}

// WindowSetSlot changes an itemstack in one of the slots in a window.
//
// Currently the packet id is: 0x2F
type WindowSetSlot struct {
	ID        byte
	Slot      int16
	ItemStack ItemStack `as:"raw"`
}

// WindowItems sets every item in a window.
//
// Currently the packet id is: 0x30
type WindowItems struct {
	ID    byte
	Items []ItemStack `length:"int16" as:"raw"`
}

// WindowProperty changes the value of a property of a window. Properties
// vary depending on the window type.
//
// Currently the packet id is: 0x31
type WindowProperty struct {
	ID       byte
	Property int16
	Value    int16
}

// ConfirmTransaction notifies the client whether a transaction was successful
// or failed (e.g. due to lag).
//
// Currently the packet id is: 0x32
type ConfirmTransaction struct {
	ID           byte
	ActionNumber int16
	Accepted     bool
}

// UpdateSign sets or changes the text on a sign.
//
// Currently the packet id is: 0x33
type UpdateSign struct {
	Location Position
	Line1    chat.AnyComponent `as:"json"`
	Line2    chat.AnyComponent `as:"json"`
	Line3    chat.AnyComponent `as:"json"`
	Line4    chat.AnyComponent `as:"json"`
}

// Maps updates a single map's contents
//
// Currently the packet id is: 0x34
type Maps struct {
	ItemDamage VarInt
	Scale      int8
	Icons      []MapIcon `length:"VarInt"`
	Columns    byte
	Rows       byte   `if:".Columns>0"`
	X          byte   `if:".Columns>0"`
	Z          byte   `if:".Columns>0"`
	Data       []byte `if:".Columns>0" length:"VarInt"`
}

// MapIcon is used by Maps
type MapIcon struct {
	DirectionType int8
	X, Z          int8
}

// UpdateBlockEntity updates the nbt tag of a block entity in the
// world.
//
// Currently the packet id is: 0x35
type UpdateBlockEntity struct {
	Location Position
	Action   byte
	NBT      *nbt.Compound
}

// SignEditorOpen causes the client to open the editor for a sign so that
// it can write to it. Only sent in vanilla when the player places a sign.
//
// Currently the packet id is: 0x36
type SignEditorOpen struct {
	Location Position
}

// Statistics is used to update the statistics screen for the client.
//
// Currently the packet id is: 0x37
type Statistics struct {
	Statistics []Statistic `length:"VarInt"`
}

// Statistic is used by Statistics
type Statistic struct {
	Name  string
	Value VarInt
}

// PlayerInfo is sent by the server for every player connected to the server
// to provide skin and username information as well as ping and gamemode info.
//
// Currently the packet id is: 0x38
type PlayerInfo struct {
	Action  VarInt
	Players []PlayerDetail `length:"VarInt"`
}

// PlayerDetail is used by PlayerInfo
type PlayerDetail struct {
	UUID        UUID              `as:"raw"`
	Name        string            `if:"..Action==0"`
	Properties  []PlayerProperty  `length:"VarInt" if:"..Action==0"`
	GameMode    VarInt            `if:"..Action==0 ..Action == 1"`
	Ping        VarInt            `if:"..Action==0 ..Action == 2"`
	HasDisplay  bool              `if:"..Action==0 ..Action == 3"`
	DisplayName chat.AnyComponent `as:"json" if:".HasDisplay==true"`
}

// PlayerProperty is used by PlayerDetail
type PlayerProperty struct {
	Name      string
	Value     string
	IsSigned  bool
	Signature string `if:".IsSigned==true"`
}

// PlayerAbilities is used to modify the players current abilities. Flying,
// creative, god mode etc.
//
// Currently the packet id is: 0x39
type PlayerAbilities struct {
	Flags        byte
	FlyingSpeed  float32
	WalkingSpeed float32
}

// TabCompleteReply is sent as a reply to a tab completion request.
// The matches should be possible completions for the command/chat the
// player sent.
//
// Currently the packet id is: 0x3A
type TabCompleteReply struct {
	Matches []string `length:"VarInt"`
}

// ScoreboardObjective creates/updates a scoreboard objective.
//
// Currently the packet id is: 0x3B
type ScoreboardObjective struct {
	Name  string
	Mode  byte
	Value string `if:".Mode == 0 .Mode == 2"`
	Type  string `if:".Mode == 0 .Mode == 2"`
}

// UpdateScore is used to update or remove an item from a scoreboard
// objective.
//
// Currently the packet id is: 0x3C
type UpdateScore struct {
	Name       string
	Action     byte
	ObjectName string
	Value      VarInt `if:".Action != 1"`
}

// ScoreboardDisplay is used to set the display position of a scoreboard.
//
// Currently the packet id is: 0x3D
type ScoreboardDisplay struct {
	Position byte
	Name     string
}

// Teams creates and updates teams
//
// Currently the packet id is: 0x3E
type Teams struct {
	Name              string
	Mode              byte
	DisplayName       string   `if:".Mode == 0 .Mode == 2"`
	Prefix            string   `if:".Mode == 0 .Mode == 2"`
	Suffix            string   `if:".Mode == 0 .Mode == 2"`
	Flags             byte     `if:".Mode == 0 .Mode == 2"`
	NameTagVisibility string   `if:".Mode == 0 .Mode == 2"`
	Color             byte     `if:".Mode == 0 .Mode == 2"`
	Players           []string `length:"VarInt" if:".Mode == 0 .Mode == 3 .Mode == 4"`
}

// PluginMessageClientbound is used for custom messages between the client
// and server. This is mainly for plugins/mods but vanilla has a few channels
// registered too.
//
// Currently the packet id is: 0x3F
type PluginMessageClientbound struct {
	Channel string
	Data    []byte `length:"remaining"`
}

// Disconnect causes the client to disconnect displaying the passed reason.
//
// Currently the packet id is: 0x40
type Disconnect struct {
	Reason chat.AnyComponent `as:"json"`
}

// ServerDifficulty changes the displayed difficulty in the client's menu
// as well as some ui changes for hardcore.
//
// Currently the packet id is: 0x41
type ServerDifficulty struct {
	Difficulty byte
}

// CombatEvent is used for... you know, I never checked. I have no
// clue.
//
// Currently the packet id is: 0x42
type CombatEvent struct {
	Event    VarInt
	Duration VarInt `if:".Event == 1"`
	PlayerID VarInt `if:".Event == 2"`
	EntityID int32  `if:".Event == 1 .Event == 2"`
	Message  string `if:".Event == 2"`
}

// Camera causes the client to spectate the entity with the passed id.
// Use the player's id to de-spectate.
//
// Currently the packet id is: 0x43
type Camera struct {
	TargetID VarInt
}

// WorldBorder configures the world's border.
//
// Currently the packet id is: 0x44
type WorldBorder struct {
	Action         VarInt
	OldRadius      float64 `if:".Action == 3 .Action == 1"`
	NewRadius      float64 `if:".Action == 3 .Action == 1 .Action == 0"`
	Speed          VarLong `if:".Action == 3 .Action == 1"`
	X, Z           float64 `if:".Action == 3 .Action == 2"`
	PortalBoundary VarInt  `if:".Action == 3"`
	WarningTime    VarInt  `if:".Action == 3 .Action == 4"`
	WarningBlocks  VarInt  `if:".Action == 3 .Action == 5"`
}

// Title configures an on-screen title.
//
// Currently the packet id is: 0x45
type Title struct {
	Action   VarInt
	Title    chat.AnyComponent `as:"json" if:".Action == 0"`
	SubTitle chat.AnyComponent `as:"json" if:".Action == 1"`
	FadeIn   int32             `if:".Action == 2"`
	FadeStay int32             `if:".Action == 2"`
	FadeOut  int32             `if:".Action == 2"`
}

// SetCompression updates the compression threshold.
//
// Currently the packet id is: 0x46
type SetCompression struct {
	Threshold VarInt
}

// PlayerListHeaderFooter updates the header/footer of the player list.
//
// Currently the packet id is: 0x47
type PlayerListHeaderFooter struct {
	Header chat.AnyComponent `as:"json"`
	Footer chat.AnyComponent `as:"json"`
}

// ResourcePackSend causes the client to check its cache for the requested
// resource packet and download it if its missing. Once the resource pack
// is obtained the client will use it.
//
// Currently the packet id is: 0x48
type ResourcePackSend struct {
	URL  string
	Hash string
}

// UpdateEntityNBT updates the nbt tag for an entity.
//
// Currently the packet id is: 0x49
type UpdateEntityNBT struct {
	EntityID VarInt
	Tag      *nbt.Compound
}
