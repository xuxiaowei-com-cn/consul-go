package command

import (
	"errors"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/consul-go/api"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func KvImportCommand() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "导入",
		Flags: append(Common(), ImportFolderFlag()),
		Action: func(context *cli.Context) error {
			var baseUrl = context.String(BaseUrl)
			var dc = context.String(Dc)
			var importFolder = context.String(ImportFolder)

			return ImExport(baseUrl, dc, importFolder)
		},
	}
}

func ImExport(baseUrl string, dc string, importFolder string) error {

	importFolder = strings.ReplaceAll(importFolder, "\\", "/")
	if !strings.HasSuffix(importFolder, "/") {
		importFolder += "/"
	}

	_, err := os.Stat(importFolder)
	if os.IsNotExist(err) {
		return errors.New("文件夹：" + importFolder + " 不存在")
	}

	err = filepath.Walk(importFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是文件，而不是文件夹，则读取文件内容
		if !info.IsDir() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			tmp1 := strings.ReplaceAll(path, "\\", "/")
			tmp2 := strings.Replace(tmp1, importFolder, "", 1)

			var requestBody = string(fileContent)

			err = ImExportPutKvName(baseUrl, dc, tmp2, requestBody)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ImExportPutKvName(baseURL string, dc string, name string, requestBody string) error {

	client, err := api.NewClient(baseURL, "", "")
	if err != nil {
		return err
	}

	var putKvNameRequestQuery = &api.PutKvNameRequestQuery{
		Dc: dc,
	}

	result, response, err := client.Kv.PutKvName(name, putKvNameRequestQuery, requestBody)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("响应状态码: " + response.Status + " 不正常")
	}

	if !*result {
		return errors.New("创建 Key / Value: " + name + " 状态码: " + response.Status + " 不正常")
	}

	return nil
}
