package export

import "github.com/Asolmn/go-gin-example/pkg/setting"

// 获取excel文件的url地址
func GetExcelFullUrl(name string) string {
	// http://127.0.0.1:8000/export/name
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

// excel的相对地址
func GetExcelPath() string {
	// export/
	return setting.AppSetting.ExportSavePath
}

// excel的绝对地址
func GetExcelFullPath() string {
	// runtime/export/
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
