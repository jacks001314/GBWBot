package backend

type Config struct {

	//rpc host
	RHost string `json:"rhost"`

	//rpc port
	RPort int `json:"rport"`

	//the download file dir rpc
	RDownloadDir string `json:"rDownloadDir"`

	//the ssl cert file,if this empty,then not use ssl
	CertFlag string `json:"certFlag"`

	//the ssl key file,if this empty,then not use ssl
	KeyFlag string `json:"keyFlag"`

	//the databse host
	DBHost string `json:"dbHost"`

	//the database port
	DBPort int `json:"dbPort"`

	//the database user
	DBUser string `json:"dbUser"`

	//the database password
	DBPass string `json:"dbPass"`

	//the dnslog subdomain
	Subdomain string `json:"subDomain"`

	//the ip address that dnslog query
	DnslogReplyIP string `json:"dnslogReplyIP"`

	//the cbot file store dir
	CBotFileStoreDir string `json:"cbotFileStoreDir"`

	//the attack file server dir
	AttackFileServerDir string `json:"attackFileServerDir"`

	//the attack file server host
	AttackFileServerHost string `json:"attackFileServerHost"`

	//the attack file server port
	AttackFileServerPort int `json:"attackFileServerPort"`

	//the java version used to generator jar attack payload
	JavaVersion string  `json:"javaVersion"`

	//the attack targets config file path
	AttackTargetsCFile string `json:"attackTargetsConfigPath"`

	//the attack targets wait queue capacity
	AttackTargetsQueueCapacity uint32 `json:"attackTargetsQueueCapacity"`

	//the attack targets wait time out
	AttackTargetsWaitTimeout uint64 `json:"attackTargetsWaitTimeout"`


	//the attack scripts config file path
	AttackScriptsCFile string `json:"attackScriptsConfigPath"`

	//the attack scripts wait queue capacity
	AttackScriptsQueueCapacity uint32 `json:"attackScriptsQueueCapacity"`

	//the attack scripts wait time out
	AttackScriptsWaitTimeout uint64 `json:"attackScriptsWaitTimeout"`
}
