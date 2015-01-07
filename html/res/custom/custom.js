
//Initialize the JSTree after the page is ready
$(function(){
    $("#foldertree").jstree({
        "core" : {
            data :{
                url:"/servetree",
                dataType : "json"
            },
            multiple:false
        },
        "plugins":["wholerow","types","sort"], //"types"
        "types": {              
            "audio":{"icon":"fa fa-file-audio-o"},
            "vnd.ms-excel":{"icon":"fa fa-file-excel-o"},
            "default":{"icon":"fa fa-file-o"},
            "directory":{"icon":"fa fa-folder-open-o"},
            "image":{"icon":"fa fa-file-image-o"},
            "text":{"icon":"fa fa-file-text-o"},
            "pdf":{"icon":"fa fa-file-pdf-o"},
            "vnd.ms-powerpoint":{"icon":"fa fa-file-powerpoint-o"},
            "video":{"icon":"fa fa-file-video-o"},
            "vnd.ms-word":{"icon":"fa fa-file-word-o"},
            "msword":{"icon":"fa fa-file-word-o"},
            "zip":{"icon":"fa fa-file-archive-o"}
        }
    });
});


//Make JSTree automatically select and open the first node after JSTree is ready
$(function(){
    $('#foldertree').on("ready.jstree", function(){
        $('#foldertree').jstree("select_node","ul > li:first")
        $('#foldertree').jstree("open_node","ul > li:first")
    });
});

//on double click, start the node
$(function(){
    $('#foldertree').on('dblclick', '.jstree-node', function (e) {
        startnode(this.id); 
        e.stopPropagation();
     });
});

//onselecting a node, fetch node details and then show the node
$(function(){
    $('#foldertree').on("select_node.jstree", function(e, data){
        fetchnode(data.node.id, shownode);
    });
});

//issues a windows start or linux xdg-open command on the node
var startnode = function(nodeid){
    $.ajax({
        url:"/startnode",
        data:{
            "id":nodeid
        },
        error: function( xhr, status, errorThrown ) {
            show_alert(ALERT_WARNING,"Could not start the requested node")
        },
        success: function(){
            show_alert(ALERT_SUCCESS,"Started the requested node")  
        }
    });   
}

//fetchnode calls server on /servenode?id= to fetch the node with the required nodeid
//and on success, it calls the function onsuccess(data) with the returned node data
var fetchnode = function(nodeid, onsuccess){
    $.ajax({
        url:"/servenode",
        datatype:"json",
        data:{
            "id":nodeid
        },
        success:function(data){
            onsuccess(data);
        },
        error: function( xhr, status, errorthrown ) {
            show_alert(ALERT_DANGER, "Error fetching the requested node");
        }
    });
}


var shownode = function(nodedata){
    $("#contentheader")
        .html(buildheader(nodedata.Name, nodedata.Id,startnode))
        .append(builddetails({
            "File Type":nodedata.Type,
            "Last Modified":nodedata.ModTime,
            "Size":nodedata.Size
        }));

    $("#articles").html(buildarticle(nodedata.Id, "CREATE","Click on the pencil to create a new comment...",""));

    for( var i = 0; i < nodedata.Comments.length; i++){
        c = nodedata.Comments[i]
        $("#articles").append(buildarticle(nodedata.Id, c.Id,c.ModTime,c.Text,function(){
            fetchnode(nodedata.Id,shownode);
        }));
    }

};

//save the editor comment. Move to single id vs. multiple ids
var editorsave = function(nodeid,commentid, text){
   $.ajax({
        url:"/editcomment",
        method:"POST",
        data:{
            "commentid":commentid,
            "nodeid":nodeid,
            "text":text
        },
        error:function(){
            show_alert(ALERT_WARNING,"Could not save the comment");
        },
        success:function(){
            show_alert(ALERT_SUCCESS,"Comment saved")
        }
    }); 
}


//nodedetails is a hashmap which must have the keys Type, ModTime and Size
var buildheader = function(title,id){
    header = ich.pagetitle({Title:title});
    but = $(header).find("button")

    if(id !== undefined){
        $(but).click(function(){
            startnode(id);
        });
    }else{
        $(but).hide();
    }

    return header;
}

var builddetails = function(hashmap){
    headings = []
    values = []
    for (var heading in hashmap){
        headings.push(heading)
        values.push(hashmap[heading])
    }
    return ich.nodedetails({"headings":headings, "values":values});
}

var buildarticle = function(nodeid, commentid, title, text, editrefresh, titleclick){
    article = ich.articletemplate({"ArticleTitle":title, "ArticleBody":marked(text)});

    $(article).find(".editbutton").click(function(evt){
        showeditor(nodeid, commentid, text, editrefresh);
        evt.stopPropagation();
    });

    if(titleclick !== undefined){
        $(article).find(".articletitle").click(function(evt){
            titleclick();
            evt.stopPropagation();
        });
    }
    return article;
}


var showeditor = function(nodeid, commentid, text, editrefresh){
    $('#editmodalcontainer').html(ich.modaltemplate({Text:text}));
    editor = new Editor();
    editor.render();

    $('#commenteditor').on('shown.bs.modal', function (e) {
        editor.codemirror.refresh();
    });

    $('#commenteditor').find("button").click(function(){
        editorsave(nodeid, commentid, editor.codemirror.getValue());
        $('#commenteditor').modal('hide');
    });


    $('#commenteditor').on('hidden.bs.modal', function (e) {
        $('#editmodalcontainer').children().remove();
        editrefresh();
    });

    $('#commenteditor').modal();
}


var ALERT_SUCCESS = "alert-success";
var ALERT_INFO = "alert-info";
var ALERT_WARNING = "alert-warning"
var ALERT_DANGER = "alert-danger"

var show_alert = function(alert_type, alert_text){
    //$('#alertcontainer').html('<div class="alert '+alert_type+' alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'+alert_text+'</div>');
    $('#alertcontainer').html(ich.alertcontent({"alert_text":alert_text, "alert_type":alert_type}));
}

//tags
$(function(){
    $('#viewtagbtn').click(function(evt){
        show_alert(ALERT_INFO, "Will try to fetch and show the tag details");
        evt.stopPropagation();
    });
});

