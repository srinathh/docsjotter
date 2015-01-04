
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
    $("#contentheader").html(ich.filenode(nodedata));

    $("#openbutton").click(function(){
        startnode($(this).data("id"));
    });

    nodedata.Comments.splice(0, 0, {Id:"CREATE",ModTime:"Click on the pencil to create a new comment...",Text:""});

    $("#Comments").children().remove();

    for (var i = 0; i < nodedata.Comments.length; i++){
        nodedata.Comments[i]["MDParsedText"]=marked(nodedata.Comments[i]["Text"]);
        nodedata.Comments[i]["IconName"]="glyphicon-pencil";
        $("#Comments").append(ich.nodecomment(nodedata.Comments[i]));
    }

    $(".cedit").click(function(evt){
        $('#commenttextarea').data("nodeid",$("#NodeName").data("id"));
        $('#commenttextarea').data("commentid",$(this).data("id"));
        editor.codemirror.setValue($(this).data("text"));
        $('#commenteditor').modal();
        evt.preventDefault();
    });
};

//-------------------------------------------------------------------------
//Editor javascript

var intervalid;

//initialize the editor on page ready
var editor;
$(function(){
    editor = new Editor();
    editor.render();
});

//refresh codemirror after edit window is show and start autosave
$('#commenteditor').on('shown.bs.modal', function (e) {
    editor.codemirror.refresh(); 
    intervalid = setInterval(function(){
        dosave(false);
    }, 2500);
});

//shut down autosave when editor window is hidden and update the node displayed
$('#commenteditor').on('hidden.bs.modal', function (e) {
    clearInterval(intervalid);
    fetchnode($("#NodeName").data("id"),shownode);
});

//save the editor comment
var dosave = function(dohide){
   $.ajax({
        url:"/editcomment",
        method:"POST",
        data:{
            "nodeid":$("#commenttextarea").data("nodeid"),
            "commentid":$("#commenttextarea").data("commentid"),
            //"text":$("#commenttextarea").val()
            "text":editor.codemirror.getValue()
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
    $('#alertcontainer').html('<div class="alert '+alert_type+' alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'+alert_text+'</div>');
}



//tags
$(function(){
    $('#viewtagbtn').click(function(e){
        show_alert(ALERT_INFO, "Will try to fetch and show the tag details");
        e.preventDefault();
    });
});