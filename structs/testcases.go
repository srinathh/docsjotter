package structs

import (
	"fmt"
	"github.com/srinathh/powerdocs/utils"
	"math/rand"
	"time"
)

const LOREM string = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

func GenComments(maxcomments int) []Comment {

	c := make([]Comment, rand.Intn(maxcomments))

	for i, _ := range c {
		c[i].Text = LOREM[:rand.Intn(len(LOREM))]
		//adding current nanosecond to prevent comments with same text getting clashing ids
		c[i].Id = utils.DoHash(fmt.Sprintf("%s%d", c[i].Text, time.Now().UnixNano()))
		c[i].ModTime = time.Now().Format(time.Stamp)
	}
	return c
}

func GenNodes() map[string]Node {

	Nodes := make(map[string]Node)

	Nodes["0"] = Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["1"] = Node{
		Id:       "1",
		Name:     "Folder 1",
		Type:     "directory",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["2"] = Node{
		Id:       "2",
		Name:     "SubFolder",
		Type:     "directory",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["3"] = Node{
		Id:       "3",
		Name:     "Audio File in SubFolder",
		Type:     "audio",
		Parent:   "2",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["4"] = Node{
		Id:       "4",
		Name:     "Video File in SubFolder",
		Type:     "video",
		Parent:   "2",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["5"] = Node{
		Id:       "5",
		Name:     "Image File in Folder 1",
		Type:     "image",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["6"] = Node{
		Id:       "6",
		Name:     "Word File in Folder 1",
		Type:     "msword",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["7"] = Node{
		Id:       "7",
		Name:     "Folder 2",
		Type:     "directory",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["8"] = Node{
		Id:       "8",
		Name:     "Text File in Folder 2",
		Type:     "text",
		Parent:   "7",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["9"] = Node{
		Id:       "9",
		Name:     "PowerPoint File in Folder 2",
		Type:     "vnd.ms-powerpoint",
		Parent:   "7",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["10"] = Node{
		Id:       "10",
		Name:     "PDF File in Root",
		Type:     "pdf",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes["11"] = Node{
		Id:       "11",
		Name:     "HTML File in Root",
		Type:     "html",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	return Nodes
}
