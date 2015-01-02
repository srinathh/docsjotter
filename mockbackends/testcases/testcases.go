//package testcases generates the testcases used in fakebackend. Seed a different number to
//rand.Seed to get a different distribution of comments. The file structure is fixed
package main

import (
	"encoding/json"
	"fmt"
	"github.com/srinathh/powerdocs/structs"
	"github.com/srinathh/powerdocs/utils"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(9)
	buf, err := json.MarshalIndent(GenNodes(), "", "    ")
	if err != nil {
		log.Fatalf("Error in encoding comments to JSON : %s", err)
	}
	fmt.Println(string(buf))
}

const LOREM string = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

func GenComments(maxcomments int) []structs.Comment {

	c := make([]structs.Comment, rand.Intn(maxcomments))

	var hashtag string

	for i, _ := range c {

		switch rand.Intn(3) {
		case 0:
			hashtag = "#todo\n"
		case 1:
			hashtag = "#talk\n"
		default:
			hashtag = ""
		}

		c[i].Text = hashtag + LOREM[:rand.Intn(len(LOREM))]
		//adding current nanosecond to prevent comments with same text getting clashing ids
		c[i].Id = utils.DoHash(fmt.Sprintf("%s%d", c[i].Text, time.Now().UnixNano()))
		c[i].ModTime = time.Now().Format(time.Stamp)
	}
	return c
}

func GenNodes() map[string]structs.Node {

	Nodes := make(map[string]structs.Node)

	Nodes[utils.DoHash("root")] = structs.Node{
		Id:       utils.DoHash("root"),
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1"),
		Name:     "Folder 1",
		Type:     "directory",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder"),
		Name:     "SubFolder",
		Type:     "directory",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder/Audio File in SubFolder")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder/Audio File in SubFolder"),
		Name:     "Audio File in SubFolder",
		Type:     "audio",
		Parent:   utils.DoHash("root/Folder 1/SubFolder"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/SubFolder/Video File in SubFolder")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1/SubFolder/Video File in SubFolder"),
		Name:     "Video File in SubFolder",
		Type:     "video",
		Parent:   utils.DoHash("root/Folder 1/SubFolder"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/Image File in Folder 1")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1/Image File in Folder 1"),
		Name:     "Image File in Folder 1",
		Type:     "image",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 1/Word File in Folder 1")] = structs.Node{
		Id:       utils.DoHash("root/Folder 1/Word File in Folder 1"),
		Name:     "Word File in Folder 1",
		Type:     "msword",
		Parent:   utils.DoHash("root/Folder 1"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2")] = structs.Node{
		Id:       utils.DoHash("root/Folder 2"),
		Name:     "Folder 2",
		Type:     "directory",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2/Text File in Folder 2")] = structs.Node{
		Id:       utils.DoHash("root/Folder 2/Text File in Folder 2"),
		Name:     "Text File in Folder 2",
		Type:     "text",
		Parent:   utils.DoHash("root/Folder 2"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/Folder 2/PowerPoint File in Folder 2")] = structs.Node{
		Id:       utils.DoHash("root/Folder 2/PowerPoint File in Folder 2"),
		Name:     "PowerPoint File in Folder 2",
		Type:     "vnd.ms-powerpoint",
		Parent:   utils.DoHash("root/Folder 2"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/PDF File in Root")] = structs.Node{
		Id:       utils.DoHash("root/PDF File in Root"),
		Name:     "PDF File in Root",
		Type:     "pdf",
		Parent:   utils.DoHash("root"),
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: GenComments(5),
	}
	Nodes[utils.DoHash("root/HTML File in Root")] = structs.Node{
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
