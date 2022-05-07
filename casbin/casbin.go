package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/miaocansky/go-tool/casbin/dto"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

type CasbinUtil struct {
	ModelPath string
	Db        *gorm.DB
	Enforcer  *casbin.Enforcer
}

func NewCasbinUtil(modelPath string, db *gorm.DB) (*CasbinUtil, error) {
	util := &CasbinUtil{
		ModelPath: modelPath,
		Db:        db,
	}
	_, err := util.GetEnforcer()
	return util, err
}

//单列模式生成 Enforcer
func (casbinUtil *CasbinUtil) GetEnforcer() (*casbin.Enforcer, error) {
	//modelPath:="./resource/casbin/rbac_model.conf"
	var err error
	once.Do(func() {
		a, _ := gormadapter.NewAdapterByDB(casbinUtil.Db)
		enforcer, err = casbin.NewEnforcer(casbinUtil.ModelPath, a)
	})
	casbinUtil.Enforcer = enforcer
	return enforcer, err

}

func (casbinUtil *CasbinUtil) LoadPolicy() {
	//casbinUtil.GetEnforcer()
	casbinUtil.Enforcer.LoadPolicy()

}

//
//  Enforce
//  @Description: 验证用户是否有权限
//  @receiver casbinUtil
//  @param authorityId 角色id
//  @param obj 请求地址
//  @param act 请求方法 POST GET
//  @return bool  是否有权限
//
func (casbinUtil *CasbinUtil) Enforce(authorityId, obj, act string) bool {
	//casbinUtil.GetEnforcer()
	isEnforce, err := casbinUtil.Enforcer.Enforce(authorityId, obj, act)
	if err != nil {
		return false
	}
	return isEnforce

}

func (casbinUtil *CasbinUtil) AddPolicy(authorityId uint64, casbinInfos []dto.CasbinInfoDto) bool {
	var policyLists [][]string
	if casbinInfos == nil {
		return false
	}
	for _, info := range casbinInfos {
		authorityIdStr := strconv.FormatUint(authorityId, 10)
		policyLists = append(policyLists, []string{authorityIdStr, info.Path, info.Method})
	}
	isAdd, _ := casbinUtil.Enforcer.AddPolicies(policyLists)
	return isAdd
}
func (casbinUtil *CasbinUtil) DeleteAllPolicy() bool {
	return casbinUtil.DeletePolicy(0)
}

func (casbinUtil *CasbinUtil) DeletePolicyForAuthorityId(authorityId uint64) bool {
	authorityIdStr := strconv.FormatUint(authorityId, 10)
	return casbinUtil.DeletePolicy(0, authorityIdStr)
}

func (casbinUtil *CasbinUtil) DeletePolicy(v int, p ...string) bool {
	isDel, _ := casbinUtil.Enforcer.RemoveFilteredPolicy(v, p...)
	return isDel
}
