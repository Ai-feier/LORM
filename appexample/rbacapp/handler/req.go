package handler

import "github.com/Ai-feier/rbacapp/model"

type UserRequest struct {
	User *model.Users `json:"user"`
}

type RoleRequest struct {
	Role *model.Roles                `json:"role"`
	RoleSubRefs []*model.RoleSubRefs `json:"roleSubRef"`
}

type RoleBindingRequest struct {
	RoleBinding *model.RoleBindings `json:"roleBinding"`
	Role        *model.Roles        `json:"role"`
}
