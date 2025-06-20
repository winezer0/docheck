package generate

import (
	"cdnCheck/cdncheck"
	"cdnCheck/fileutils"
	"cdnCheck/maputils"
	"fmt"
)

// TransferNaliCdnYaml  实现nali cdn.yml到json的转换
func TransferNaliCdnYaml(path string) *cdncheck.CDNData {
	// 数据来源 https://github.com/4ft35t/cdn/blob/master/src/cdn.yml
	// CloudKeysData 是整个 YAML 文件的结构
	type naliCdnData map[string]struct {
		Name string `yaml:"name"`
		Link string `yaml:"link"`
	}
	// 1. 读取 YAML 到 CloudKeysData
	var yamlData naliCdnData
	err := fileutils.ReadYamlToStruct(path, &yamlData)
	if err != nil {
		panic(err)
	}

	// 2. 构建 cname map[string][]string 并赋值给 cdnData.CDN.CNAME
	// 初始化 CDNData 结构
	cdnData := cdncheck.NewEmptyCDNData()
	for domain, info := range yamlData {
		cdnData.CDN.CNAME[info.Name] = append(cdnData.CDN.CNAME[info.Name], domain)
	}

	return cdnData
}

// TransferPDCdnCheckJson 实现cdn check json 数据源的转换
func TransferPDCdnCheckJson(path string) *cdncheck.CDNData {
	// PDCdnCheckData 表示整个配置结构
	type PDCdnCheckData struct {
		CDN    map[string][]string `json:"cdn"`
		WAF    map[string][]string `json:"waf"`
		Cloud  map[string][]string `json:"cloud"`
		Common map[string][]string `json:"common"`
	}

	// 加载cdn check json数据源
	var pdCdnCheckData PDCdnCheckData
	err := fileutils.ReadJsonToStruct(path, &pdCdnCheckData)
	if err != nil {
		panic(err)
	}

	// 将 cdn/waf/cloud 的值作为 IP 数据填充到对应字段
	cdnData := cdncheck.NewEmptyCDNData()
	cdnData.CDN.IP = maputils.CopyMap(pdCdnCheckData.CDN)
	cdnData.WAF.IP = maputils.CopyMap(pdCdnCheckData.WAF)
	cdnData.CLOUD.IP = maputils.CopyMap(pdCdnCheckData.Cloud)

	// 合并 common 到 cdn.cname
	for provider, cnames := range pdCdnCheckData.Common {
		cdnData.CDN.CNAME[provider] = append([]string{}, cnames...)
	}
	return cdnData
}

// TransferCloudKeysYaml  实现 cloud keys yml到json的转换
func TransferCloudKeysYaml(path string) *cdncheck.CDNData {
	// 数据来源 用户自己数据到cloud_keys.yml中
	// 是整个 YAML 文件的结构
	var cloudKeysYaml map[string]struct {
		Keys []string `yaml:"keys"`
	}

	// 1. 读取 YAML 到 CloudKeysData
	err := fileutils.ReadYamlToStruct(path, &cloudKeysYaml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("yamlData: %v", cloudKeysYaml)
	//
	// 2. 构建 cname map[string][]string 并赋值给 cdnData.CDN.CNAME
	// 初始化 CDNData 结构
	cdnData := cdncheck.NewEmptyCDNData()
	for cloudName, yamEntry := range cloudKeysYaml {
		cdnData.CLOUD.KEYS[cloudName] = yamEntry.Keys
	}

	return cdnData
}
