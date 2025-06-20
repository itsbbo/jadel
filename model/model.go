package model

import (
	"time"

	"github.com/guregu/null/v6"
	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
)

type Timestamp struct {
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp()" json:"updated_at"`
}

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            ulid.ULID `bun:"id,pk"       json:"id"`
	Name          string    `bun:"name"        json:"name"`
	Email         string    `bun:"email"       json:"email"`
	Password      string    `bun:"password"    json:"-"`
	Timestamp
}

type Session struct {
	bun.BaseModel `bun:"table:sessions"`
	ID            string      `bun:"id,pk"          json:"id"`
	UserID        ulid.ULID   `bun:"user_id"        json:"user_id"`
	IPAddress     null.String `bun:"ip_address"     json:"ip_address"`
	UserAgent     null.String `bun:"user_agent"     json:"user_agent"`
	ExpiredAt     time.Time   `bun:"expired_at"     json:"expired_at"`

	User User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type Project struct {
	bun.BaseModel `bun:"table:projects"`
	ID            ulid.ULID   `bun:"id,pk"          json:"id"`
	Name          string      `bun:"name"           json:"name"`
	UserID        ulid.ULID   `bun:"user_id"        json:"user_id"`
	Description   null.String `bun:"description"    json:"description"`
	Timestamp

	Environments []Environment `bun:"rel:has-many,join:id=project_id" json:"environments"`
}

type PrivateKey struct {
	bun.BaseModel `bun:"table:private_keys"`
	ID            ulid.ULID   `bun:"id,pk"              json:"id"`
	Name          string      `bun:"name"               json:"name"`
	UserID        ulid.ULID   `bun:"user_id"            json:"user_id"`
	Description   null.String `bun:"description"        json:"description"`
	PrivateKey    string      `bun:"private_key"        json:"private_key"`
	IsGitRelated  bool        `bun:"is_git_related"     json:"is_git_related"`
	Timestamp
}

type Server struct {
	bun.BaseModel `bun:"table:servers"`
	ID            ulid.ULID   `bun:"id,pk"          json:"id"`
	Name          string      `bun:"name"           json:"name"`
	UserID        ulid.ULID   `bun:"user_id"        json:"user_id"`
	Description   null.String `bun:"description"    json:"description"`
	IP            string      `bun:"ip"             json:"ip"`
	Port          int         `bun:"port"           json:"port"`
	Username      string      `bun:"username"       json:"username"`
	PrivateKeyID  ulid.ULID   `bun:"private_key_id" json:"private_key_id"`
	Timestamp
}

type Environment struct {
	bun.BaseModel `bun:"table:environments"`
	ID            ulid.ULID `bun:"id,pk"              json:"id"`
	Name          string    `bun:"name"               json:"name"`
	ProjectID     ulid.ULID `bun:"project_id"         json:"project_id"`
	Timestamp

	Project      Project       `bun:"rel:belongs-to,join:project_id=id"   json:"project"`
	Applications []Application `bun:"rel:has-many,join:id=environment_id" json:"applications"`
	Databases    []Database    `bun:"rel:has-many,join:id=environment_id" json:"databases"`
}

type Application struct {
	bun.BaseModel        `bun:"table:applications"`
	ID                   ulid.ULID         `bun:"id,pk"                  json:"id"`
	EnvironmentID        ulid.ULID         `bun:"environment_id"         json:"environment_id"`
	Name                 string            `bun:"name"                   json:"name"`
	Description          null.String       `bun:"description"            json:"description"`
	BuildTool            string            `bun:"build_tool"             json:"build_tool"`
	Domain               string            `bun:"domain"                 json:"domain"`
	EnableHTTPS          bool              `bun:"enable_https"           json:"enable_https"`
	PreDeploymentScript  null.String       `bun:"pre_deployment_script"  json:"pre_deployment_script"`
	PostDeploymentScript null.String       `bun:"post_deployment_script" json:"post_deployment_script"`
	PortMappings         map[string]string `bun:"port_mappings"          json:"port_mappings"`
	Metadata             map[string]any    `bun:"metadata"               json:"metadata"`
	Timestamp

	Environment Environment `bun:"rel:belongs-to,join:environment_id=id" json:"environment"`
}

type Database struct {
	bun.BaseModel `bun:"table:databases"`
	ID            ulid.ULID         `bun:"id,pk"           json:"id"`
	EnvironmentID ulid.ULID         `bun:"environment_id"  json:"environment_id"`
	DatabaseType  string            `bun:"database_type"   json:"database_type"`
	Name          string            `bun:"name"            json:"name"`
	Description   null.String       `bun:"description"     json:"description"`
	Image         string            `bun:"image"           json:"image"`
	Username      string            `bun:"username"        json:"username"`
	Password      null.String       `bun:"password"        json:"password"`
	PortMappings  map[string]string `bun:"port_mappings"   json:"port_mappings"`
	CustomConfig  null.String       `bun:"custom_config"   json:"custom_config"`
	Metadata      map[string]any    `bun:"metadata"        json:"metadata"`
	Timestamp

	Environment Environment `bun:"rel:belongs-to,join:environment_id=id" json:"environment"`
}
