package command

import (
	"encoding/base64"
	"errors"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/consul-go/api"
	"net/http"
	"os"
	"strings"
	"time"
)

func KvExportCommand() *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: "导出",
		Flags: append(Common(), ExportFolderFlag()),
		Action: func(context *cli.Context) error {
			var baseUrl = context.String(BaseUrl)
			var dc = context.String(Dc)
			var exportFolder = context.String(ExportFolder)

			currentTime := time.Now()
			folderTime := currentTime.Format("2006-01-02_15-04-05")

			return KvExport(baseUrl, dc, exportFolder+"/"+folderTime)
		},
	}
}

func KvExport(baseUrl string, dc string, exportFolder string) error {
	client, err := api.NewClient(baseUrl, "", "")
	if err != nil {
		return err
	}

	_, err = os.Stat(exportFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(exportFolder, 0755)
		if err != nil {
			return err
		}
	}

	var getKvRequestQuery = &api.GetKvRequestQuery{
		Keys:      "",
		Dc:        dc,
		Separator: "/",
	}

	contents, response, err := client.Kv.GetKv("", getKvRequestQuery)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("响应状态码: " + response.Status + " 不正常")
	}

	for _, name := range contents {

		if strings.HasSuffix(name, "/") {
			err = KvFolder(dc, name, exportFolder, client)
			if err != nil {
				return err
			}
		} else {
			err = KvGetKvName(dc, name, exportFolder, client)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func KvFolder(dc string, path string, exportFolder string, client *api.Client) error {
	var getKvRequestQuery = &api.GetKvRequestQuery{
		Keys:      "",
		Dc:        dc,
		Separator: "/",
	}

	_, err := os.Stat(exportFolder + "/" + path)
	if os.IsNotExist(err) {
		err := os.Mkdir(exportFolder+"/"+path, 0755)
		if err != nil {
			return err
		}
	}

	contents, response, err := client.Kv.GetKv(path, getKvRequestQuery)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("响应状态码: " + response.Status + " 不正常")
	}

	for _, name := range contents {

		if name == path {
			continue
		}

		if strings.HasSuffix(name, "/") {
			err = KvFolder(dc, name, exportFolder, client)
			if err != nil {
				return err
			}
		} else {
			err = KvGetKvName(dc, name, exportFolder, client)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func KvGetKvName(dc string, name string, exportFolder string, client *api.Client) error {
	var getKvNameRequestQuery = &api.GetKvNameRequestQuery{
		Dc: dc,
	}
	responses, response, err := client.Kv.GetKvName(name, getKvNameRequestQuery)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("响应状态码: " + response.Status + " 不正常")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(responses[0].Value)
	if err != nil {
		return err
	}

	decodedString := string(decodedBytes)

	err = os.WriteFile(exportFolder+"/"+name, []byte(decodedString), 0644)
	if err != nil {
		return err
	}

	return nil
}
