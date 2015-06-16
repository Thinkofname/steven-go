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

//go:generate protocol_builder $GOFILE Play serverbound

package protocol

import (
	"github.com/thinkofdeath/steven/format"
)

// KeepAliveServerbound is sent by a client as a response to a
// KeepAliveClientbound. If the client doesn't reply the server
// may disconnect the client.
//
// Currently the packet id is: 0x00
type KeepAliveServerbound struct {
	ID VarInt
}

// ChatMessage is sent by the client when it sends a chat message or
// executes a command (prefixed by '/').
//
// Currently the packet id is: 0x01
type ChatMessage struct {
	Message string
}

// UseEntity is sent when the user interacts (right clicks) or attacks
// (left clicks) an entity.
//
// Currently the packet id is: 0x02
type UseEntity struct {
	TargetID VarInt
	Type     VarInt
	TargetX  float32 `if:".Type==2"`
	TargetY  float32 `if:".Type==2"`
	TargetZ  float32 `if:".Type==2"`
}

// Player is used to update whether the player is on the ground or not.
//
// Currently the packet id is: 0x03
type Player struct {
	OnGround bool
}

// PlayerPosition is used to update the player's position.
//
// Currently the packet id is: 0x04
type PlayerPosition struct {
	X, Y, Z  float64
	OnGround bool
}

// PlayerLook is used to update the player's rotation.
//
// Currently the packet id is: 0x05
type PlayerLook struct {
	Yaw, Pitch float32
	OnGround   bool
}

// PlayerPositionLook is a combination of PlayerPosition and
// PlayerLook.
//
// Currently the packet id is: 0x06
type PlayerPositionLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool
}

// PlayerDigging is sent when the client starts/stops digging a block.
// It also can be sent for droppping items and eating/shooting.
//
// Currently the packet id is: 0x07
type PlayerDigging struct {
	Status   byte
	Location Position
	Face     byte
}

// PlayerBlockPlacement is sent when the client tries to place a block.
//
// Currently the packet id is: 0x08
type PlayerBlockPlacement struct {
	Location                  Position
	Face                      byte
	HeldItem                  ItemStack `as:"raw"`
	CursorX, CursorY, CursorZ byte
}

// HeldItemChange is sent when the player changes the currently active
// hotbar slot.
//
// Currently the packet id is: 0x09
type HeldItemChange struct {
	Slot int16
}

// ArmSwing is sent by the client when the player left clicks (to swing their
// arm).
//
// Currently the packet id is: 0x0A
type ArmSwing struct {
}

// PlayerAction is sent when a player preforms various actions.
//
// Currently the packet id is: 0x0B
type PlayerAction struct {
	EntityID  VarInt
	ActionID  VarInt
	JumpBoost VarInt
}

// SteerVehicle is sent by the client when steers or preforms an action
// on a vehicle.
//
// Currently the packet id is: 0x0C
type SteerVehicle struct {
	Sideways float32
	Forward  float32
	Flags    byte
}

// CloseWindow is sent when the client closes a window.
//
// Currently the packet id is: 0x0D
type CloseWindow struct {
	ID byte
}

// ClickWindow is sent when the client clicks in a window.
//
// Currently the packet id is: 0x0E
type ClickWindow struct {
	ID           byte
	Slot         int16
	Button       byte
	ActionNumber int16
	Mode         byte
	ClickedItem  ItemStack `as:"raw"`
}

// ConfirmTransactionServerbound is a reply to ConfirmTransaction.
//
// Currently the packet id is: 0x0F
type ConfirmTransactionServerbound struct {
	ID           byte
	ActionNumber int16
	Accepted     bool
}

// CreativeInventoryAction is sent when the client clicks in the creative
// inventory. This is used to spawn items in creative.
//
// Currently the packet id is: 0x10
type CreativeInventoryAction struct {
	Slot        int16
	ClickedItem ItemStack `as:"raw"`
}

// EnchantItem is sent when the client enchants an item.
//
// Currently the packet id is: 0x11
type EnchantItem struct {
	ID          byte
	Enchantment byte
}

// SetSign sets the text on a sign after placing it.
//
// Currently the packet id is: 0x12
type SetSign struct {
	Location Position
	Line1    format.AnyComponent `as:"json"`
	Line2    format.AnyComponent `as:"json"`
	Line3    format.AnyComponent `as:"json"`
	Line4    format.AnyComponent `as:"json"`
}

// ClientAbilities is used to modify the players current abilities.
// Currently flying is the only one
//
// Currently the packet id is: 0x13
type ClientAbilities struct {
	Flags        byte
	FlyingSpeed  float32
	WalkingSpeed float32
}

// TabComplete is sent by the client when the client presses tab in
// the chat box.
//
// Currently the packet id is: 0x14
type TabComplete struct {
	Text      string
	HasTarget bool
	Target    Position `if:".HasTarget==true"`
}

// ClientSettings is sent by the client to update its current settings.
//
// Currently the packet id is: 0x15
type ClientSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           byte
	ChatColors         bool
	DisplayedSkinParts byte
}

// ClientStatus is sent to update the client's status
//
// Currently the packet id is: 0x16
type ClientStatus struct {
	ActionID VarInt
}

// PluginMessageServerbound is used for custom messages between the client
// and server. This is mainly for plugins/mods but vanilla has a few channels
// registered too.
//
// Currently the packet id is: 0x17
type PluginMessageServerbound struct {
	Channel string
	Data    []byte `length:"remaining"`
}

// SpectateTeleport is sent by clients in spectator mode to teleport to a player.
//
// Currently the packet id is: 0x18
type SpectateTeleport struct {
	Target UUID `as:"raw"`
}

// ResourcePackStatus informs the server of the client's current progress
// in activating the requested resource pack
type ResourcePackStatus struct {
	Hash   string
	Result VarInt
}
