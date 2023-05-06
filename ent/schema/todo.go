package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("oid"),
		field.String("task"),
		field.Bool("completed").Default(false).StructTag(`json:"completed"`),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}
