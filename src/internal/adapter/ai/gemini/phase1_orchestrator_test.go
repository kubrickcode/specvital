package gemini

import (
	"testing"

	"github.com/specvital/worker/internal/domain/specview"
)

func TestMergeAssignmentsToPhase1Output_BasicMerge(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Name:        "Authentication",
				Description: "Auth features",
				Features: []specview.TaxonomyFeature{
					{Name: "Login", FileIndices: []int{0}},
					{Name: "Logout", FileIndices: []int{1}},
				},
			},
			{
				Name:        "Payment",
				Description: "Payment processing",
				Features: []specview.TaxonomyFeature{
					{Name: "Stripe", FileIndices: []int{2}},
				},
			},
		},
	}

	assignments := []specview.AssignmentOutput{
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Authentication", Feature: "Login", TestIndices: []int{0, 1}},
				{Domain: "Authentication", Feature: "Logout", TestIndices: []int{2}},
			},
		},
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Payment", Feature: "Stripe", TestIndices: []int{3, 4}},
			},
		},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "auth/login_test.go", Tests: []specview.TestInfo{{Index: 0}, {Index: 1}}},
			{Path: "auth/logout_test.go", Tests: []specview.TestInfo{{Index: 2}}},
			{Path: "payment/stripe_test.go", Tests: []specview.TestInfo{{Index: 3}, {Index: 4}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, assignments, input)

	if len(output.Domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(output.Domains))
	}

	authDomain := output.Domains[0]
	if authDomain.Name != "Authentication" {
		t.Errorf("expected first domain 'Authentication', got %q", authDomain.Name)
	}
	if len(authDomain.Features) != 2 {
		t.Errorf("expected 2 features in Authentication, got %d", len(authDomain.Features))
	}

	loginFeature := authDomain.Features[0]
	if loginFeature.Name != "Login" {
		t.Errorf("expected first feature 'Login', got %q", loginFeature.Name)
	}
	if len(loginFeature.TestIndices) != 2 {
		t.Errorf("expected 2 test indices in Login, got %d", len(loginFeature.TestIndices))
	}
}

func TestMergeAssignmentsToPhase1Output_EmptyAssignments(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Name: "Auth",
				Features: []specview.TaxonomyFeature{
					{Name: "Login", FileIndices: []int{0}},
				},
			},
		},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, nil, input)

	// No assignments means fallback to Uncategorized
	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}
	if output.Domains[0].Name != uncategorizedDomainName {
		t.Errorf("expected %q domain, got %q", uncategorizedDomainName, output.Domains[0].Name)
	}
}

func TestMergeAssignmentsToPhase1Output_SkipsEmptyFeatures(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Name: "Auth",
				Features: []specview.TaxonomyFeature{
					{Name: "Login", FileIndices: []int{0}},
					{Name: "Logout", FileIndices: []int{1}},
				},
			},
		},
	}

	// Only assign to Login, not Logout
	assignments := []specview.AssignmentOutput{
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Auth", Feature: "Login", TestIndices: []int{0, 1}},
			},
		},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}, {Index: 1}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, assignments, input)

	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}

	// Logout feature should be skipped because no tests assigned
	if len(output.Domains[0].Features) != 1 {
		t.Errorf("expected 1 feature (Logout skipped), got %d", len(output.Domains[0].Features))
	}
}

func TestMergeAssignmentsToPhase1Output_CombinesMultipleBatchAssignments(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Name: "Auth",
				Features: []specview.TaxonomyFeature{
					{Name: "Login", FileIndices: []int{0}},
				},
			},
		},
	}

	// Two batches both assign to Login
	assignments := []specview.AssignmentOutput{
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Auth", Feature: "Login", TestIndices: []int{0, 1}},
			},
		},
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Auth", Feature: "Login", TestIndices: []int{2, 3}},
			},
		},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}, {Index: 1}, {Index: 2}, {Index: 3}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, assignments, input)

	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}

	loginFeature := output.Domains[0].Features[0]
	if len(loginFeature.TestIndices) != 4 {
		t.Errorf("expected 4 test indices from combined batches, got %d", len(loginFeature.TestIndices))
	}
}

