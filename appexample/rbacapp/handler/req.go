package handler

import "github.com/Ai-feier/rbacapp/model"

type UserRequest struct {
	User *model.User `json:"user"`
}

type RoleRequest struct {
	Role *model.Role `json:"role"`
	RoleSubRefs []*model.RoleSubRef `json:"roleSubRef"`
}

type RoleBindingRequest struct {
	RoleBinding *model.RoleBinding `json:"roleBinding"`
	Role        *model.Role        `json:"role"`
}
