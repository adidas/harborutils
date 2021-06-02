package main

import (
	"fmt"
	"time"
)

type ClientPrt struct {
	Url         string
	Method      string
	ContentType string
	Password    string
	User        string
	Body        interface{}
}

//HTTP body request struct
type AddMember struct {
	RoleID      int            `json:"role_id,omitempty"`
	MemberGroup AddMemberGroup `json:"member_group,omitempty"`
	MemberUser  AddMemberUser  `json:"member_user,omitempty"`
}

type AddMemberGroup struct {
	GroupName   string `json:"group_name,omitempty"`
	LdapGroupDN string `json:"ldap_group_dn,omitempty"`
	GroupType   int    `json:"group_type,omitempty"`
	ID          int    `json:"id,omitempty"`
}
type AddMemberUser struct {
	Username string `json:"username,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
}

type Matcher struct {
	Id     string `json:"id,omitempty"`
	Active bool   `json:"active,omitempty"`
	Type   `json:"type,omitempty"`
}

type Type struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

//HTTP body response struct
type ProjectListResponse []ProjectResponse

type ProjectResponse struct {
	ProjectId         int          `json:"project_id,omitempty"`
	OwnerId           int          `json:"owner_id,omitempty"`
	Name              string       `json:"name,omitempty"`
	CreationTime      string       `json:"creation_time,omitempty"`
	UpdateTime        string       `json:"update_time,omitempty"`
	Deleted           bool         `json:"deleted,omitempty"`
	OwnerName         bool         `json:"owner_name,omitempty"`
	Toggeable         bool         `json:"toggeable,omitempty"`
	CurrentUserRoleId int          `json:"current_user_role_id,omitempty"`
	RepoCount         int          `json:"repo_count,omitempty"`
	ChartCount        int          `json:"chart_count,omitempty"`
	Metadata          Metadata     `json:"metadata,omitempty"`
	CveWhitelist      CveWhitelist `json:"cve_whitelist,omitempty"`
}

type LabelListResponse []LabelResponse

type LabelResponse struct {
	UpdateTime   string `json:"update_time,omitempty"`
	Description  string `json:"description,omitempty"`
	Color        string `json:"color,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	Deleted      bool   `json:"deleted,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ProjectID    int    `json:"project_id,omitempty"`
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
}

type ArtifactResponse struct {
	Digest string `json:"digest,omitempty"`
}

type Metadata struct {
	AutoScan           bool   `json:"auto_scan,omitempty"`
	EnableContentTrust bool   `json:"enable_content_trust,omitempty"`
	PreventVul         bool   `json:"prevent_vul,omitempty"`
	Public             bool   `json:"public,omitempty"`
	Severity           string `json:"severity,omitempty"`
}

type CveWhitelist struct {
	Id        int `json:"id,omitempty"`
	ProjectId int `json:"project_id,omitempty"`
	//Items  "items": null,
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
}

type MemberListResponse []MemberResponse

