package framework_test

import (
	"testing"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/framework"
)

func TestRegistry_Register(t *testing.T) {
	r := framework.NewRegistry()

	def := &framework.Definition{
		Name:      "test-framework",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityGeneric,
	}

	r.Register(def)

	found := r.Find("test-framework")
	if found == nil {
		t.Fatal("expected framework to be registered")
	}
	if found.Name != "test-framework" {
		t.Errorf("got name %q, want %q", found.Name, "test-framework")
	}
}

func TestRegistry_FindByLanguage(t *testing.T) {
	r := framework.NewRegistry()

	r.Register(&framework.Definition{
		Name:      "framework-ts",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityGeneric,
	})

	r.Register(&framework.Definition{
		Name:      "framework-go",
		Languages: []domain.Language{domain.LanguageGo},
		Priority:  framework.PriorityGeneric,
	})

	r.Register(&framework.Definition{
		Name:      "framework-both",
		Languages: []domain.Language{domain.LanguageTypeScript, domain.LanguageGo},
		Priority:  framework.PriorityE2E,
	})

	tsFrameworks := r.FindByLanguage(domain.LanguageTypeScript)
	if len(tsFrameworks) != 2 {
		t.Errorf("got %d TypeScript frameworks, want 2", len(tsFrameworks))
	}

	goFrameworks := r.FindByLanguage(domain.LanguageGo)
	if len(goFrameworks) != 2 {
		t.Errorf("got %d Go frameworks, want 2", len(goFrameworks))
	}
}

func TestRegistry_PriorityOrdering(t *testing.T) {
	r := framework.NewRegistry()

	r.Register(&framework.Definition{
		Name:      "generic",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityGeneric,
	})

	r.Register(&framework.Definition{
		Name:      "specialized",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PrioritySpecialized,
	})

	r.Register(&framework.Definition{
		Name:      "e2e",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityE2E,
	})

	all := r.All()
	if len(all) != 3 {
		t.Fatalf("got %d frameworks, want 3", len(all))
	}

	// Should be sorted by priority (highest first)
	if all[0].Name != "specialized" {
		t.Errorf("first framework should be 'specialized', got %q", all[0].Name)
	}
	if all[1].Name != "e2e" {
		t.Errorf("second framework should be 'e2e', got %q", all[1].Name)
	}
	if all[2].Name != "generic" {
		t.Errorf("third framework should be 'generic', got %q", all[2].Name)
	}
}

func TestRegistry_Clear(t *testing.T) {
	r := framework.NewRegistry()

	r.Register(&framework.Definition{
		Name:      "test-framework",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Priority:  framework.PriorityGeneric,
	})

	if len(r.All()) != 1 {
		t.Fatal("framework should be registered")
	}

	r.Clear()

	if len(r.All()) != 0 {
		t.Error("registry should be empty after Clear()")
	}

	if r.Find("test-framework") != nil {
		t.Error("framework should not be found after Clear()")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := framework.NewRegistry()

	r.Register(&framework.Definition{Name: "zebra", Priority: framework.PriorityGeneric})
	r.Register(&framework.Definition{Name: "alpha", Priority: framework.PriorityGeneric})
	r.Register(&framework.Definition{Name: "middle", Priority: framework.PriorityGeneric})

	names := r.Names()
	if len(names) != 3 {
		t.Fatalf("got %d names, want 3", len(names))
	}

	// Should be sorted alphabetically
	expected := []string{"alpha", "middle", "zebra"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("names[%d] = %q, want %q", i, name, expected[i])
		}
	}
}
