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

	Nodes[utils.DoHash("root")] = Node{
		Id:       utils.DoHash("root"),
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1")] = Node{
		Id:       utils.DoHash("root/Folder 1"),
		Name:     "Folder 1",
		Type:     "directory",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder")] = Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder"),
		Name:     "SubFolder",
		Type:     "directory",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder/Audio File in SubFolder")] = Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder/Audio File in SubFolder"),
		Name:     "Audio File in SubFolder",
		Type:     "audio",
		Parent:   utils.DoHash("root/Folder 1/SubFolder"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder/Video File in SubFolder")] = Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder/Video File in SubFolder"),
		Name:     "Video File in SubFolder",
		Type:     "video",
		Parent:   utils.DoHash("root/Folder 1/SubFolder"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/Image File in Folder 1")] = Node{
		Id:       utils.DoHash("root/Folder 1/Image File in Folder 1"),
		Name:     "Image File in Folder 1",
		Type:     "image",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/Word File in Folder 1")] = Node{
		Id:       utils.DoHash("root/Folder 1/Word File in Folder 1"),
		Name:     "Word File in Folder 1",
		Type:     "msword",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2")] = Node{
		Id:       utils.DoHash("root/Folder 2"),
		Name:     "Folder 2",
		Type:     "directory",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2/Text File in Folder 2")] = Node{
		Id:       utils.DoHash("root/Folder 2/Text File in Folder 2"),
		Name:     "Text File in Folder 2",
		Type:     "text",
		Parent:   utils.DoHash("root/Folder 2"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2/PowerPoint File in Folder 2")] = Node{
		Id:       utils.DoHash("root/Folder 2/PowerPoint File in Folder 2"),
		Name:     "PowerPoint File in Folder 2",
		Type:     "vnd.ms-powerpoint",
		Parent:   utils.DoHash("root/Folder 2"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/PDF File in Root")] = Node{
		Id:       utils.DoHash("root/PDF File in Root"),
		Name:     "PDF File in Root",
		Type:     "pdf",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/HTML File in Root")] = Node{
		Id:       utils.DoHash("root/HTML File in Root"),
		Name:     "HTML File in Root",
		Type:     "html",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	return Nodes
}
