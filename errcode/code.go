package errcode

const (
	ResNotFound             = 10001 // 未找到
	ResExisted              = 10002 // 资源已存在
	PermissionDenied        = 10003 // 权限不足
	ParamError              = 10004 // web参数错误
	InvalidSecrets          = 10005 // 密码错误
	Unauthorized            = 10006 // 未登录
	LastAuthorization       = 10007 // 最后一个认证类型
	RateLimit               = 10008 // 频率限制
	ResNotMatch             = 10009 // 资源未能匹配
	InvalidContentSignature = 10010 // 内容签名错误（不是接口签名）
	InvalidWechatSession    = 10011 // 微信session无效
	ResNotAvailable         = 10012 // 资源暂不可用
	ResBusy                 = 10013 // 资源繁忙（在使用）
	ResUsedTooManyTimes     = 10014 // 资源被使用资数过多
	ResExpired              = 10015 // 资源过期
	LoginDeviceLimit        = 10016 // 登录设备超出限制

	RemoteServerError           = 20001 // 第三方服务错误
	QiNiuRemoteServerError      = 20002 // 七牛服务错误
	AppStoreRemoveServerError   = 20003 // AppStore远程错误
	RemoveServerTimeout         = 20004 // 远程服务器访问超时
	DecodeJSONError             = 20005 // 解析JSON失败
	InvalidXML                  = 20006 // 无效的XML
	WechatPayError              = 20007 // 解析JSON失败
	AlipayAmountError           = 20008 // 支付宝金额错误
	UnexpectRemoteResponse      = 20009 // 远程服务器格式非预期
	InvalidRemoteSessionKey     = 20010 // 远程服务器SessionKey无效
	UnAuthorizedAppStoreReceipt = 20011 // 未授权的苹果收据
	InvalidHuaweiPaymentData    = 20012 // 无效的华为交易数据
	NotProvisioned              = 55555 // 非预期
	Undefined                   = 99999 // 未定义
)
