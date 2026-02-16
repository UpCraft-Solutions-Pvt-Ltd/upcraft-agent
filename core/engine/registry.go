package engine

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/UpCraft-Solutions-Pvt-Ltd/upcraft-agent/core/skills"
)

// ActionHandler executes one registered skill action.
type ActionHandler func(ctx context.Context, input map[string]interface{}) *ActionResult

// RegisteredAction binds a skill action contract to executable code.
type RegisteredAction struct {
	Skill       string
	Action      string
	Description string
	InputSchema map[string]interface{}
	Handler     ActionHandler
}

type Registry struct {
	mu      sync.RWMutex
	actions map[string]RegisteredAction
}

func NewRegistry() *Registry {
	return &Registry{actions: make(map[string]RegisteredAction)}
}

func (r *Registry) Register(a RegisteredAction) error {
	if strings.TrimSpace(a.Skill) == "" || strings.TrimSpace(a.Action) == "" {
		return fmt.Errorf("skill and action are required")
	}
	if a.Handler == nil {
		return fmt.Errorf("handler is required for %s.%s", a.Skill, a.Action)
	}

	key := actionKey(a.Skill, a.Action)

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.actions[key]; exists {
		return fmt.Errorf("action already registered: %s", key)
	}
	if a.InputSchema == nil {
		a.InputSchema = map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}
	}
	r.actions[key] = a
	return nil
}

func (r *Registry) Execute(ctx context.Context, skillName, actionName string, input map[string]interface{}) *ActionResult {
	if input == nil {
		input = map[string]interface{}{}
	}

	r.mu.RLock()
	a, ok := r.actions[actionKey(skillName, actionName)]
	r.mu.RUnlock()
	if !ok {
		return ErrorResult(fmt.Sprintf("unknown action: %s.%s", skillName, actionName), fmt.Errorf("action not registered"))
	}

	result := a.Handler(ctx, input)
	if result == nil {
		return ErrorResult(fmt.Sprintf("action returned nil result: %s.%s", skillName, actionName), fmt.Errorf("nil action result"))
	}
	return result
}

func (r *Registry) ToProviderDefs() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defs := make([]ToolDefinition, 0, len(r.actions))
	for _, a := range r.actions {
		defs = append(defs, ToolDefinition{
			Type: "function",
			Function: ToolFunctionDefinition{
				Name:        a.Skill + "." + a.Action,
				Description: a.Description,
				Parameters:  a.InputSchema,
			},
		})
	}
	return defs
}

func (r *Registry) SkillDefinitions() []skills.SkillDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bySkill := map[string][]skills.ActionDefinition{}
	for _, a := range r.actions {
		inputSchema := "{}"
		if a.InputSchema != nil {
			inputSchema = fmt.Sprintf("%v", a.InputSchema)
		}
		bySkill[a.Skill] = append(bySkill[a.Skill], skills.ActionDefinition{
			Name:        a.Action,
			Description: a.Description,
			InputSchema: inputSchema,
		})
	}

	names := make([]string, 0, len(bySkill))
	for name := range bySkill {
		names = append(names, name)
	}
	sort.Strings(names)

	definitions := make([]skills.SkillDefinition, 0, len(names))
	for _, name := range names {
		actions := bySkill[name]
		sort.Slice(actions, func(i, j int) bool {
			return actions[i].Name < actions[j].Name
		})
		definitions = append(definitions, skills.SkillDefinition{
			Name:        name,
			Description: "Registered skill actions",
			Actions:     actions,
		})
	}
	return definitions
}

func actionKey(skillName, actionName string) string {
	return strings.ToLower(strings.TrimSpace(skillName)) + "." + strings.ToLower(strings.TrimSpace(actionName))
}