func TestCollectAllTestIndices(t *testing.T) {
	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{
				Path: "file1.go",
				Tests: []specview.TestInfo{
					{Index: 0},
					{Index: 1},
				},
			},
			{
				Path: "file2.go",
				Tests: []specview.TestInfo{
					{Index: 2},
				},
			},
		},
	}

	indices := collectAllTestIndices(input)

	if len(indices) != 3 {
		t.Fatalf("expected 3 indices, got %d", len(indices))
	}
	if indices[0] != 0 || indices[1] != 1 || indices[2] != 2 {
		t.Errorf("expected [0, 1, 2], got %v", indices)
	}
}

func TestCollectAllTestIndices_Empty(t *testing.T) {
	input := specview.Phase1Input{
		Files: []specview.FileInfo{},
	}

	indices := collectAllTestIndices(input)

	if len(indices) != 0 {
		t.Errorf("expected 0 indices, got %d", len(indices))
	}
}

func TestAggregateUsage(t *testing.T) {
	total := &specview.TokenUsage{
		PromptTokens:     100,
		CandidatesTokens: 50,
		TotalTokens:      150,
	}

	usage := &specview.TokenUsage{
		PromptTokens:     200,
		CandidatesTokens: 100,
		TotalTokens:      300,
	}

	aggregateUsage(total, usage)

	if total.PromptTokens != 300 {
		t.Errorf("expected PromptTokens 300, got %d", total.PromptTokens)
	}
	if total.CandidatesTokens != 150 {
		t.Errorf("expected CandidatesTokens 150, got %d", total.CandidatesTokens)
	}
	if total.TotalTokens != 450 {
		t.Errorf("expected TotalTokens 450, got %d", total.TotalTokens)
	}
}

func TestAggregateUsage_NilUsage(t *testing.T) {
	total := &specview.TokenUsage{
		PromptTokens:     100,
		CandidatesTokens: 50,
		TotalTokens:      150,
	}

	aggregateUsage(total, nil)

	if total.PromptTokens != 100 {
		t.Errorf("expected PromptTokens unchanged at 100, got %d", total.PromptTokens)
	}
}

func TestGenerateHeuristicTaxonomy_GroupsByDirectory(t *testing.T) {
	files := []specview.FileInfo{
		{Path: "auth/login_test.go"},
		{Path: "auth/logout_test.go"},
		{Path: "payment/stripe_test.go"},
	}

	taxonomy := generateHeuristicTaxonomy(files)

	if len(taxonomy.Domains) != 2 {
		t.Fatalf("expected 2 domains (auth, payment), got %d", len(taxonomy.Domains))
	}

	// Verify deterministic ordering (alphabetical by directory name)
	// "auth" < "payment", so Authentication should come first
	if taxonomy.Domains[0].Name != "Authentication" {
		t.Errorf("expected first domain 'Authentication', got %q", taxonomy.Domains[0].Name)
	}
	if taxonomy.Domains[1].Name != "Payment" {
		t.Errorf("expected second domain 'Payment', got %q", taxonomy.Domains[1].Name)
	}

	// Check that domains contain correct file indices
	domainMap := make(map[string][]int)
	for _, d := range taxonomy.Domains {
		if len(d.Features) > 0 {
			domainMap[d.Name] = d.Features[0].FileIndices
		}
	}

	authIndices, ok := domainMap["Authentication"]
	if !ok {
		t.Error("expected 'Authentication' domain")
	}
	if len(authIndices) != 2 {
		t.Errorf("expected 2 files in Authentication, got %d", len(authIndices))
	}

	paymentIndices, ok := domainMap["Payment"]
	if !ok {
		t.Error("expected 'Payment' domain")
	}
	if len(paymentIndices) != 1 {
		t.Errorf("expected 1 file in Payment, got %d", len(paymentIndices))
	}
}

func TestGenerateHeuristicTaxonomy_EmptyFiles(t *testing.T) {
	taxonomy := generateHeuristicTaxonomy(nil)

	if len(taxonomy.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(taxonomy.Domains))
	}
	if taxonomy.Domains[0].Name != uncategorizedDomainName {
		t.Errorf("expected %q domain, got %q", uncategorizedDomainName, taxonomy.Domains[0].Name)
	}
}

