package aggregate

/*
func Test_HierarchicScopeStrategy_ReturnsTrue_IfHaystackContainsExactMatch(t *testing.T) {
	haystack := []string{"account:read:admin"}
	needle := "account:read:admin"
	if res := HierarchicScopeStrategy(haystack, needle); !res {
		t.Error("Wrongfully did not find a hierarchical scope match.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsTrue_IfHaystackContainsHierarchicalMatch(t *testing.T) {
	haystack := []string{"account:read:admin"}
	needle := "account.name:read:admin"

	if res := HierarchicScopeStrategy(haystack, needle); !res {
		t.Error("Wrongfully did not find a hierarchical scope match.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsFalse_IfNeedleFormatIsWrong(t *testing.T) {
	haystack := []string{"account:read:admin"}
	needle := "account:read:wrong"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Wrongfully found a hierarchical scope match even though needle is the wrong format.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsFalse_IfHaystackItemFormatIsWrong(t *testing.T) {
	haystack := []string{"account:read:wrong"}
	needle := "account:read:admin"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Wrongfully found a hierarchical scope match even though haystack item is the wrong format.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsFalse_IfNeedleIsLessGranular(t *testing.T) {
	haystack := []string{"account.name:read:admin"}
	needle := "account:read:admin"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Hierarchical scope match was wrongfully found even though the needle is more granular than haystack items.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsFalse_IfHaystackAndNeedleIsDifferentAction(t *testing.T) {
	haystack := []string{"account.name:write:admin"}
	needle := "account.name:read:admin"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Haystack item and needle wrongfully allowed different actions.")
	}
}

func Test_HierarchicScopeStrategy_ReturnsFalse_IfHaystackAndNeedleIsDifferentPerspective(t *testing.T) {
	haystack := []string{"account.name:write:admin"}
	needle := "account.name:read:user"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Haystack item and needle wrongfully allowed different perspectives.")
	}
}

func Test_ExactScopeStrategy_ReturnsFalse_IfHaystackItemHasDifferentSubClass(t *testing.T) {
	haystack := []string{"account.email:read:user"}
	needle := "account.name:read:user"
	if res := HierarchicScopeStrategy(haystack, needle); res {
		t.Error("Hasystack item and needle wrongfully allowed different sub classes.")
	}
}

func Test_ExactScopeStrategy_ReturnsFalse_IfHaystackitemsAndNeedleHasDifferentActions(t *testing.T) {
	haystack := []string{"account.email:read:user"}
	needle := "account.name:write:user"
	if res := ExactScopeStrategy(haystack, needle); res {
		t.Error("Wrongfully found a match even though the needle and haystack items have different actions.")
	}
}

func Test_ExactScopeStrategy_ReturnsFalse_IfNeedleAndHaystackHaveDifferentClasses(t *testing.T) {
	haystack := []string{"account.email:read:user"}
	needle := "account.name:read:user"
	if res := ExactScopeStrategy(haystack, needle); res {
		t.Error("Wrongfully found a match even though the classes of the needle and haystack item are different")
	}
}

func Test_ExactScopeStrategy_ReturnsTrue_IfTheNeedleExactlyMatchesAnHaystackItem(t *testing.T) {
	haystack := []string{"account.name:read:user"}
	needle := "account.name:read:user"
	if res := ExactScopeStrategy(haystack, needle); !res {
		t.Error("Did not find exact match even though there is one")
	}
}

func Test_ParseScope_ReturnsCorrectScope_IfScopeIsValid(t *testing.T) {
	scope := "account.name:read:admin"
	result, err := ParseScope(scope)
	if err != nil {
		t.Error("Correct scope string resulted in error")
	}

	if len(result.Classes) != 2 ||
		result.Classes[0] != "account" ||
		result.Classes[1] != "name" {
		t.Error("Parsed scope does not contain the correct classes")
	}

	if result.Action != "read" {
		t.Error("Parsed scope action is not 'read'")
	}

	if result.Perspective != "admin" {
		t.Error("Parsed scope perspective is not 'admin'")
	}
}

func Test_ParseScope_ReturnsErrNoClasses_IfNoClass(t *testing.T) {
	scope := ":wrong:admin"
	if _, err := ParseScope(scope); err != ErrNoClasses {
		t.Error("Scope parsing did not result in 'ErrNoClasses'.")
	}
}

func Test_ParseScope_ReturnsErrInvalidAction_IfInvalidAction(t *testing.T) {
	scope := "account.name:wrong:admin"
	if _, err := ParseScope(scope); err != ErrInvalidAction {
		t.Error("Scope parsing did not result in 'ErrInvalidAction'.")
	}
}

func Test_ParseScope_ReturnsErrInvalidScopeString_IfTooManyColons(t *testing.T) {
	scope := "account.name:read:admin:admin"
	if _, err := ParseScope(scope); err != ErrTooManyColons {
		t.Error("Scope parsing did not result in 'ErrTooManyColons'.")
	}
}

func Test_ParseScope_ReturnsErrDoesNotContainPerspective_IfDoesNotContainPerspective(t *testing.T) {
	scope := "account.name:read"
	if _, err := ParseScope(scope); err != ErrDoesNotContainPerspective {
		t.Error("Scope parsing did not result in 'ErrDoesNotContainPerspective'.")
	}
}

func Test_ParseScope_ReturnsErrInvalidAction_IfInvalidPerspective(t *testing.T) {
	scope := "account.name:read:wrong"
	if _, err := ParseScope(scope); err != ErrInvalidPerspectiove {
		t.Error("Scope parsing did not result in 'ErrInvalidPerspectiove'.")
	}
}
*/
