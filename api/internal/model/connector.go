package model
import "time"
type Connector struct {
ID        string    `db:"id" json:"id"`
Name      string    `db:"name" json:"name"`
Type      string    `db:"type" json:"type" // pg | mysql | s3 | excel | gsheets | sf | rest`
Config    string    `db:"config_json" json:"-" // encrypted JSON string`
CreatedBy string    `db:"created_by_user_id" json:"created_by_user_id"`
Workspace string    `db:"workspace_id" json:"workspace_id"`
CreatedAt time.Time `db:"created_at" json:"created_at"`
UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}