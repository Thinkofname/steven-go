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
// This is a Minecraft packet
type KeepAliveServerbound struct {
	ID VarInt
}

// ChatMessage is sent by the client when it sends a chat message or
// executes a command (prefixed by '/').
//
// This is a Minecraft packet
type ChatMessage struct {
	Message string
}

// UseEntity is sent when the user interacts (right clicks) or attacks
// (left clicks) an entity.
//
// This is a Minecraft packet
type UseEntity struct {
	TargetID VarInt
	Type     VarInt
	TargetX  float32 `if:".Type==2"`
	TargetY  float32 `if:".Type==2"`
	TargetZ  float32 `if:".Type==2"`
}

// Player is used to update whether the player is on the ground or not.
//
// This is a Minecraft packet
type Player struct {
	OnGround bool
}

// PlayerPosition is used to update the player's position.
//
// This is a Minecraft packet
type PlayerPosition struct {
	X, Y, Z  float64
	OnGround bool
}

// PlayerLook is used to update the player's rotation.
//
// This is a Minecraft packet
type PlayerLook struct {
	Yaw, Pitch float32
	OnGround   bool
}

// PlayerPositionLook is a combination of PlayerPosition and
// PlayerLook.
//
// This is a Minecraft packet
type PlayerPositionLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool
}

// PlayerDigging is sent when the client starts/stops digging a block.
// It also can be sent for droppping items and eating/shooting.
//
// This is a Minecraft packet
type PlayerDigging struct {
	Status   byte
	Location Position
	Face     byte
}

// PlayerBlockPlacement is sent when the client tries to place a block.
//
// This is a Minecraft packet
type PlayerBlockPlacement struct {
	Location                  Position
	Face                      byte
	HeldItem                  ItemStack `as:"raw"`
	CursorX, CursorY, CursorZ byte
}

// HeldItemChange is sent when the player changes the currently active
// hotbar slot.
//
// This is a Minecraft packet
type HeldItemChange struct {
	Slot int16
}

// ArmSwing is sent by the client when the player left clicks (to swing their
// arm).
//
// This is a Minecraft packetA
type ArmSwing struct {
}

// PlayerAction is sent when a player preforms various actions.
//
// This is a Minecraft packetB
type PlayerAction struct {
	EntityID  VarInt
	ActionID  VarInt
	JumpBoost VarInt
}

// SteerVehicle is sent by the client when steers or preforms an action
// on a vehicle.
//
// This is a Minecraft packetC
type SteerVehicle struct {
	Sideways float32
	Forward  float32
	Flags    byte
}

// CloseWindow is sent when the client closes a window.
//
// This is a Minecraft packetD
type CloseWindow struct {
	ID byte
}

// ClickWindow is sent when the client clicks in a window.
//
// This is a Minecraft packetE
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
// This is a Minecraft packetF
type ConfirmTransactionServerbound struct {
	ID           byte
	ActionNumber int16
	Accepted     bool
}

// CreativeInventoryAction is sent when the client clicks in the creative
// inventory. This is used to spawn items in creative.
//
// This is a Minecraft packet
type CreativeInventoryAction struct {
	Slot        int16
	ClickedItem ItemStack `as:"raw"`
}

// EnchantItem is sent when the client enchants an item.
//
// This is a Minecraft packet
type EnchantItem struct {
	ID          byte
	Enchantment byte
}

// SetSign sets the text on a sign after placing it.
//
// This is a Minecraft packet
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
// This is a Minecraft packet
type ClientAbilities struct {
	Flags        byte
	FlyingSpeed  float32
	WalkingSpeed float32
}

// TabComplete is sent by the client when the client presses tab in
// the chat box.
//
// This is a Minecraft packet
type TabComplete struct {
	Text      string
	HasTarget bool
	Target    Position `if:".HasTarget==true"`
}

// ClientSettings is sent by the client to update its current settings.
//
// This is a Minecraft packet
type ClientSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           byte
	ChatColors         bool
	DisplayedSkinParts byte
}

// ClientStatus is sent to update the client's status
//
// This is a Minecraft packet
type ClientStatus struct {
	ActionID VarInt
}

// PluginMessageServerbound is used for custom messages between the client
// and server. This is mainly for plugins/mods but vanilla has a few channels
// registered too.
//
// This is a Minecraft packet
type PluginMessageServerbound struct {
	Channel string
	Data    []byte `length:"remaining"`
}

// SpectateTeleport is sent by clients in spectator mode to teleport to a player.
//
// This is a Minecraft packet
type SpectateTeleport struct {
	Target UUID `as:"raw"`
}

// ResourcePackStatus informs the server of the client's current progress
// in activating the requested resource pack
type ResourcePackStatus struct {
	Hash   string
	Result VarInt
}
