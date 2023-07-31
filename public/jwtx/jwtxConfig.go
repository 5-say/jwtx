package jwtx

type Config struct {

	// 是否开启 IP一致性校验，当前 IP 必须与登录时 IP 一致
	CheckIP bool
	// 是否开启 单端登录
	SingleEnd bool

	// 刷新时间间隔（秒）
	//
	// 超过该时间即可触发 token 刷新、旧 token 进入容错，间隔越小安全性越高
	RefreshInterval int64
	// 并发容错时间（秒）
	FaultTolerance int64

	// 是否开启 自动续期
	//
	// 假设：单次登录有效期（7天）
	// true  只要 7天 内有操作触发刷新，登录状态就一直延续
	// false 不论 7天 内是否操作，7天 后必定要求重新登录
	AutomaticRenewal bool
	// 单次登录有效期（小时）
	AccessExpireByHour int64
}
