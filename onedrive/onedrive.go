package onedrive

import (
	"errors"
	"gonelist/conf"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 初始化登陆状态
// 如果初始化时获取失败直接退出
// 如果在自动刷新时失败给出 error 警告，见 onedrive/timer.go
func InitOnedive() {
	// 获取文件内容和初始化 README 缓存
	err := RefreshOnedriveAll()
	if err != nil {
		log.WithField("err", err).Fatal("InitOnedrive 出现错误")
	}
	// 设置 onedrive 登陆状态
	FileTree.SetLogin(true)
	cacheGoOnce.Do(func() {
		go SetAutoRefresh()
	})
}

// 刷新所有 onedrive 的内容
// 包括 文件列表，README，password，搜索索引
func RefreshOnedriveAll() error {
	log.Info("开始刷新文件缓存")
	if _, err := GetAllFiles(); err != nil { // 获取所有文件并且刷新树结构
		log.WithField("err", err).Error("刷新文件缓存遇到错误")
		return err
	}
	log.Infof("结束刷新文件缓存")
	log.Debug(FileTree)
	log.Info("开始刷新 README 缓存")
	if err := RefreshREADME(); err != nil {
		log.WithField("err", err).Error("刷新 README 缓存遇到错误")
		return err
	}
	log.Info("结束刷新 README 缓存")
	// 构建搜索
	return nil
}

// 从缓存获取某个路径下的所有内容
func CacheGetPathList(oPath string, host string) (*FileNode, error) {
	var (
		root    *FileNode
		isFound bool
	)

	root = FileTree.GetRoot()
	pArray := strings.Split(oPath, "/")

	if oPath == "" || oPath == "/" || len(pArray) < 2 {
		if conf.UserSet.DomainBasedSubFolders.Enable {
			hostSubNode, error := GetHostSpecifiedNode(root, host)
			if hostSubNode != nil {
				return ConvertReturnNode(hostSubNode, host), nil
			} else {
				return nil, error
			}
		}
		return ConvertReturnNode(root, host), nil
	}

	hostSubNode := root
	if conf.UserSet.DomainBasedSubFolders.Enable {
		subNode, error := GetHostSpecifiedNode(root, host)
		if hostSubNode == nil {
			return nil, error
		} else {
			hostSubNode = subNode
		}
	}
	for i := 1; i < len(pArray); i++ {
		isFound = false
		for _, item := range hostSubNode.Children {
			if pArray[i] == item.Name {
				hostSubNode = item
				isFound = true
			}
		}
		if isFound == false {
			log.WithFields(log.Fields{
				"oPath":    oPath,
				"pArray":   pArray,
				"position": pArray[i],
			})
			return nil, errors.New("未找到该路径")
		}
	}

	// 只返回当前层的内容
	reNode := ConvertReturnNode(hostSubNode, host)
	return reNode, nil
}

func GetHostSpecifiedNode(root *FileNode, host string) (*FileNode, error) {
	for _, pair := range conf.UserSet.DomainBasedSubFolders.Pairs {
		if host == pair.Domain {
			for _, childNode := range root.Children {
				if childNode.Path == pair.SubFolder {
					return childNode, nil
				}
			}
			return nil, errors.New("未找到站点子文件夹:" + pair.SubFolder)
		}
	}
	return nil, errors.New("站点:" + host + "未配置")
}

func ConvertReturnNode(node *FileNode, host string) *FileNode {
	if node == nil {
		return nil
	}

	reNode := CopyFileNode(node, host)

	if reNode == nil {
		return nil
	}
	for key := range node.Children {
		if node.Children[key].Name == ".password" {
			continue
		}
		tmpNode := node.Children[key]
		reNode.Children = append(reNode.Children, CopyFileNode(tmpNode, host))
	}
	return reNode
}

func CopyFileNode(node *FileNode, host string) *FileNode {
	if node == nil {
		return nil
	}
	path := GetReplacePath(node.Path, host)
	return &FileNode{
		Name:           node.Name,
		Path:           path,
		IsFolder:       node.IsFolder,
		DownloadUrl:    node.DownloadUrl,
		LastModifyTime: node.LastModifyTime,
		Size:           node.Size,
		Children:       nil,
		Password:       node.Password,
	}
}

func GetDownloadUrl(filePath string, host string) (string, error) {
	var (
		fileInfo    *FileNode
		err         error
		downloadUrl string
	)

	if fileInfo, err = CacheGetPathList(filePath, host); err != nil || fileInfo == nil || fileInfo.IsFolder == true {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"err":      err,
		}).Info("请求的文件未找到")
		return "", err
	}

	// 如果有重定向前缀，就加上
	downloadUrl = conf.UserSet.DownloadRedirectPrefix + fileInfo.DownloadUrl

	return downloadUrl, nil
}

// 替换路径
// 如果设置了 folderSub 为 /public
// 那么 /public 替换为 /, /public/test 替换为 /test
func GetReplacePath(pSrc string, host string) string {
	if conf.UserSet.DomainBasedSubFolders.Enable {
		for i := range conf.UserSet.DomainBasedSubFolders.Pairs {
			pair := conf.UserSet.DomainBasedSubFolders.Pairs[i]
			if pair.Domain == host {
				return ReplaceLeadingPath(pSrc, pair.SubFolder)
			}
		}
		return ReplaceLeadingPath(pSrc, conf.UserSet.DomainBasedSubFolders.DefaultFolder)
	} else if conf.UserSet.Server.FolderSub != "/" {
		return ReplaceLeadingPath(pSrc, conf.UserSet.Server.FolderSub)
	}
	return pSrc
}

func ReplaceLeadingPath(source string, target string) string {
	if strings.Index(source, target) == 0 {
		str := strings.Replace(source, target, "", 1)
		if str == "" {
			return "/"
		} else {
			return str
		}
	} else {
		return source
	}

}
