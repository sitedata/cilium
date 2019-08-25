package endpoint

import (
	"time"

	"github.com/cilium/cilium/common/addressing"
	"github.com/cilium/cilium/pkg/endpoint/regeneration"
	"github.com/cilium/cilium/pkg/identity"

	. "gopkg.in/check.v1"
)

// PrepareEndpointForTesting creates an endpoint useful for testing purposes.
func PrepareEndpointForTesting(owner regeneration.Owner, id uint16, identity *identity.Identity, ipv4 addressing.CiliumIPv4, ipv6 addressing.CiliumIPv6) *Endpoint {
	e := NewEndpointWithState(owner, id, StateWaitingForIdentity)
	e.IPv6 = ipv6
	e.IPv4 = ipv4
	e.SetIdentity(identity, true)

	e.unconditionalLock()
	e.setState(StateWaitingToRegenerate, "test")
	e.unlock()
	return e
}

func (e *Endpoint) RegenerateEndpointTest(c *C, regenMetadata *regeneration.ExternalRegenerationMetadata) {
	e.unconditionalLock()
	ready := e.setState(StateWaitingToRegenerate, "test")
	e.unlock()
	c.Assert(ready, Equals, true)
	buildSuccess := <-e.Regenerate(regenMetadata)
	c.Assert(buildSuccess, Equals, true)
}

func (e *Endpoint) WaitForIdentity(timeoutDuration time.Duration) *identity.Identity {
	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()
	var secID *identity.Identity
	for {
		select {
		case <-timeout.C:
			return nil
		case <-tick.C:
			e.unconditionalRLock()
			secID = e.SecurityIdentity
			e.RUnlock()
			if secID != nil {
				return secID
			}
		}
	}
}
