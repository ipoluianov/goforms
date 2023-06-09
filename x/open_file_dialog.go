package x

/*
import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/goforms/paths"
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type OpenFileDialog struct {
	uicontrols.Dialog

	panelNavButtons *uicontrols.Panel
	panelTree       *uicontrols.Panel
	panelResult     *uicontrols.Panel
	panelButtons    *uicontrols.Panel

	btnRoot      *uicontrols.Button
	btnHome      *uicontrols.Button
	btnDocuments *uicontrols.Button
	btnDownloads *uicontrols.Button
	btnDesktop   *uicontrols.Button
	btnPictures  *uicontrols.Button

	tvItems *uicontrols.TreeView

	txtFilePath *uicontrols.TextBox

	filePath         string
	defaultDirectory string
}

func NewOpenFileDialog(parent uiinterfaces.Widget, defaultDirectory string) *OpenFileDialog {
	var c OpenFileDialog
	c.InitControl(parent, &c)
	c.SetTitle("Open file")
	c.Resize(600, 750)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.panelNavButtons = pContent.AddPanelOnGrid(0, 0)
	c.panelTree = pContent.AddPanelOnGrid(0, 1)
	c.panelResult = pContent.AddPanelOnGrid(0, 2)
	c.panelButtons = pContent.AddPanelOnGrid(0, 3)

	// NavButtons
	c.btnRoot = c.panelNavButtons.AddButtonOnGrid(0, 0, "Root", func(event *uievents.Event) {
		c.SetRoot("")
	})
	c.btnHome = c.panelNavButtons.AddButtonOnGrid(1, 0, "Home", func(event *uievents.Event) {
		c.SetRoot(paths.HomeFolder())
	})
	c.btnDocuments = c.panelNavButtons.AddButtonOnGrid(2, 0, "Documents", func(event *uievents.Event) {
		c.SetRoot(paths.DocumentsFolder())
	})
	c.btnDownloads = c.panelNavButtons.AddButtonOnGrid(0, 1, "Downloads", func(event *uievents.Event) {
		c.SetRoot(paths.DownloadsFolder())
	})
	c.btnDesktop = c.panelNavButtons.AddButtonOnGrid(1, 1, "Desktop", func(event *uievents.Event) {
		c.SetRoot(paths.DesktopFolder())
	})
	c.btnPictures = c.panelNavButtons.AddButtonOnGrid(2, 1, "Pictures", func(event *uievents.Event) {
		c.SetRoot(paths.PicturesFolder())
	})

	// Tree
	c.tvItems = c.panelTree.AddTreeViewOnGrid(0, 0)
	c.tvItems.AddColumn("Ext", 70)
	c.tvItems.AddColumn("Size", 70)
	//c.tvItems.AddColumn("Date/Time", 150)
	c.tvItems.SetColumnWidth(0, 400)

	c.tvItems.OnExpand = func(treeView *uicontrols.TreeView, node *uicontrols.TreeNode) {
		c.loadNode(node)
	}
	c.tvItems.OnSelectedNode = func(treeView *uicontrols.TreeView, node *uicontrols.TreeNode) {
		fi, ok := node.UserData.(*FileInfo)
		if ok {
			c.txtFilePath.SetText(strings.ReplaceAll(fi.Path, "//", "/"))
		}
	}

	// Result
	c.txtFilePath = c.panelResult.AddTextBoxOnGrid(0, 0)
	c.txtFilePath.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.filePath = newValue
	}
	// Buttons
	c.panelButtons.AddHSpacerOnGrid(0, 0)
	btnOK := c.panelButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		return true
	}
	btnOK.SetMinWidth(70)
	btnCancel := c.panelButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)
	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	c.btnRoot.Press()

	return &c
}

func (c *OpenFileDialog) SetRoot(path string) {
	c.tvItems.RemoveAllNodes()

	if path == "" {
		rootItems := GetRootItems()
		rootNode := c.tvItems.AddNode(nil, "Local Computer")
		rootNode.UserData = "root"
		//rootNode.Icon = uiresources.ResImageAdjusted("icons/material/file/drawable-hdpi/baseline_folder_black_48dp.png", c.ForeColor())

		for _, rootItem := range rootItems {
			node := c.tvItems.AddNode(rootNode, rootItem)
			var fi FileInfo
			fi.Name = rootItem
			fi.Dir = true
			fi.Path = rootItem
			node.UserData = &fi
			//node.Icon = uiresources.ResImageAdjusted("icons/material/file/drawable-hdpi/baseline_folder_open_black_48dp.png", c.ForeColor())
			c.tvItems.AddNode(node, "loading ...")
		}
		c.tvItems.ExpandNode(rootNode)
	} else {
		var fi FileInfo
		fi.Name = path
		fi.Dir = true
		fi.Path = path
		rootNode := c.tvItems.AddNode(nil, path)
		rootNode.UserData = &fi
		//rootNode.Icon = uiresources.ResImageAdjusted("icons/material/file/drawable-hdpi/baseline_folder_black_48dp.png", c.ForeColor())
		c.tvItems.ExpandNode(rootNode)
	}

}

func (c *OpenFileDialog) loadNode(loadingNode *uicontrols.TreeNode) {
	nodeData, ok := loadingNode.UserData.(string)
	if ok {
		if nodeData == "root" {
			return
		}
	}

	c.tvItems.RemoveNodes(loadingNode)
	fi, ok := loadingNode.UserData.(*FileInfo)
	if ok {
		fis, err := GetDir(fi.Path)
		if err == nil {
			for _, f := range fis {
				n := c.tvItems.AddNode(loadingNode, f.Name)
				if !f.Dir {
					c.tvItems.SetNodeValue(n, 0, f.NameWithoutExt)
					c.tvItems.SetNodeValue(n, 1, f.Ext)
				}
				c.tvItems.SetNodeValue(n, 2, f.SizeAsString())
				ff := f
				n.UserData = &ff
			}
		}
	}
}

type FileInfo struct {
	Path           string
	Name           string
	Dir            bool
	NameWithoutExt string
	Ext            string
	Size           int64
	Date           time.Time
	Attr           string
}

func (c *FileInfo) SizeAsString() string {
	if c.Dir {
		return "<DIR>"
	}
	div := int64(1)
	uom := ""
	if c.Size >= 1024 {
		div = 1024
		uom = "KB"
	}
	if c.Size >= 1024*1024 {
		div = 1024 * 1024
		uom = "MB"
	}
	if c.Size >= 1024*1024*1024 {
		div = 1024 * 1024 * 1024
		uom = "GB"
	}

	result := strconv.FormatInt(c.Size/div, 10) + " " + uom
	return result
}

func GetDir(path string) ([]FileInfo, error) {
	result := make([]FileInfo, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result, err
	}

	dirsList := make([]FileInfo, 0)
	filesList := make([]FileInfo, 0)

	for _, f := range files {
		var fileInfo FileInfo
		if f.IsDir() {
			fileInfo.Name = f.Name()
		} else {
			fileInfo.Name = f.Name()
			fileInfo.Ext = filepath.Ext(f.Name())
			if len(fileInfo.Ext) > 0 {
				fileInfo.Ext = fileInfo.Ext[1:]
			}
		}

		fileInfo.NameWithoutExt = strings.TrimSuffix(fileInfo.Name, filepath.Ext(fileInfo.Name))
		fileInfo.Path = path + "/" + f.Name()
		fileInfo.Date = f.ModTime()
		fileInfo.Dir = f.IsDir()
		fileInfo.Size = f.Size()

		if f.IsDir() {
			dirsList = append(dirsList, fileInfo)
		} else {
			filesList = append(filesList, fileInfo)
		}
	}

	sort.Slice(dirsList, func(i, j int) bool {
		return dirsList[i].Name < dirsList[j].Name
	})

	sort.Slice(filesList, func(i, j int) bool {
		return filesList[i].Name < filesList[j].Name
	})

	for _, d := range dirsList {
		result = append(result, d)
	}

	for _, f := range filesList {
		result = append(result, f)
	}

	return result, nil
}

func ShowOpenFile(parent uiinterfaces.Widget, selected func(filePath string)) {
	dialog := NewOpenFileDialog(parent, "")
	dialog.ShowDialog()
	dialog.OnAccept = func() {
		if selected != nil {
			selected(dialog.filePath)
		}
	}
}
*/
