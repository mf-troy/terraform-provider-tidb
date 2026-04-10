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
			input:    "teleport_reader",
			expected: "teleport_reader",
		},
		{
			name:     "quoted percent-host role normalizes to bare name",
			input:    "'teleport_reader'@'%'",
			expected: "teleport_reader",
		},
		{
			name:     "explicit percent-host role normalizes to bare name",
			input:    "teleport_reader@%",
			expected: "teleport_reader",
		},
		{
			name:     "non-default host is preserved",
			input:    "'teleport_reader'@'internal'",
			expected: "teleport_reader@internal",
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
			input:    "teleport_reader",
			expected: "'teleport_reader'@'%'",
		},
		{
			name:     "legacy quoted role stays canonical",
			input:    "'teleport_reader'@'%'",
			expected: "'teleport_reader'@'%'",
		},
		{
			name:     "non-default host is preserved",
			input:    "teleport_reader@internal",
			expected: "'teleport_reader'@'internal'",
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
			"teleport_reader",
			"'teleport_creator'@'%'",
		},
		UserOrRole: UserOrRole{
			Name: "networker",
			Host: "%",
		},
	}

	expectedGrant := "GRANT 'teleport_reader'@'%', 'teleport_creator'@'%' TO 'networker'@'%'"
	expectedRevoke := "REVOKE 'teleport_reader'@'%', 'teleport_creator'@'%' FROM 'networker'@'%'"

	if actual := grant.SQLGrantStatement(); actual != expectedGrant {
		t.Fatalf("expected grant SQL %q, got %q", expectedGrant, actual)
	}

	if actual := grant.SQLRevokeStatement(); actual != expectedRevoke {
		t.Fatalf("expected revoke SQL %q, got %q", expectedRevoke, actual)
	}
}
