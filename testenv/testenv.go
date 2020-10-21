package testenv

import (
	"os"

	"github.com/chenjie4255/env"
	"github.com/joho/godotenv"
)

type IntegratedTestEnv struct {
	MongoHost                     string `env:"TEST_MONGO_HOST"`
	MongoUsername                 string `env:"TEST_MONGO_USERNAME"`
	MongoPassword                 string `env:"TEST_MONGO_PASSWORD"`
	MongoSource                   string `env:"TEST_MONGO_SOURCE"`
	ElasticHost                   string `env:"TEST_ELASTIC_HOST"`
	LeanCloudID                   string `env:"TEST_LEAN_CLOUD_ID"`
	LeanCloudKey                  string `env:"TEST_LEAN_CLOUD_KEY"`
	LeanCloudSyncFlagObject       string `env:"TEST_LEAN_CLOUD_SYNC_FLAG_OBJECT"`
	LeanCloudSyncFlagTimingObject string `env:"TEST_LEAN_CLOUD_SYNC_FLAG_TIMING_OBJECT"`
	LeanCloudPWDUsername          string `env:"TEST_LEAN_CLOUD_PWD_USERNAME"`
	LeanCloudPWDPassword          string `env:"TEST_LEAN_CLOUD_PWD_PASSWORD"`
	RedisHost                     string `env:"TEST_REDIS_HOST"`
	RedisPassword                 string `env:"TEST_REDIS_PASSWORD"`
	WeChatOpenID                  string `env:"TEST_WECHAT_OEPN_ID"`
	WeChatToken                   string `env:"TEST_WECHAT_TOKEN"`
	WeChatAppID                   string `env:"TEST_WECHAT_APP_ID"`
	WeChatPayAppKey               string `env:"TEST_WECHAT_PAY_APP_KEY"`
	WeChatPayMchID                string `env:"TEST_WECHAT_PAY_MCH_ID"`
	WeChatSecretKey               string `env:"TEST_WECHAT_SECRET_KEY"`
	WeChatCode                    string `env:"TEST_WECHAT_CODE"`
	FacebookToken                 string `env:"TEST_FACEBOOK_TOKEN"`
	GoogleToken                   string `env:"TEST_GOOGLE_TOKEN"`
	WebiBoToken                   string `env:"TEST_WEIBO_TOKEN"`
	QQAppID                       string `env:"TEST_QQ_APP_ID"`
	QQToken                       string `env:"TEST_QQ_TOKEN"`
	CMQAppID                      string `env:"TEST_CMQ_APP_ID"`
	CMQAppKey                     string `env:"TEST_CMQ_APP_KEY"`
	AppStoreReceipt               string `env:"TEST_APPSTORE_RECEIPT"`
	AppStoreKey                   string `env:"TEST_APPSTORE_KEY"`
	BearyChatURI                  string `env:"BEARY_CHAT_URI"`
	WechatMPAppID                 string `env:"TEST_WECHAT_MP_APP_ID"`
	WechatMPAppKey                string `env:"TEST_WECHAT_MP_APP_KEY"`
	QiNiuAccessID                 string `env:"TEST_QINIU_ACCESS_ID"`
	QiNiuAccessKey                string `env:"TEST_QINIU_ACCESS_KEY"`
	UMengAppKey                   string `env:"TEST_UMENG_APP_KEY"`
	UMengAppSecret                string `env:"TEST_UMENG_APP_SECRET"`
	UMengIOSAppKey                string `env:"TEST_UMENG_IOS_APP_KEY"`
	UMengIOSAppSecret             string `env:"TEST_UMENG_IOS_APP_SECRET"`
	UMengDeviceToken              string `env:"TEST_UMENG_DEVICE_TOKEN"`
	UMengAndroidDeviceToken       string `env:"TEST_UMENG_ANDROID_DEVICE_TOKEN"`

	// QQ MP
	QQMPAppID  string `env:"TEST_QQ_MP_APP_ID"`
	QQMPAppKey string `env:"TEST_QQ_MP_APP_KEY"`

	// AliPay
	AlipayAppID        string `env:"TEST_ALIPAY_APP_ID"`
	AlipayAliPublicKey string `env:"TEST_ALIPAY_ALI_PUBLIC_KEY"`
	AlipayPrivateKey   string `env:"TEST_ALIPAY_PRIVATE_KEY"`

	// Android Play Subscription
	AndroidPlayCredJSON          string `env:"TEST_ANDROID_PLAY_CRED_JSON"`
	AndroidPlaySubscriptionToken string `env:"TEST_ANDROID_PLAY_SUBSCRIPTION_TOKEN"`
	AndroidPlayPurchaseToken     string `env:"TEST_ANDROID_PLAY_PURCHASE_TOKEN"`

	// HUAWEI
	HuaWeiAccessToken string `env:"TEST_HUAWEI_ACCESS_TOKEN"`
	HuaWeiOpenID      string `env:"TEST_HUAWEI_OPEN_ID"`
	HuaweiPrivateKey  string `env:"TEST_HUAWEI_PRIVATE_KEY"`
	HuaweiMerchantID  string `env:"TEST_HUAWEI_MERCHANT_ID"`
	HuaweiAppID       string `env:"TEST_HUAWEI_APP_ID"`
	HuaweiAppKey      string `env:"TEST_HUAWEI_APP_KEY"`
	HuaweiIDToken     string `env:"TEST_HUAWEI_ID_TOKEN"`

	// baidu
	BaiduClientID     string `env:"TEST_BAIDU_CLIENT_ID"`
	BaiduClientSecret string `env:"TEST_BAIDU_CLIENT_SECRET"`

	// Google Play
	GooglePlayPublicKey string `env:"GOOGLE_PLAY_PUBLIC_KEY"`

	FacebookAppID     string `env:"FACEBOOK_APP_ID"`
	FacebookAppSecret string `env:"FACEBOOK_APP_SECRET"`
}

// GetIntegratedTestEnv 获取集成测试用环境变量
func GetIntegratedTestEnv() *IntegratedTestEnv {
	testEnv := IntegratedTestEnv{}
	Parse(&testEnv)
	return &testEnv
}

// ParseEnv ParseEnvFile
func Parse(output interface{}, filenames ...string) error {
	envFile := os.Getenv("GOTEST_ENV_FILE")
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return err
		}
	}

	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil {
			return err
		}
	}

	return env.Parse(output)
}
