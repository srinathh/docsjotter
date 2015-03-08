package comments

import (
	"bytes"
	"fmt"
	//"github.com/srinathh/docsjotter/utils"
	"github.com/srinathh/hashtag"
	"github.com/srinathh/toml"
	"io"
	"io/ioutil"
	//"regexp"
	"time"
)

//Comment represents the data in a single file comment. BaseName stores the filepath.Base() name
//of the file that the comment refers to (Comments are stored in the same folder as the file). If
//BaseName is "." then the comment refers to the folder rather than to a file. Comment Text should
//be MarkDown formatted text which can contain hashtags
type Comment struct {
	Id       string
	BaseName string    //for parent folder: "."; for regular files filepath.Base()
	ModTime  time.Time //last modified time of the comment
	Text     string    `modifier:"multiline_string"` //Markdown formatted text
}

func (c *Comment) ExtractHashTags() []string {
	return hashtag.ExtractHashtags(c.Text)
}

//CommentFile is a TOML encoded file where comments are stored by Comment ID
type CommentsFile map[string]Comment

func (c *CommentsFile) ReadFrom(r io.Reader) (n int64, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return int64(len(data)), fmt.Errorf("CommentFile: Error in reading data:%s", err)
	}

	datar := bytes.NewReader(data)
	if _, err := toml.DecodeReader(datar, c); err != nil {
		return int64(len(data)), fmt.Errorf("Commentfile: Error parsing TOML:%s", err)
	}

	return int64(len(data)), nil
}

func (c *CommentsFile) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(c); err != nil {
		return 0, fmt.Errorf("CommentFile: Error in toml encoding: %s", err)
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}
