package command

import "github.com/urfave/cli/v2"

const (
	BaseUrl          = "base-url"
	Dc               = "dc"
	ExportFolder     = "export-folder"
	ExportFolderDate = "export-folder-date"
	ImportFolder     = "import-folder"
)

func BaseUrlFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    BaseUrl,
		Usage:   "Consul URL",
		Value:   "http://127.0.0.1:8500",
		EnvVars: []string{"CONSUL_GO_BASE_URL"},
	}
}

func DcFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    Dc,
		Usage:   "dc",
		Value:   "dc1",
		EnvVars: []string{"CONSUL_GO_DC"},
	}
}

func ExportFolderFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    ExportFolder,
		Usage:   "导出文件夹",
		Value:   "consul-go-export-folder",
		EnvVars: []string{"CONSUL_GO_EXPORT_FOLDER"},
	}
}

func ExportFolderUseDateFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    ExportFolderDate,
		Usage:   "导出文件夹是否使用日期",
		Value:   true,
		EnvVars: []string{"CONSUL_GO_EXPORT_FOLDER_DATE"},
	}
}

func ImportFolderFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    ImportFolder,
		Usage:   "导入文件夹",
		Value:   "consul-go-import-folder",
		EnvVars: []string{"CONSUL_GO_IMPORT_FOLDER"},
	}
}

func Common() []cli.Flag {
	return []cli.Flag{
		BaseUrlFlag(),
		DcFlag(),
	}
}
