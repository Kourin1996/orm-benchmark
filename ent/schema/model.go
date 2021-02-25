package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Model holds the schema definition for the Model entity.
type Model struct {
	ent.Schema
}

// Fields of the Model.
func (Model) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("fax").NotEmpty(),
		field.String("web").NotEmpty(),
		field.Int("age"),
		field.Bool("right"),
		field.Int64("counter"),
	}
}

// Edges of the Model.
func (Model) Edges() []ent.Edge {
	return nil
}
