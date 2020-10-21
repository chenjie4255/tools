package event

import "sync"

// FIXME: this is a temporary implement, refactor me ASAP

// NotifyAndWait publish one event and wait until all listener finished their job
func NotifyAndWait(group, name string, data interface{}) error {
	mgr.rwLock.RLock()
	defer mgr.rwLock.RUnlock()

	if listeners, found := mgr.listenerGroups[group]; found {
		for i := range listeners {
			if err := listeners[i](name, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// Propose publish one event and wait until all proposal reviewers' opinion
func Propose(group, name string, data interface{}) error {
	mgr.rwLock.RLock()
	defer mgr.rwLock.RUnlock()

	if rejecters, found := mgr.rejecterGroups[group]; found {
		for i := range rejecters {
			if err := rejecters[i](name, data); err != nil {
				return err
			}
		}
	}

	return nil
}

type ProposalEventCB func(name string, data interface{}) error
type ListenerEventCB func(name string, data interface{}) error

func RegisterListener(group string, cb ListenerEventCB) {
	if group == "" {
		panic("listener group cannot be empty")
	}

	mgr.rwLock.Lock()
	defer mgr.rwLock.Unlock()

	if listeners, found := mgr.listenerGroups[group]; found {
		listeners = append(listeners, cb)
		mgr.listenerGroups[group] = listeners
	} else {
		mgr.listenerGroups[group] = []ListenerEventCB{cb}
	}
}

func RegisterProposalReviewer(group string, cb ProposalEventCB) {
	if group == "" {
		panic("listener group cannot be empty")
	}

	mgr.rwLock.Lock()
	defer mgr.rwLock.Unlock()

	if list, found := mgr.rejecterGroups[group]; found {
		list = append(list, cb)
		mgr.rejecterGroups[group] = list
	} else {
		mgr.rejecterGroups[group] = []ProposalEventCB{cb}
	}
}

type manager struct {
	rwLock         sync.RWMutex
	listenerGroups map[string][]ListenerEventCB
	rejecterGroups map[string][]ProposalEventCB
}

var mgr manager

func init() {
	mgr.rwLock = sync.RWMutex{}
	mgr.listenerGroups = map[string][]ListenerEventCB{}
	mgr.rejecterGroups = map[string][]ProposalEventCB{}
}
