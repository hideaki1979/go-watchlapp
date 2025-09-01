package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Movie holds the schema definition for the Movie entity.
type Movie struct {
	ent.Schema
}

// Fields of the Movie.
func (Movie) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty().Comment("映画・ドラマのタイトル"),
		field.Text("description").Optional().Comment("概要・あらすじ"),
		field.String("genre").
			Optional().
			Comment("ジャンル"),
		field.Int("release_year").
			Optional().
			Comment("公開年"),
		field.String("poster_url").
			Optional().
			Comment("ポスター画像URL"),
		field.Enum("media_type").
			Values("movie", "tv_series", "documentary", "anime").
			Default("movie").
			Comment("メディアタイプ"),
		field.Enum("watch_status").
			Values("want_to_watch", "watching", "completed", "dropped").
			Default("want_to_watch").
			Comment("視聴ステータス"),
		field.Int("rating").
			Optional().
			Range(1, 5).
			Comment("評価（1-5）"),
		field.Text("review").
			Optional().
			Comment("レビュー・感想"),
		field.Time("watched_at").
			Optional().
			Comment("視聴完了日"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("作成日時"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新日時"),
	}
}

// Edges of the Movie.
func (Movie) Edges() []ent.Edge {
	return nil
}

// Indexes of the Movie.
func (Movie) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
		index.Fields("genre"),
		index.Fields("watch_status"),
		index.Fields("media_type"),
		index.Fields("created_at"),
	}
}
