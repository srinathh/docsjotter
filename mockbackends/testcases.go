package mockbackends

import (
	"github.com/mailgun/gocql/uuid"
	"github.com/srinathh/powerdocs/structs"
	"github.com/srinathh/powerdocs/utils"
	"path/filepath"
	"regexp"
	"time"
)

/*
func GetTestCases() map[string]structs.Node {
	dec := json.NewDecoder(bytes.NewBufferString(testcasesjson))
	ret := make(map[string]structs.Node)
	dec.Decode(&ret)
	return ret
}
*/

const LOREM string = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

func GetTestCases() map[string]structs.Node {

	re := regexp.MustCompile("\\W+")
	lenlorem := len(LOREM)

	ret := make(map[string]structs.Node, len(testcases))

	for path, nodetype := range testcases {
		var n structs.Node
		n.Id = utils.DoHash(path)
		n.ModTime = time.Now().Format(utils.MyTimeStamp)
		n.Name = filepath.Base(path)
		if n.Name != "root" {
			n.Parent = utils.DoHash(filepath.Dir(path))
		} else {
			n.Parent = "#"
		}
		n.Size = "100 kb"
		n.Type = nodetype

		numcomments := len(re.FindAllString(path, -1)) + 1
		comments := make([]structs.Comment, numcomments)

		for j := 0; j < numcomments; j++ {
			var c structs.Comment
			c.Id = uuid.RandomUUID().String()
			c.ModTime = time.Now().Format(utils.MyTimeStamp)
			c.Text = LOREM[:((j+1)*lenlorem)/numcomments]
			switch j % 4 {
			case 0:
				c.Text = "#todo\n" + c.Text
			case 1:
				c.Text = "#talk\n" + c.Text
			}
			comments[j] = c
		}
		n.Comments = comments
		ret[n.Id] = n
	}
	return ret

}

var testcases = map[string]string{
	"root":                                            "directory",
	"root/Folder 1":                                   "directory",
	"root/Folder 1/SubFolder":                         "directory",
	"root/Folder 1/SubFolder/Audio File in SubFolder": "audio",
	"root/Folder 1/SubFolder/Video File in SubFolder": "video",
	"root/Folder 1/Image File in Folder 1":            "image",
	"root/Folder 1/Word File in Folder 1":             "msword",
	"root/Folder 2":                                   "directory",
	"root/Folder 2/Text File in Folder 2":             "text",
	"root/Folder 2/PowerPoint File in Folder 2":       "vnd.ms-powerpoint",
	"root/PDF File in Root":                           "pdf",
	"root/HTML File in Root":                          "html",
}
