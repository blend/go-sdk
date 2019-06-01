package migration

import "github.com/blend/go-sdk/logger"

// SuiteOption is an option for migration Suites
type SuiteOption func(s *Suite)

// OptGroups allows you to add groups to the Suite. If you want, multiple OptGroups can be applied to the same Suite.
// They are additive.
func OptGroups(groups ...*Group) SuiteOption {
	return func(s *Suite) {
		if len(s.Groups) == 0 {
			s.Groups = groups
		} else {
			s.Groups = append(s.Groups, groups...)
		}

	}
}

// OptLog allows you to add a logger to the Suite.
func OptLog(log logger.Log) SuiteOption {
	return func(s *Suite) {
		s.Log = log
	}
}

// GroupOption is an option for migration Groups (Group)
type GroupOption func(g *Group)

// OptActions allows you to add actions to the NewGroup. If you want, multiple OptActions can be applied to the same NewGroup.
// They are additive.
func OptActions(actions ...Actionable) GroupOption {
	return func(g *Group) {
		if len(g.Actions) == 0 {
			g.Actions = actions
		} else {
			g.Actions = append(g.Actions, actions...)
		}
	}
}

// OptSkipTransaction will allow this group to be run outside of a transaction. Use this to concurrently create indices
// and perform other actions that cannot be executed in a Tx
func OptSkipTransaction() GroupOption {
	return func(g *Group) {
		g.SkipTransaction = true
	}
}
