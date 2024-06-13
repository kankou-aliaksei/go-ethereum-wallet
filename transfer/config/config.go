package config

type Config struct {
	PublicNodeUrl       string
	EthereumExplorerUrl string
	UsdtContractAddress string
}

var EthereumMainnet = Config{
	PublicNodeUrl:       "https://cloudflare-eth.com",
	EthereumExplorerUrl: "https://etherscan.io",
	UsdtContractAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
}

var SepoliaTestnet = Config{
	PublicNodeUrl:       "https://rpc.sepolia.org",
	EthereumExplorerUrl: "https://sepolia.etherscan.io",
	UsdtContractAddress: "0xE3d2B274Ec5a0F4e9FA12911F76BA052faFeA6aE",
}