type MemberResponse struct {
	EntityId   int    `json:"entity_id,omitempty"`
	RoleName   string `json:"role_name,omitempty"`
	EntityName string `json:"entity_name,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	ProjectId  int    `json:"project_id,omitempty"`
	Id         int    `json:"id,omitempty"`
	RoleId     int    `json:"role_id,omitempty"`
}

type UserListResponse []UserResponse

type UserResponse struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	//LdapRealname string   `json:"ldap_realname,omitempty"`
	//LdapGroupdn  []string `json:"ldap_groupdn,omitempty"`
}

type UserImport struct {
	LdapUIDList []string `json:"ldap_uid_list"`
}

type GroupListResponse []GroupResponse

type GroupResponse struct {
	GroupName   string `json:"group_name,omitempty"`
	LdapGroupDN string `json:"ldap_group_dn,omitempty"`
	GroupType   int    `json:"group_type,omitempty"`
	GroupID     int    `json:"id,omitempty"`
}

type RobotResponse struct {
	UpdateTime   time.Time `json:"update_time"`
	Description  string    `json:"description"`
	Level        string    `json:"level"`
	Editable     bool      `json:"editable"`
	CreationTime time.Time `json:"creation_time"`
	ExpiresAt    int       `json:"expires_at"`
	Name         string    `json:"name"`
	Secret       string    `json:"secret"`
	Disable      bool      `json:"disable"`
	Duration     int       `json:"duration"`
	ID           int       `json:"id,omitempty"`
	Permissions  []struct {
		Access []struct {
			Action   string `json:"action"`
			Resource string `json:"resource"`
			Effect   string `json:"effect"`
		} `json:"access"`
		Kind      string `json:"kind"`
		Namespace string `json:"namespace"`
	} `json:"permissions"`
}

func (robot RobotResponse) String() string {
	return fmt.Sprintf("<robot %d> name: %s level: %s", robot.ID, robot.Name, robot.Level)
}

type RobotCreatedReponse struct {
	Secret       string    `json:"secret"`
	CreationTime time.Time `json:"creation_time"`
	ID           int       `json:"id"`
	ExpiresAt    int       `json:"expires_at"`
	Name         string    `json:"name"`
}

type AuditLog struct {
	Username     string `json:"username"`
	Resource     string `json:"resource"`
	Operation    string `json:"operation"`
	OpTime       string `json:"op_time"`
	ID           int    `json:"id"`
	ResourceType string `json:"resource_type"`
}

type ReplicationExecution struct {
	Status     string    `json:"status"`
	StatusText string    `json:"status_text"`
	Trigger    string    `json:"trigger"`
	StartTime  time.Time `json:"start_time"`
	Failed     int       `json:"failed"`
	Succeed    int       `json:"succeed"`
	Stopped    int       `json:"stopped"`
	EndTime    time.Time `json:"end_time"`
	InProgress int       `json:"in_progress"`
	Total      int       `json:"total"`
	ID         int       `json:"id"`
	PolicyID   int       `json:"policy_id"`
}

type StartReplicationExecution struct {
	PolicyID int `json:"policy_id"`
}

type RPFilter struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ReplicationPolicy struct {
	UpdateTime   string     `json:"update_time"`
	Description  string     `json:"description"`
	Enabled      bool       `json:"enabled"`
	Filters      []RPFilter `json:"filters"`
	DestRegistry struct {
		Status     string `json:"status"`
		Credential struct {
			AccessKey    string `json:"access_key"`
			AccessSecret string `json:"access_secret"`
			Type         string `json:"type"`
		} `json:"credential"`
		UpdateTime   string `json:"update_time"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		Insecure     bool   `json:"insecure"`
		CreationTime string `json:"creation_time"`
		Type         string `json:"type"`
		ID           int    `json:"id"`
		Description  string `json:"description"`
	} `json:"dest_registry"`
	CreationTime string `json:"creation_time"`
	SrcRegistry  struct {
		Status     string `json:"status"`
		Credential struct {
			AccessKey    string `json:"access_key"`
			AccessSecret string `json:"access_secret"`
			Type         string `json:"type"`
		} `json:"credential"`
		UpdateTime   string `json:"update_time"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		Insecure     bool   `json:"insecure"`
		CreationTime string `json:"creation_time"`
		Type         string `json:"type"`
		ID           int    `json:"id"`
		Description  string `json:"description"`
	} `json:"src_registry"`
	DestNamespace string `json:"dest_namespace"`
	Trigger       struct {
		Type            string `json:"type"`
		TriggerSettings struct {
			Cron string `json:"cron"`
		} `json:"trigger_settings"`
	} `json:"trigger"`
	Deletion bool   `json:"deletion"`
	Override bool   `json:"override"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

type ReplicationTask struct {
	Status       string    `json:"status"`
	JobID        string    `json:"job_id"`
	StartTime    time.Time `json:"start_time"`
	DstResource  string    `json:"dst_resource"`
	SrcResource  string    `json:"src_resource"`
	ResourceType string    `json:"resource_type"`
	Operation    string    `json:"operation"`
	ID           int       `json:"id"`
	ExecutionID  int       `json:"execution_id"`
	EndTime      time.Time `json:"end_time"`
}

// prefix Db, to reference database models

type DbHarborUser struct {
	// gorm.Model
	UserId          uint        `gorm:"primary_key;column:user_id"`
	Username        string      `gorm:"column:username"`
	Email           string      `gorm:"column:email"`
	Password        string      `gorm:"column:password"`
	Realname        string      `gorm:"column:realname"`
	Comment         string      `gorm:"column:comment"`
	Deleted         bool        `gorm:"column:deleted"`
	ResetUuid       string      `gorm:"column:reset_uuid"`
	Salt            string      `gorm:"column:salt"`
	SysadminFlag    bool        `gorm:"column:sysadmin_flag"`
	CreationTime    time.Time   `gorm:"column:creation_time"`
	UpdateTime      time.Time   `gorm:"column:update_time"`
	PasswordVersion string      `gorm:"column:password_version"`
	OidcUser        *DBOidcUser `gorm:"foreignKey:UserId"`
}

func (DbHarborUser) TableName() string {
	return "harbor_user"
}

func (u DbHarborUser) String() string {
	return fmt.Sprintf("<harbor_user %d> - username: %s email:%s", u.UserId, u.Username, u.Email)
}

type DBOidcUser struct {
	Id           uint      `gorm:"primary_key;column:id"`
	UserId       uint      `gorm:"column:user_id"`
	Secret       string    `gorm:"column:secret"` // secret token to use in harbor
	Subiss       string    `gorm:"column:subiss"` //
	Token        string    `gorm:"column:token"`
	CreationTime time.Time `gorm:"column:creation_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (DBOidcUser) TableName() string {
	return "oidc_user"
}

func (oidc DBOidcUser) String() string {
	return fmt.Sprintf("<oidc_user %d> - user_id: %d", oidc.UserId, oidc.UserId)
}

type DbRobot struct {
	Id     uint   `gorm:"primary_key;column:id"`
	Name   string `gorm:"column:name"`
	Secret string `gorm:"column:secret"`
	Salt   string `gorm:"column:salt"`
}

func (DbRobot) TableName() string {
	return "robot"
}

func (robot DbRobot) String() string {
	return fmt.Sprintf("<robot_db %d> - name: %s", robot.Id, robot.Name)
}
