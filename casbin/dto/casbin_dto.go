package dto

// Casbin info structure
type CasbinInfoDto struct {
	AuthorityId string `json:"authorityId"` // 路径
	Path        string `json:"path"`        // 路径
	Method      string `json:"method"`      // 方法
}

// Casbin structure for input parameters
type CasbinInReceive struct {
	AuthorityId string          `json:"authorityId"` // 权限id
	CasbinInfos []CasbinInfoDto `json:"casbinInfos"`
}
