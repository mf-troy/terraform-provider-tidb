package mysql

import "testing"

func TestNormalizeRoleIdentifier(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bare role name stays bare",
			input:    "readonly_role",
			expected: "readonly_role",
		},
		{
			name:     "quoted percent-host role normalizes to bare name",
			input:    "'readonly_role'@'%'",
			expected: "readonly_role",
		},
		{
			name:     "explicit percent-host role normalizes to bare name",
			input:    "readonly_role@%",
			expected: "readonly_role",
		},
		{
			name:     "non-default host is preserved",
			input:    "'readonly_role'@'internal'",
			expected: "readonly_role@internal",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := normalizeRoleIdentifier(testCase.input)
			if actual != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, actual)
			}
		})
	}
}

func TestRoleIdentifierSQL(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bare role name gets percent host",
			input:    "readonly_role",
			expected: "'readonly_role'@'%'",
		},
		{
			name:     "legacy quoted role stays canonical",
			input:    "'readonly_role'@'%'",
			expected: "'readonly_role'@'%'",
		},
		{
			name:     "non-default host is preserved",
			input:    "readonly_role@internal",
			expected: "'readonly_role'@'internal'",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := roleIdentifierSQL(testCase.input)
			if actual != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, actual)
			}
		})
	}
}

func TestRoleGrantSQLStatements(t *testing.T) {
	grant := &RoleGrant{
		Roles: []string{
			"readonly_role",
			"'writer_role'@'%'",
		},
		UserOrRole: UserOrRole{
			Name: "app_user",
			Host: "%",
		},
	}

	expectedGrant := "GRANT 'readonly_role'@'%', 'writer_role'@'%' TO 'app_user'@'%'"
	expectedRevoke := "REVOKE 'readonly_role'@'%', 'writer_role'@'%' FROM 'app_user'@'%'"

	if actual := grant.SQLGrantStatement(); actual != expectedGrant {
		t.Fatalf("expected grant SQL %q, got %q", expectedGrant, actual)
	}

	if actual := grant.SQLRevokeStatement(); actual != expectedRevoke {
		t.Fatalf("expected revoke SQL %q, got %q", expectedRevoke, actual)
	}
}
