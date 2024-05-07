package model

// Users 用户模型
type Users struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	UserName   string `gorm:"not null" json:"username,omitempty"`
	Password   string `gorm:"not null" json:"password,omitempty"`
	CreateTime int64  `gorm:"not null" json:"createTime,omitempty"`
	UpdateTime int64  `gorm:"not null" json:"updateTime,omitempty"`
}

// Roles 角色模型
type Roles struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name      string `gorm:"not null" json:"name,omitempty"`
	Namespace string `gorm:"not null" json:"namespace,omitempty"`
}

// RoleBindings 角色绑定模型
type RoleBindings struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name      string `gorm:"not null" json:"name,omitempty"`
	Namespace string `gorm:"not null" json:"namespace,omitempty"`
	// 可有多个 user
	Users  string `gorm:"not null" json:"users,omitempty"`
	RoleID int    `gorm:"not null" json:"roleID,omitempty"`
}

// RoleSubRefs 角色子参考模型
type RoleSubRefs struct {
	ID     int `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	RoleID int `gorm:"not null" json:"roleID,omitempty"`
	// 可有多个 verb
	Verbs string `gorm:"not null" json:"verbs,omitempty"`
	// 可有多个 resource
	Resources string `gorm:"not null" json:"resources,omitempty"`
}

// ClusterRoleBindings 集群角色绑定模型
type ClusterRoleBindings struct {
	ID     int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name   string `gorm:"not null" json:"name,omitempty"`
	Users  string `gorm:"not null" json:"users,omitempty"`
	RoleID int    `gorm:"not null" json:"roleID,omitempty"`
}

// ClusterRoles 集群角色模型
type ClusterRoles struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name string `gorm:"not null" json:"name,omitempty"`
}

// ClusterRoleSubRefs 集群角色子参考模型
type ClusterRoleSubRefs struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	ClusterRoleID int    `gorm:"not null" json:"clusterroleID,omitempty"`
	Verbs         string `gorm:"not null" json:"verbs,omitempty"`
	Resources     string `gorm:"not null" json:"resources,omitempty"`
}

