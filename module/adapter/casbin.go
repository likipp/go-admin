package adapter

import (
	"context"
	"fmt"
	casbinModel "github.com/casbin/casbin/v2/model"
)

type CasbinAdapter struct {
	User string
	Role string
}

func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	fmt.Println("进入加载策略")
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		fmt.Println("错误信息")
		return err
	}

	err = a.loadUserPolicy(ctx)
	if err != nil {
		fmt.Println("错误信息")
		return err
	}
	return nil
}

// 加载角色策略(p,role_id,path,method)
func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	return nil
}

// 加载用户策略(g,user_id,role_id)
func (a *CasbinAdapter) loadUserPolicy(ctx context.Context) error {
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
