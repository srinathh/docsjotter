DirTree / Nodes
---------------
Root Node should have the id of 0. It is auto selectd on directory load

Node Types
----------
They `type` parameter within each JSTree Node entry is used to assign the right icon to each node. Thea available node types are as follows.
```
    "audio":{"icon":"fa fa-file-audio-o"},
    "image":{"icon":"fa fa-file-image-o"},
    "text":{"icon":"fa fa-file-text-o"},
    "video":{"icon":"fa fa-file-video-o"},

    "vnd.ms-excel":{"icon":"fa fa-file-excel-o"},
    "default":{"icon":"fa fa-file-o"},
    "directory":{"icon":"fa fa-folder-open-o"},
    "html":{"icon":"fa fa-html5"},
    "pdf":{"icon":"fa fa-file-pdf-o"},
    "vnd.ms-powerpoint":{"icon":"fa fa-file-powerpoint-o"},
    "vnd.ms-word":{"icon":"fa fa-file-word-o"},
    "msword":{"icon":"fa fa-file-word-o"},
    "zip":{"icon":"fa fa-file-archive-o"}
```

If a MIME magic function is being 
        //MIME parsing rules 
        // if the node is a directory use "directory" as type else run magic identification
        // if the first word is audio, image or video just use it directly
        // if the first word is application, parse the second word before the dot
        // if it is any of the types above, use them, else use default
        // parse text for html or any others         
