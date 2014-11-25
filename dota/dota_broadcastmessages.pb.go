// Code generated by protoc-gen-go.
// source: dota_broadcastmessages.proto
// DO NOT EDIT!

package dota

import proto "github.com/golang/protobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import google_protobuf "github.com/dotabuff/yasha/dota/google/protobuf/descriptor.pb"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type EDotaBroadcastMessages int32

const (
	EDotaBroadcastMessages_DOTA_BM_LANLobbyRequest EDotaBroadcastMessages = 1
	EDotaBroadcastMessages_DOTA_BM_LANLobbyReply   EDotaBroadcastMessages = 2
)

var EDotaBroadcastMessages_name = map[int32]string{
	1: "DOTA_BM_LANLobbyRequest",
	2: "DOTA_BM_LANLobbyReply",
}
var EDotaBroadcastMessages_value = map[string]int32{
	"DOTA_BM_LANLobbyRequest": 1,
	"DOTA_BM_LANLobbyReply":   2,
}

func (x EDotaBroadcastMessages) Enum() *EDotaBroadcastMessages {
	p := new(EDotaBroadcastMessages)
	*p = x
	return p
}
func (x EDotaBroadcastMessages) String() string {
	return proto.EnumName(EDotaBroadcastMessages_name, int32(x))
}
func (x *EDotaBroadcastMessages) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EDotaBroadcastMessages_value, data, "EDotaBroadcastMessages")
	if err != nil {
		return err
	}
	*x = EDotaBroadcastMessages(value)
	return nil
}

type CDOTABroadcastMsg struct {
	Type             *EDotaBroadcastMessages `protobuf:"varint,1,req,name=type,enum=dota.EDotaBroadcastMessages,def=1" json:"type,omitempty"`
	Msg              []byte                  `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *CDOTABroadcastMsg) Reset()         { *m = CDOTABroadcastMsg{} }
func (m *CDOTABroadcastMsg) String() string { return proto.CompactTextString(m) }
func (*CDOTABroadcastMsg) ProtoMessage()    {}

const Default_CDOTABroadcastMsg_Type EDotaBroadcastMessages = EDotaBroadcastMessages_DOTA_BM_LANLobbyRequest

func (m *CDOTABroadcastMsg) GetType() EDotaBroadcastMessages {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_CDOTABroadcastMsg_Type
}

func (m *CDOTABroadcastMsg) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

type CDOTABroadcastMsg_LANLobbyRequest struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *CDOTABroadcastMsg_LANLobbyRequest) Reset()         { *m = CDOTABroadcastMsg_LANLobbyRequest{} }
func (m *CDOTABroadcastMsg_LANLobbyRequest) String() string { return proto.CompactTextString(m) }
func (*CDOTABroadcastMsg_LANLobbyRequest) ProtoMessage()    {}

type CDOTABroadcastMsg_LANLobbyReply struct {
	Id               *uint64                                         `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	TournamentId     *uint32                                         `protobuf:"varint,2,opt,name=tournament_id" json:"tournament_id,omitempty"`
	TournamentGameId *uint32                                         `protobuf:"varint,3,opt,name=tournament_game_id" json:"tournament_game_id,omitempty"`
	Members          []*CDOTABroadcastMsg_LANLobbyReply_CLobbyMember `protobuf:"bytes,4,rep,name=members" json:"members,omitempty"`
	RequiresPassKey  *bool                                           `protobuf:"varint,5,opt,name=requires_pass_key" json:"requires_pass_key,omitempty"`
	LeaderAccountId  *uint32                                         `protobuf:"varint,6,opt,name=leader_account_id" json:"leader_account_id,omitempty"`
	XXX_unrecognized []byte                                          `json:"-"`
}

func (m *CDOTABroadcastMsg_LANLobbyReply) Reset()         { *m = CDOTABroadcastMsg_LANLobbyReply{} }
func (m *CDOTABroadcastMsg_LANLobbyReply) String() string { return proto.CompactTextString(m) }
func (*CDOTABroadcastMsg_LANLobbyReply) ProtoMessage()    {}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetTournamentId() uint32 {
	if m != nil && m.TournamentId != nil {
		return *m.TournamentId
	}
	return 0
}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetTournamentGameId() uint32 {
	if m != nil && m.TournamentGameId != nil {
		return *m.TournamentGameId
	}
	return 0
}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetMembers() []*CDOTABroadcastMsg_LANLobbyReply_CLobbyMember {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetRequiresPassKey() bool {
	if m != nil && m.RequiresPassKey != nil {
		return *m.RequiresPassKey
	}
	return false
}

func (m *CDOTABroadcastMsg_LANLobbyReply) GetLeaderAccountId() uint32 {
	if m != nil && m.LeaderAccountId != nil {
		return *m.LeaderAccountId
	}
	return 0
}

type CDOTABroadcastMsg_LANLobbyReply_CLobbyMember struct {
	AccountId        *uint32 `protobuf:"varint,1,opt,name=account_id" json:"account_id,omitempty"`
	PlayerName       *string `protobuf:"bytes,2,opt,name=player_name" json:"player_name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CDOTABroadcastMsg_LANLobbyReply_CLobbyMember) Reset() {
	*m = CDOTABroadcastMsg_LANLobbyReply_CLobbyMember{}
}
func (m *CDOTABroadcastMsg_LANLobbyReply_CLobbyMember) String() string {
	return proto.CompactTextString(m)
}
func (*CDOTABroadcastMsg_LANLobbyReply_CLobbyMember) ProtoMessage() {}

func (m *CDOTABroadcastMsg_LANLobbyReply_CLobbyMember) GetAccountId() uint32 {
	if m != nil && m.AccountId != nil {
		return *m.AccountId
	}
	return 0
}

func (m *CDOTABroadcastMsg_LANLobbyReply_CLobbyMember) GetPlayerName() string {
	if m != nil && m.PlayerName != nil {
		return *m.PlayerName
	}
	return ""
}

func init() {
	proto.RegisterEnum("dota.EDotaBroadcastMessages", EDotaBroadcastMessages_name, EDotaBroadcastMessages_value)
}
