package main

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/config"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
	"github.com/Ai-feier/rbacapp/router"
	"github.com/Ai-feier/rbacapp/service"
	"net/http"
)

func main() {
	config.InitConfig()
	dao.InitDB()
	dao.Migrate()

	router := router.NewRouter()
	http.ListenAndServe("192.168.1.100:8080", router)

	//ctx := context.Background()
	//userdao := dao.NewUserDao(ctx)
	//user := &model.User{
	//	ID: 1,
	//	UserName:   "test",
	//	Password:   "123456",
	//}
	//userdao.UpdateUser(user)

	//TestClusterRoleCreate()
	//TestClusterRoleDelete()
	//TestClusterRoleBindingCreate()
	//TestClusterRoleBindingDelete()
	//TestRoleCreate()
	//TestRoleDelete()
	//TestRoleBindingCreate()
	//TestRoleBindingDelete()
}

func TestClusterRoleCreate() {
	// clusterrole
	testCases := []struct{
		name string
		clustername string
		subs []*model.ClusterRoleSubRef
		err error
	} {
		{
			name: "cr1",
			clustername: "cr1",
			subs: []*model.ClusterRoleSubRef{
				{
					Verbs: "get,create",
					Resources: "pods,deploymentes",
				},
				{
					Verbs: "delete,list",
					Resources: "daemonsets,statefulsets",
				},
			},
		},
	}

	ctx := context.Background()
	crSvc := service.NewClusterRoleSvc()
	for _, tc := range testCases {
		cr := &model.ClusterRole{Name: tc.name}
		err := crSvc.CreateClusterRole(ctx, cr, tc.subs...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestClusterRoleDelete() {
	// clusterrole
	testCases := []struct{
		name string
		clustername string
		subs []*model.ClusterRoleSubRef
		err error
	} {
		{
			name: "cr1",
			clustername: "cr1",
			subs: []*model.ClusterRoleSubRef{
				{
					Verbs: "get,create",
					Resources: "pods,deploymentes",
				},
				{
					Verbs: "delete,list",
					Resources: "daemonsets,statefulsets",
				},
			},
		},
	}

	ctx := context.Background()
	crSvc := service.NewClusterRoleSvc()
	for _, tc := range testCases {
		err := crSvc.DeleteClusterRole(ctx, tc.clustername)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestClusterRoleBindingCreate() {
	crbSvc :=service.NewClusterRoleBindingSvc()
	testCases := []struct{
		name string
		clustername string
		crbName string
		crbUser string
		subs []*model.ClusterRoleSubRef
		err error
	} {
		{
			name: "cr1",
			clustername: "cr1",
			crbName: "crb1",
			crbUser: "test1,test2",
			subs: []*model.ClusterRoleSubRef{
				{
					Verbs: "get,create",
					Resources: "pods,deploymentes",
				},
				{
					Verbs: "delete,list",
					Resources: "daemonsets,statefulsets",
				},
			},
		},
	}

	ctx := context.Background()
	crSvc := service.NewClusterRoleSvc()
	for _, tc := range testCases {
		cr := &model.ClusterRole{Name: tc.name}
		err := crSvc.CreateClusterRole(ctx, cr, tc.subs...)
		if err != nil {
			fmt.Println(err)
		}
		err = crbSvc.CreateClusterRoleBinding(ctx, tc.crbName, tc.crbUser, cr)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestClusterRoleBindingDelete() {
	crbSvc :=service.NewClusterRoleBindingSvc()
	testCases := []struct{
		name string
		crbName string
		err error
	} {
		{
			name: "cr1",
			crbName: "crb1",
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		err := crbSvc.DeleteClusterRoleBindingByName(ctx, tc.crbName)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestRoleCreate() {
	ctx := context.Background()
	roleSvc := service.NewRoleSvc()

	role := &model.Role{
		Name:      "role2",
		Namespace: "ns1",
	}
	rolesubs := []*model.RoleSubRef{
		{
			Verbs:     "get,list",
			Resources: "pods",
		},
		{
			Verbs: "watch",
			Resources: "deployments,daemonsets",
		},
	}
	err := roleSvc.CreateRole(ctx, role, rolesubs...)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRoleDelete() {
	ctx := context.Background()
	roleSvc := service.NewRoleSvc()

	role := &model.Role{
		Name:      "role2",
		Namespace: "ns1",
	}
	err := roleSvc.DeleteRole(ctx, role)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRoleBindingCreate() {
	ctx := context.Background()
	roleSvc := service.NewRoleSvc()

	role := &model.Role{
		Name:      "role2",
		Namespace: "ns1",
	}
	rolesubs := []*model.RoleSubRef{
		{
			Verbs:     "get,list",
			Resources: "pods",
		},
		{
			Verbs: "watch",
			Resources: "deployments,daemonsets",
		},
	}
	err := roleSvc.CreateRole(ctx, role, rolesubs...)
	if err != nil {
		fmt.Println(err)
	}

	rbSvc := service.NewRoleBindingSvc()
	rb := &model.RoleBinding{
		Name:      "rb2",
		Namespace: "ns1",
		Users:     "user1,user2",
	}
	err = rbSvc.CreateRoleBinding(ctx, rb, role)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRoleBindingDelete() {
	rb := &model.RoleBinding{
		Name:      "rb2",
		Namespace: "ns1",
		Users:     "user1,user2",
	}
	rbSvc := service.NewRoleBindingSvc()
	err := rbSvc.DeleteRoleBinding(context.Background(), rb)
	if err != nil {
		fmt.Println(err)
	}
}