func TestGenerateHeuristicTaxonomy_RootFiles(t *testing.T) {
	files := []specview.FileInfo{
		{Path: "main_test.go"},
		{Path: "util_test.go"},
	}

	taxonomy := generateHeuristicTaxonomy(files)

	if len(taxonomy.Domains) != 1 {
		t.Fatalf("expected 1 domain (root), got %d", len(taxonomy.Domains))
	}
	if taxonomy.Domains[0].Name != uncategorizedDomainName {
		t.Errorf("expected %q domain for root files, got %q", uncategorizedDomainName, taxonomy.Domains[0].Name)
	}
}

func TestExtractTopLevelDir(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"auth/login_test.go", "auth"},
		{"payment/stripe/api_test.go", "payment"},
		{"test.go", "root"},
		{"/root/test.go", "root"},
		{"src/internal/auth_test.go", "src"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := extractTopLevelDir(tt.path)
			if got != tt.want {
				t.Errorf("extractTopLevelDir(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestHumanizeDirName(t *testing.T) {
	tests := []struct {
		dir  string
		want string
	}{
		{"auth", "Authentication"},
		{"api", "API"},
		{"db", "Database"},
		{"payment", "Payment"},
		{"user", "User Management"},
		{"users", "User Management"},
		{"admin", "Administration"},
		{"config", "Configuration"},
		{"util", "Utilities"},
		{"utils", "Utilities"},
		{"root", uncategorizedDomainName},
		{"custom", "Custom"},    // Capitalizes unknown names
		{"Auth", "Auth"},        // Already capitalized - no change
		{"123tests", "123tests"}, // Starts with digit - no panic
		{"", uncategorizedDomainName},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			got := humanizeDirName(tt.dir)
			if got != tt.want {
				t.Errorf("humanizeDirName(%q) = %q, want %q", tt.dir, got, tt.want)
			}
		})
	}
}

func TestMergeAssignmentsToPhase1Output_PreservesConfidence(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Name: "Auth",
				Features: []specview.TaxonomyFeature{
					{Name: "Login", FileIndices: []int{0}},
				},
			},
		},
	}

	assignments := []specview.AssignmentOutput{
		{
			Assignments: []specview.TestAssignment{
				{Domain: "Auth", Feature: "Login", TestIndices: []int{0}},
			},
		},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, assignments, input)

	if output.Domains[0].Confidence != defaultClassificationConfidence {
		t.Errorf("expected domain confidence %f, got %f", defaultClassificationConfidence, output.Domains[0].Confidence)
	}
	if output.Domains[0].Features[0].Confidence != defaultClassificationConfidence {
		t.Errorf("expected feature confidence %f, got %f", defaultClassificationConfidence, output.Domains[0].Features[0].Confidence)
	}
}

func TestMergeAssignmentsToPhase1Output_NilTaxonomy(t *testing.T) {
	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}, {Index: 1}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(nil, nil, input)

	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}
	if output.Domains[0].Name != uncategorizedDomainName {
		t.Errorf("expected %q domain, got %q", uncategorizedDomainName, output.Domains[0].Name)
	}
	if len(output.Domains[0].Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(output.Domains[0].Features))
	}
	if len(output.Domains[0].Features[0].TestIndices) != 2 {
		t.Errorf("expected 2 test indices, got %d", len(output.Domains[0].Features[0].TestIndices))
	}
}

func TestMergeAssignmentsToPhase1Output_EmptyTaxonomyDomains(t *testing.T) {
	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{},
	}

	input := specview.Phase1Input{
		Files: []specview.FileInfo{
			{Path: "test.go", Tests: []specview.TestInfo{{Index: 0}}},
		},
	}

	output := mergeAssignmentsToPhase1Output(taxonomy, nil, input)

	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}
	if output.Domains[0].Name != uncategorizedDomainName {
		t.Errorf("expected %q domain, got %q", uncategorizedDomainName, output.Domains[0].Name)
	}
}
