package endpoint

import (
	"github.com/cilium/cilium/common/addressing"
	"github.com/cilium/cilium/pkg/endpoint/regeneration"
	"github.com/cilium/cilium/pkg/identity"
)

// PrepareEndpointForTesting creates an endpoint useful for testing purposes.
func PrepareEndpointForTesting(owner regeneration.Owner, id uint16, identity *identity.Identity, ipv4 addressing.CiliumIPv4, ipv6 addressing.CiliumIPv6) *Endpoint {
	e := NewEndpointWithState(owner, id, StateWaitingForIdentity)
	e.IPv6 = ipv6
	e.IPv4 = ipv4
	e.SetIdentity(identity, true)

	e.UnconditionalLock()
	e.SetStateLocked(StateWaitingToRegenerate, "test")
	e.Unlock()
	return e
}
