package goclash

import "sync"

type LiteralKeysProvider struct {
	keys    []string
	nextKey int
	m       sync.Mutex
}

func (p *LiteralKeysProvider) GetKey() string {
	p.m.Lock()
	key := p.keys[p.nextKey]
	p.nextKey = (p.nextKey + 1) % len(p.keys)
	p.m.Unlock()
	return key
}

func (*LiteralKeysProvider) RevalidateKeys() error {
	return nil // Nothing to do
}
