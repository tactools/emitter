/**********************************************************************************
* Copyright (c) 2009-2017 Misakai Ltd.
* This program is free software: you can redistribute it and/or modify it under the
* terms of the GNU Affero General Public License as published by the  Free Software
* Foundation, either version 3 of the License, or(at your option) any later version.
*
* This program is distributed  in the hope that it  will be useful, but WITHOUT ANY
* WARRANTY;  without even  the implied warranty of MERCHANTABILITY or FITNESS FOR A
* PARTICULAR PURPOSE.  See the GNU Affero General Public License  for  more details.
*
* You should have  received a copy  of the  GNU Affero General Public License along
* with this program. If not, see<http://www.gnu.org/licenses/>.
************************************************************************************/

package subscription

import (
	"sync"

	"github.com/emitter-io/emitter/security"
)

// Ssid represents a subscription ID which contains a contract and a list of hashes
// for various parts of the channel.
type Ssid []uint32

// NewSsid creates a new Ssid.
func NewSsid(contract uint32, c *security.Channel) Ssid {
	ssid := make([]uint32, 0, len(c.Query)+1)
	ssid = append(ssid, uint32(contract))
	ssid = append(ssid, c.Query...)
	return ssid
}

// Contract gets the contract part from Ssid.
func (s Ssid) Contract() uint32 {
	return uint32(s[0])
}

// GetHashCode combines the ssid into a single hash.
func (s Ssid) GetHashCode() uint32 {
	h := s[0]
	for _, i := range s[1:] {
		h ^= i
	}
	return h
}

// ------------------------------------------------------------------------------------

// SubscriberType represents a type of subscriber
type SubscriberType uint8

// Subscriber types
const (
	SubscriberDirect = SubscriberType(iota)
	SubscriberRemote
)

// Subscriber is a value associated with a subscription.
type Subscriber interface {
	ID() string
	Type() SubscriberType
	Send(ssid Ssid, channel []byte, payload []byte) error
}

// ------------------------------------------------------------------------------------

// Subscribers represents a subscriber set which can contain only unique values.
type Subscribers []Subscriber

// AddUnique adds a subscriber to the set.
func (s *Subscribers) AddUnique(value Subscriber) {
	if s.Contains(value) == false {
		*s = append(*s, value)
	}
}

// Contains checks whether a subscriber is in the set.
func (s *Subscribers) Contains(value Subscriber) bool {
	for _, v := range *s {
		if v == value {
			return true
		}
	}
	return false
}

// Subscription represents a topic subscription.
type Subscription struct {
	Ssid       Ssid       // Gets or sets the SSID (parsed channel) for this subscription.
	Subscriber Subscriber // Gets or sets the subscriber for this subscription.
}

// ------------------------------------------------------------------------------------

// Set represents a set of subscriptions with the channel names. This
// is really only useful for automatic unsubscription when a Conn/Peer disconnects.
type Set sync.Map

// Add adds a subscription to the set.
func (s *Set) Add(channel string, ssid Ssid) {
	(*sync.Map)(s).Store(channel, ssid)
}

// Remove removes a subscription from the set.
func (s *Set) Remove(channel string) {
	(*sync.Map)(s).Delete(channel)
}

// Range iterates through all of the
func (s *Set) Range(f func(string, Ssid) bool) {
	(*sync.Map)(s).Range(func(k, v interface{}) bool {
		return f(k.(string), v.(Ssid))
	})
}
