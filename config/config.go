package config

// 存储类型(表示文件存到哪里)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreOSS : 阿里OSS
	StoreOSS
	// StoreMix : 混合(Ceph及OSS)
	StoreMix
	// StoreAll : 所有类型的存储都存一份数据
	StoreAll
)
const (
	UploadServiceHost = "0.0.0.0:8088" // UploadServiceHost : 上传服务监听的地址

	MySQLDSN = "root:123456@tcp(127.0.0.1:3306)/file2?charset=utf8&parseTime=true"

	RedisHost = "127.0.0.1:6379"
	RedisPass = ""

	CurrentStoreType = StoreLocal // 设置当前文件的存储类型

	AsyncTransferEnable  = true                                 // 是否开启文件异步转移(默认同步)
	RabbitURL            = "amqp://guest:guest@127.0.0.1:5672/" // rabbitmq服务的入口url
	TransExchangeName    = "uploadserver.trans"                 // 用于文件transfer的交换机
	TransOSSQueueName    = "uploadserver.trans.oss"             // oss转移队列名
	TransOSSErrQueueName = "uploadserver.trans.oss.err"         // oss转移失败后写入另一个队列的队列名
	TransOSSRoutingKey   = "oss"                                // routingkey

	OSSBucket          = "buckettest-filestore-miao"    // oss bucket名
	OSSEndpoint        = "oss-cn-shenzhen.aliyuncs.com" // oss endpoint
	OSSAccesskeyID     = ""                             // oss访问key
	OSSAccessKeySecret = ""                             // oss访问key secret

	CephAccessKey  = "8WOOFAOAZ3SKQK3Y5I2L"                     // 访问Key
	CephSecretKey  = "syYWcEmF0Dx7BXrpyDvuAZ3yRe4EmNC9oDrucx3M" // 访问密钥
	CephGWEndpoint = "http://127.0.0.1:9080"                    // gateway地址

	UserPwdSalt = "#666#"
)
