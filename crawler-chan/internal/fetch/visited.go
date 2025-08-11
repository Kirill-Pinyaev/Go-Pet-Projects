// internal/fetch/visited.go
package fetch

import "sync"

type Visited struct{ m sync.Map }

func (v *Visited) Seen(raw string) bool {
	_, existed := v.m.LoadOrStore(raw, struct{}{})
	return existed // true, если URL уже был
}
