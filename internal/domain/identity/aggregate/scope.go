package aggregate

import (
	"errors"
	"strings"

	"github.com/hywmongous/example-service/internal/domain/identity/values"
)

// order -> Can read evyerthing regarding the order
// order:items -> Can read items in order
// order:items.write -> Can write items to order
// https://datatracker.ietf.org/doc/html/rfc6749#section-3.3

/* Explanation
 * haystack: The scopes of the user
 * needle: The required scope
 *
 * Grammar:
 * Scope       = Class {"." Class} ":" Action ":" Perspective
 * Class       = String
 * Action      = read | write
 * Perspective = user | admin
 *
 * Example input
 * haystack: [ "order.items:write", "basket.items:write" ]
 * needle: [ "order.items:write" ]
 * HierarchicScopeStrategy: True
 * ExactScopeStrategy: False
 * */

type Scope struct {
	id          values.ScopeID
	classes     []string
	action      string
	perspective string
}

func CreateScope(
	classes []string,
	action string,
	perspective string,
) (Scope, error) {
	CreateScope := Scope{
		values.GenerateScopeID(),
		classes,
		action,
		perspective,
	}
	return CreateScope, nil
}

func RecreateScope(
	id values.ScopeID,
	classes []string,
	action string,
	perspective string,
) Scope {
	return Scope{
		id,
		classes,
		action,
		perspective,
	}
}

func (scope Scope) ToString() string {
	var result string = scope.classes[0]
	for i := 1; i < len(scope.classes); i++ {
		result += "." + scope.classes[i]
	}
	result += ":" + scope.action
	result += ":" + scope.perspective
	return result
}

type ScopeStrategy func(haystack []string, needle string) bool
type ScopeMatch func(hay Scope, needle Scope) bool

const (
	ReadAction       = "read"
	WriteAction      = "write"
	UserPerspective  = "user"
	AdminPerspective = "admin"
)

var (
	ErrInvalidAction             = errors.New("the action of the scope is not supported")
	ErrInvalidPerspectiove       = errors.New("the perspective of the scope is invalid")
	ErrTooManyColons             = errors.New("the amount of colons are greater than two")
	ErrNoClasses                 = errors.New("no classes were found in the scope")
	ErrDoesNotContainPerspective = errors.New("the scope string does not contain a perspective")
)

func HierarchicScopeStrategy(haystack []Scope, needle Scope) bool {
	return strategyDriver(haystack, needle, HierarchicMatch)
}

func HierarchicMatch(hay Scope, needle Scope) bool {
	if hay.action != needle.action ||
		hay.perspective != needle.perspective {
		return false
	}

	// hay: "order" and needle: "order.items" then true
	// hay: "order.items" and needle: "order" then false
	// Obs action is in the example omitted for clarity

	// If the classes for the hay is greater then it is more granular
	// the result of this is that the needle will never be found
	// Eg. you are allow "order.items" but require read access to "order"
	// and not just the items found within the order
	if len(hay.classes) > len(needle.classes) {
		return false
	}

	// Since the precondition is that hay.classes <= needled.classes
	// then no out of bounds errors will occure. Also since the body
	// only performance computations without permutating the arrays
	// then the invariant (the same as the precondition) is always kept
	for idx, class := range hay.classes {
		if class != needle.classes[idx] {
			return false
		}
	}

	return true
}

func ExactScopeStrategy(haystack []Scope, needle Scope) bool {
	return strategyDriver(haystack, needle, ExactMatch)
}

func ExactMatch(hay Scope, needle Scope) bool {
	if hay.action != needle.action ||
		hay.perspective != needle.perspective ||
		len(hay.classes) != len(needle.classes) {
		return false
	}

	// len of classes for both the haystack and needle
	// are the same based on the verified precondition
	for idx, class := range hay.classes {
		if class != needle.classes[idx] {
			return false
		}
	}
	return true
}

func strategyDriver(haystack []Scope, needle Scope, comparer func(Scope, Scope) bool) bool {
	for _, curr := range haystack {
		if comparer(curr, needle) {
			return true
		}
	}
	return false
}

func ParseScope(scope string) (Scope, error) {
	dotSplit := strings.Split(scope, ".")
	dotSplitLen := len(dotSplit)

	colonSplit := strings.Split(dotSplit[dotSplitLen-1], ":")

	dotSplit[dotSplitLen-1] = colonSplit[0]
	classes := dotSplit
	if classes[0] == "" {
		return Scope{}, ErrNoClasses
	}

	action := colonSplit[1]
	if action != ReadAction &&
		action != WriteAction {
		return Scope{}, ErrInvalidAction
	}

	perspective := ""
	if len(colonSplit) > 3 {
		return Scope{}, ErrTooManyColons
	}
	if len(colonSplit) == 2 {
		return Scope{}, ErrDoesNotContainPerspective
	}

	perspective = colonSplit[2]
	if perspective != UserPerspective &&
		perspective != AdminPerspective {
		return Scope{}, ErrInvalidPerspectiove
	}

	return Scope{
		classes:     classes,
		action:      action,
		perspective: perspective,
	}, nil
}
