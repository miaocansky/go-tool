package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/miaocansky/go-tool/casbin/dto"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	enforcer       *casbin.Enforcer
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

type CasbinUtil struct {
	ModelPath      string
	Db             *gorm.DB
	SyncedEnforcer *casbin.SyncedEnforcer
	SyncedTime     time.Duration
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
		syncedEnforcer, err = casbin.NewSyncedEnforcer(casbinUtil.ModelPath, a)
		if err == nil {
			casbinUtil.StartAutoLoadPolicy()
		}
		//rediswatcher.NewWatcher()
		//rediswatcher.WithRedisPubConnection()
	})
	casbinUtil.SyncedEnforcer = syncedEnforcer
	return enforcer, err

}

func (casbinUtil *CasbinUtil) LoadPolicy() {
	//casbinUtil.GetEnforcer()
	casbinUtil.SyncedEnforcer.LoadPolicy()
	//casbinUtil.Enforcer.GetAllActions()

}

//
//  startAutoLoadPolicy
//  @Description: 开启定时拉取数据库的数据(用于定时更新casbin缓存中的数据)
//  @receiver casbinUtil
//
func (casbinUtil *CasbinUtil) StartAutoLoadPolicy() {
	if casbinUtil.SyncedTime > 0 {
		casbinUtil.SyncedEnforcer.StartAutoLoadPolicy(casbinUtil.SyncedTime)
	}
}

//
//  stopAutoLoadPolicy
//  @Description:  停止定时拉取数据库的数据(关闭定时更新casbin缓存中的数据)
//  @receiver casbinUtil
//
func (casbinUtil *CasbinUtil) StopAutoLoadPolicy() {
	if casbinUtil.SyncedTime > 0 {
		casbinUtil.SyncedEnforcer.StopAutoLoadPolicy()
	}
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
	isEnforce, err := casbinUtil.SyncedEnforcer.Enforce(authorityId, obj, act)
	if err != nil {
		return false
	}
	return isEnforce

}

//
//  AddPolicy
//  @Description: 添加一个角色多个权限
//  @receiver casbinUtil
//  @param authorityId 权限id
//  @param casbinInfos 权限lists
//  @return bool
//
func (casbinUtil *CasbinUtil) AddPolicy(authorityId string, casbinInfoLists []dto.CasbinInfoDto) bool {
	var policyLists [][]string
	if casbinInfoLists == nil {
		return false
	}
	for _, info := range casbinInfoLists {
		policyLists = append(policyLists, []string{authorityId, info.Path, info.Method})
	}
	isAdd, _ := casbinUtil.SyncedEnforcer.AddPolicies(policyLists)
	return isAdd
}

//
//  DeleteAllPolicy
//  @Description: 删除所有权限
//  @receiver casbinUtil
//  @return bool
//
func (casbinUtil *CasbinUtil) DeleteAllPolicy() bool {
	return casbinUtil.deletePolicy(0)
}

//
//  DeletePolicyForAuthorityId
//  @Description: 删除某个角色的全部权限
//  @receiver casbinUtil
//  @param authorityId
//  @return bool
//
func (casbinUtil *CasbinUtil) DeletePolicyForAuthorityId(authorityId string) bool {
	return casbinUtil.deletePolicy(0, authorityId)
}

func (casbinUtil *CasbinUtil) deletePolicy(v int, p ...string) bool {
	isDel, _ := casbinUtil.SyncedEnforcer.RemoveFilteredPolicy(v, p...)
	return isDel
}

//
//  GetAllPolicy
//  @Description: 获取所有的权限
//  @receiver casbinUtil
//
func (casbinUtil *CasbinUtil) GetAllPolicy() []dto.CasbinInfoDto {
	policyStringsLists := casbinUtil.SyncedEnforcer.GetPolicy()
	policyLists := policyStringsToListsStruct(policyStringsLists)
	return policyLists
}

//
//  GetAuthorityAllPolicy
//  @Description: 获取某个角色id的所有权限
//  @receiver casbinUtil
//  @param authorityId
//
func (casbinUtil *CasbinUtil) GetAuthorityAllPolicy(authorityId string) []dto.CasbinInfoDto {

	policyStringsLists := casbinUtil.SyncedEnforcer.GetFilteredPolicy(0, authorityId)
	policyLists := policyStringsToListsStruct(policyStringsLists)
	return policyLists
}

func (casbinUtil *CasbinUtil) Close() {
	casbinUtil.SyncedEnforcer = nil
	casbinUtil.StopAutoLoadPolicy()
}

//
//  policyStringsToListsStruct
//  @Description: [][]string  转化  []dto.CasbinInfoDto
//  @param policyStringsLists 权限数据
//  @return []dto.CasbinInfoDto
//
func policyStringsToListsStruct(policyStringsLists [][]string) []dto.CasbinInfoDto {
	policyLists := make([]dto.CasbinInfoDto, 0, 16)
	if len(policyStringsLists) > 0 {
		for _, policyStrings := range policyStringsLists {
			c1 := dto.CasbinInfoDto{}
			//list[0]
			c1.Path = policyStrings[1]
			c1.AuthorityId = policyStrings[0]
			if len(policyStrings) > 2 {
				c1.Method = policyStrings[2]
			}

			policyLists = append(policyLists, c1)
		}
	}
	return policyLists

}
