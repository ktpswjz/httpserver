package packer

import (
	"crypto/rand"
	"fmt"
	"github.com/ktpswjz/httpserver/tool/deploy/pack/config"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Packer interface {
	Pack() error
	OutputFolder() string
}

func NewPacker(cfg *config.Config) Packer {
	return &packer{cfg: cfg}
}

type packer struct {
	cfg *config.Config
}

func (s *packer) Pack() error {
	if s.cfg == nil {
		return fmt.Errorf("配置无效：为空")
	}
	if len(s.cfg.Version) < 1 {
		return fmt.Errorf("版本号无效：为空")
	}

	outRootPath := s.OutputFolder()
	err := os.RemoveAll(outRootPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outRootPath, 0777)
	if err != nil {
		return err
	}

	err = s.packApp(outRootPath)
	if err != nil {
		return err
	}

	err = s.packWeb(outRootPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *packer) OutputFolder() string {
	return filepath.Join(s.cfg.Destination, s.cfg.Version)
}

func (s *packer) packApp(outRootPath string) error {
	binaryFileName := fmt.Sprintf("%s_rel_%s_%s_%s.%s", s.cfg.App.Name, runtime.GOOS, runtime.GOARCH, s.cfg.Version, s.pkgExt())
	fmt.Println("正在打包服务程序:", binaryFileName)

	srcFolder := s.cfg.App.Bin.Root
	_, err := os.Stat(srcFolder)
	if os.IsNotExist(err) {
		return err
	}
	if len(s.cfg.App.Bin.Files) < 1 {
		return fmt.Errorf("未指定发布文件")
	}
	tmpFolderName := s.newGuid()
	tmpFolderPath := filepath.Join(outRootPath, tmpFolderName)
	err = os.MkdirAll(tmpFolderPath, 0777)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpFolderPath)
	binFolderPath := filepath.Join(tmpFolderPath, "bin")
	err = os.MkdirAll(binFolderPath, 0777)
	if err != nil {
		return err
	}

	for srcName, destName := range s.cfg.App.Bin.Files {
		srcPath := filepath.Join(srcFolder, srcName)
		fi, err := os.Stat(srcPath)
		if os.IsNotExist(err) {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("指定的文件'%s'是个文件夹", srcName)
		}
		distPath := filepath.Join(binFolderPath, srcPath)
		if destName != "" {
			distPath = filepath.Join(binFolderPath, destName)
		}
		_, err = s.copyFile(srcPath, distPath)
		if err != nil {
			return err
		}
	}

	siteRoot := filepath.Join(tmpFolderPath, "site")
	folder := &Folder{}
	webCount := len(s.cfg.Webs)
	for webIndex := 0; webIndex < webCount; webIndex++ {
		web := s.cfg.Webs[webIndex]
		if !web.Enable {
			continue
		}

		err = folder.Copy(filepath.Join(web.Src.Root, "dist"), filepath.Join(siteRoot, web.Name))
		if err != nil {
			return err
		}
	}

	binaryFile, err := os.Create(filepath.Join(outRootPath, binaryFileName))
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	err = s.compressFolder(binaryFile, tmpFolderPath, "", nil)
	if err != nil {
		return err
	}

	// source
	if s.cfg.Source {
		sourceFileName := fmt.Sprintf("%s_src_%s.%s", s.cfg.App.Name, s.cfg.Version, s.pkgExt())
		fmt.Println("正在打包服务源代码:", sourceFileName)

		sourcePath := s.cfg.App.Src.Root
		_, err := os.Stat(sourcePath)
		if os.IsNotExist(err) {
			return err
		}

		sourceFile, err := os.Create(filepath.Join(outRootPath, sourceFileName))
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		err = s.compressFolder(sourceFile, sourcePath, "", s.cfg.App.Src.IsIgnore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *packer) packWeb(outRootPath string) error {
	count := len(s.cfg.Webs)
	for i := 0; i < count; i++ {
		web := s.cfg.Webs[i]
		if !web.Enable {
			continue
		}

		err := s.outputWeb(outRootPath, web)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *packer) outputWeb(outRootPath string, web config.Web) error {
	binaryFileName := fmt.Sprintf("web.%s_rel_%s.%s", web.Name, s.cfg.Version, s.pkgExt())
	fmt.Println("正在打包网站程序:", binaryFileName)
	err := s.outputWebFolder(filepath.Join(web.Src.Root, "dist"), filepath.Join(outRootPath, binaryFileName), nil)
	if err != nil {
		return err
	}

	if s.cfg.Source {
		binaryFileName = fmt.Sprintf("web.%s_src_%s.%s", web.Name, s.cfg.Version, s.pkgExt())
		fmt.Println("正在打包网站源代码:", binaryFileName)
		err = s.outputWebFolder(web.Src.Root, filepath.Join(outRootPath, binaryFileName), web.Src.IsIgnore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *packer) outputWebFolder(folderPath, filePath string, ignore func(name string) bool) error {
	fi, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("指定的文件夹'%s'是个文件", folderPath)
	}

	binaryFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	return s.compressFolder(binaryFile, folderPath, "", ignore)
}

func (s *packer) copyFile(source, dest string) (int64, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return 0, err
	}

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileInfo.Mode())
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, sourceFile)
}

func (s *packer) newGuid() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
