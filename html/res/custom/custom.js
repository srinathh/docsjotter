
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

//nodedetails is a hashmap which must have the keys Type, ModTime and Size
var buildheader = function(title,id, hide_button){
    header = ich.pagetitle({Title:title, Id:id});
    
    if (hide_button !== undefined){
        $(header).find('#'+id+'_button').hide();
    }
    return header;
}

var builddetails = function(hashmap){
    showkeys = []
    showvalues = []
    for (var key in hashmap){
        showkeys.push(key)
        showvalues.push(hashmap[key])
    }
    return ich.nodedetails({keys:showkeys, values:showvalues});
}

var buildarticle = function(id, title, text){
    return ich.articletemplate({"ArticleId":id, "ArticleTitle":title, "ArticleText":text, "ArticleBody":marked(text)});
}

var fetchid=function(text){
    return text.substr(0,text.lastIndexOf("_"));
}

var shownode = function(nodedata){
    header = buildheader(nodedata.Name, nodedata.Id);
    details = builddetails({"File Type":nodedata.Type,"Last Modified":nodedata.ModTime,"Size":nodedata.Size});
    $("#contentheader").html(header);
    $("#contentheader").append(details);

    $("#articles").children().remove();
    $("#articles").append(buildarticle("CREATE","Click on the pencil to create a new comment...",""));

    for( var i = 0; i < nodedata.Comments.length; i++){
        c = nodedata.Comments[i]
        $("#articles").append(buildarticle(c.Id,c.ModTime,c.Text));
    }

    $(".articletitle").click(function(e){ e.preventDefault(); });
   
    $(".articlebutton").click(function(e){
        var ondone = function(){
            fetchnode(nodedata.Id, shownode);
        }
        showeditor($(this).data('text'),fetchid($(this).attr("id")),nodedata.Id, ondone);
        e.preventDefault();
    });
};

var showeditor = function(text, id, nodeid, ondone){
    $('#editorcontainer').append('<textarea></textarea>')
    editor = new Editor();
    editor.render();
    editor.codemirror.setValue(text);
    $('#commenteditor').data("editor", editor);
    $('#commenteditor').data("id", id);
    $('#commenteditor').data("nodeid", nodeid);
    $('#commenteditor').data("ondone", ondone);
    $('#commenteditor').modal();
}

//tk modularize the editor

//-------------------------------------------------------------------------
//Editor javascript

//refresh codemirror after edit window is show and start autosave
$('#commenteditor').on('shown.bs.modal', function (e) {
    editor.codemirror.refresh();
});

//shut down autosave when editor window is hidden and update the node displayed
$('#commenteditor').on('hidden.bs.modal', function (e) {
    $('#editorcontainer').children().remove();
    $('#commenteditor').data("ondone")();
});

//save the editor comment. Move to single id vs. multiple ids
var dosave = function(dohide){
   $.ajax({
        url:"/editcomment",
        method:"POST",
        data:{
            "commentid":$("#commenteditor").data("id"),
            "nodeid":$("#commenteditor").data("nodeid"),
            "text":$("#commenteditor").data("editor").codemirror.getValue()
        },
        success:function(){
            //alert("updated");
            if(dohide == true){
                $('#commenteditor').modal('hide');
            }
        }
    }); 
}


$(function(){
    $('#saveeditor').click(function(event) {
        dosave(true);
    });
});



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
    $('#viewtagbtn').click(function(e){
        show_alert(ALERT_INFO, "Will try to fetch and show the tag details");
        e.preventDefault();
    });
});