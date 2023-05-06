// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"
	"todo-app/ent/schema"
	"todo-app/ent/todo"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	todoFields := schema.Todo{}.Fields()
	_ = todoFields
	// todoDescCompleted is the schema descriptor for completed field.
	todoDescCompleted := todoFields[2].Descriptor()
	// todo.DefaultCompleted holds the default value on creation for the completed field.
	todo.DefaultCompleted = todoDescCompleted.Default.(bool)
	// todoDescCreatedAt is the schema descriptor for created_at field.
	todoDescCreatedAt := todoFields[3].Descriptor()
	// todo.DefaultCreatedAt holds the default value on creation for the created_at field.
	todo.DefaultCreatedAt = todoDescCreatedAt.Default.(func() time.Time)
	// todoDescID is the schema descriptor for id field.
	todoDescID := todoFields[0].Descriptor()
	// todo.DefaultID holds the default value on creation for the id field.
	todo.DefaultID = todoDescID.Default.(func() uuid.UUID)
}