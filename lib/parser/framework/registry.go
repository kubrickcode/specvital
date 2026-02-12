package framework

import (
	"sort"
	"sync"

	"github.com/kubrickcode/specvital/lib/parser/domain"
)

var defaultRegistry = NewRegistry()

// Registry manages registered framework definitions (thread-safe).
type Registry struct {
	mu         sync.RWMutex
	frameworks map[string]*Definition
	byLanguage map[domain.Language][]*Definition
	byPriority []*Definition
}

func NewRegistry() *Registry {
	return &Registry{
		frameworks: make(map[string]*Definition),
		byLanguage: make(map[domain.Language][]*Definition),
		byPriority: []*Definition{},
	}
}

func DefaultRegistry() *Registry {
	return defaultRegistry
}

// Register adds a framework definition to the default registry.
// Typically called from framework package init() functions.
func Register(def *Definition) {
	defaultRegistry.Register(def)
}

func (r *Registry) Register(def *Definition) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.frameworks[def.Name] = def
	for _, lang := range def.Languages {
		r.byLanguage[lang] = append(r.byLanguage[lang], def)
	}
	r.byPriority = append(r.byPriority, def)
	r.sortByPriority()
}

func (r *Registry) sortByPriority() {
	sort.Slice(r.byPriority, func(i, j int) bool {
		if r.byPriority[i].Priority != r.byPriority[j].Priority {
			return r.byPriority[i].Priority > r.byPriority[j].Priority
		}
		return r.byPriority[i].Name < r.byPriority[j].Name
	})

	for lang := range r.byLanguage {
		defs := r.byLanguage[lang]
		sort.Slice(defs, func(i, j int) bool {
			if defs[i].Priority != defs[j].Priority {
				return defs[i].Priority > defs[j].Priority
			}
			return defs[i].Name < defs[j].Name
		})
	}
}

func Find(name string) *Definition {
	return defaultRegistry.Find(name)
}

func (r *Registry) Find(name string) *Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.frameworks[name]
}

// FindByLanguage returns frameworks supporting the language (sorted by priority).
func FindByLanguage(lang domain.Language) []*Definition {
	return defaultRegistry.FindByLanguage(lang)
}

func (r *Registry) FindByLanguage(lang domain.Language) []*Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defs := r.byLanguage[lang]
	result := make([]*Definition, len(defs))
	copy(result, defs)
	return result
}

// All returns all registered frameworks (sorted by priority).
func All() []*Definition {
	return defaultRegistry.All()
}

func (r *Registry) All() []*Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Definition, len(r.byPriority))
	copy(result, r.byPriority)
	return result
}

// Clear removes all frameworks (useful for testing).
func Clear() {
	defaultRegistry.Clear()
}

func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.frameworks = make(map[string]*Definition)
	r.byLanguage = make(map[domain.Language][]*Definition)
	r.byPriority = []*Definition{}
}

func Names() []string {
	return defaultRegistry.Names()
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.frameworks))
	for name := range r.frameworks {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
