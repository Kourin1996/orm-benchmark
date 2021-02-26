package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Models holds the schema definition for the Models entity.
type Models struct {
	ent.Schema
}

// Fields of the Models.
func (Models) Fields() []ent.Field {
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

// Edges of the Models.
func (Models) Edges() []ent.Edge {
	return nil
}
