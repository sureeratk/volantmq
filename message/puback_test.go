// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
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

package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubAckMessageFields(t *testing.T) {
	m, err := NewMessage(ProtocolV311, PUBACK)
	require.NoError(t, err)

	msg, ok := m.(*AckMessage)
	require.True(t, ok, "Couldn't cast message type")

	msg.SetPacketID(100)

	id, _ := msg.PacketID()
	require.Equal(t, PacketID(100), id)
}

func TestPubAckMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	m, n, err := Decode(ProtocolV311, msgBytes)
	msg, ok := m.(*AckMessage)
	require.Equal(t, true, ok, "Invalid message type")

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "decode length does not match")
	require.Equal(t, PUBACK, msg.Type(), "Message type does not match")

	id, _ := msg.PacketID()
	require.Equal(t, PacketID(7), id, "PacketID does not match")
}

// test insufficient bytes
func TestPubAckMessageDecode2(t *testing.T) {
	msgBytes := []byte{
		byte(PUBACK << 4),
		2,
		7, // packet ID LSB (7)
	}

	_, _, err := Decode(ProtocolV311, msgBytes)

	require.Error(t, err)
}

func TestPubAckMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	m, err := NewMessage(ProtocolV311, PUBACK)
	require.NoError(t, err)

	msg, ok := m.(*AckMessage)
	require.True(t, ok, "Couldn't cast message type")

	msg.SetPacketID(7)

	dst := make([]byte, 10)
	n, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n], "Error decoding message.")
}

// test to ensure encoding and decoding are the same
// decode, encode, and decode again
func TestPubAckDecodeEncodeEquiv(t *testing.T) {
	msgBytes := []byte{
		byte(PUBACK << 4),
		2,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
	}

	m, n, err := Decode(ProtocolV311, msgBytes)
	msg, ok := m.(*AckMessage)
	require.Equal(t, true, ok, "Invalid message type")

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n, "Error decoding message.")

	dst := make([]byte, 100)
	n2, err := msg.Encode(dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n2, "Error decoding message.")
	require.Equal(t, msgBytes, dst[:n2], "Error decoding message.")

	_, n3, err := Decode(ProtocolV311, dst)

	require.NoError(t, err, "Error decoding message.")
	require.Equal(t, len(msgBytes), n3, "Error decoding message.")
}
