package model

// User 用户模型
type User struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	UserName   string `gorm:"not null" json:"username,omitempty"`
	Password   string `gorm:"not null" json:"password,omitempty"`
	CreateTime int64  `gorm:"not null" json:"createTime,omitempty"`
	UpdateTime int64  `gorm:"not null" json:"updateTime,omitempty"`
}

// Role 角色模型
type Role struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name      string `gorm:"not null" json:"name,omitempty"`
	Namespace string `gorm:"not null" json:"namespace,omitempty"`
}

// RoleBinding 角色绑定模型
type RoleBinding struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name      string `gorm:"not null" json:"name,omitempty"`
	Namespace string `gorm:"not null" json:"namespace,omitempty"`
	// 可有多个 user
	Users  string `gorm:"not null" json:"users,omitempty"`
	RoleID int    `gorm:"not null" json:"roleID,omitempty"`
}

// RoleSubRef 角色子参考模型
type RoleSubRef struct {
	ID     int `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	RoleID int `gorm:"not null" json:"roleID,omitempty"`
	// 可有多个 verb
	Verbs string `gorm:"not null" json:"verbs,omitempty"`
	// 可有多个 resource
	Resources string `gorm:"not null" json:"resources,omitempty"`
}

// ClusterRoleBinding 集群角色绑定模型
type ClusterRoleBinding struct {
	ID     int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name   string `gorm:"not null" json:"name,omitempty"`
	Users  string `gorm:"not null" json:"users,omitempty"`
	RoleID int    `gorm:"not null" json:"roleID,omitempty"`
}

// ClusterRole 集群角色模型
type ClusterRole struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	Name string `gorm:"not null" json:"name,omitempty"`
}

// ClusterRoleSubRef 集群角色子参考模型
type ClusterRoleSubRef struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`
	ClusterRoleID int    `gorm:"not null" json:"clusterroleID,omitempty"`
	Verbs         string `gorm:"not null" json:"verbs,omitempty"`
	Resources     string `gorm:"not null" json:"resources,omitempty"`
}

