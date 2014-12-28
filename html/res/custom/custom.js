//Initialize the JSTree
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

//Make JSTree automatically select the first node after it is ready
$(function(){
    $('#foldertree').on("ready.jstree", function(){
        $('#foldertree').jstree("select_node","ul > li:first")
        $('#foldertree').jstree("open_node","ul > li:first")
    });
});

var shownode = function(nodedata){
    $("#NodeType").text(nodedata.Type);
    $("#NodeSize").text(nodedata.Size);
    $("#NodeModTime").text(nodedata.ModTime);

    $("#NodeName").html('<button id="openbutton" data-id="'+nodedata.Id+'" class="btn btn-default" role="button"><span class="glyphicon glyphicon-open" /></button>&nbsp;'+nodedata.Name);
    
    $("#openbutton").click(function(){
        $.ajax({
            url:"/startnode",
            data:{
                "id":$("#openbutton").data("id")
            },
            error: function( xhr, status, errorThrown ) {
                alert( "Sorry, there was a problem!" + errorThrown + status );
            }
        });
    });

    $("#Comments").children().remove();
    nodedata.Comments.splice(0, 0, {Id:"CREATE",ModTime:"Click on the pencil to create a new comment...",Text:""});

    for (var i = 0; i < nodedata.Comments.length; i++){
        $("#Comments").append('<article class="panel"><header><div class="panel-heading">'+nodedata.Comments[i].ModTime+'<a id="cedit'+i.toString()+'" data-nodeid="'+nodedata.Id+'" data-commentid="'+nodedata.Comments[i].Id+'" data-text="'+ nodedata.Comments[i].Text+'" href="#"><span style="color:white;" class="glyphicon glyphicon-pencil pull-right"></span></a></div></header><div class="panel-body">'+marked(nodedata.Comments[i].Text)+'</div></article>');
                
        $("#cedit"+i.toString()).click(function(evt){
            $('#commenttextarea').val($(this).data("text"));
            $('#commenttextarea').data("nodeid",$(this).data("nodeid"));
            $('#commenttextarea').data("commentid",$(this).data("commentid"));
            $('#commenteditor').modal();
            evt.preventDefault();
        });
    } 
};


$(function(){
    $('#saveeditor').click(function(event) {
        $.ajax({
            url:"/editcomment",
            method:"POST",
            data:{
                "nodeid":$("#commenttextarea").data("nodeid"),
                "commentid":$("#commenttextarea").data("commentid"),
                "text":$("#commenttextarea").val()
            },
            success:function(){
                //alert("updated");
                $('#commenteditor').modal('hide');
                updatenode($("#commenttextarea").data("nodeid"));

            }
        });
    });
});

var updatenode = function(nodeid){

        $.ajax({
            url:"/servenode",
            datatype:"json",
            data:{
                "id":nodeid
            },
            success:function(data){
                shownode(data);
            },
            error: function( xhr, status, errorthrown ) {
                alert( "sorry, there was a problem!" );
                console.log( "error: " + errorthrown );
                console.log( "status: " + status );
                console.dir( xhr );
            }

        })

}


$(function(){
    $('#foldertree').on("select_node.jstree", function(e, data){
        updatenode(data.node.id);

    });

});

$(function(){
    $('#foldertree').on('dblclick', '.jstree-node', function (e) {
            $.ajax({
                url:"/startnode",
                data:{
                    "id":this.id
                },
                error: function( xhr, status, errorThrown ) {
                    alert( "Sorry, there was a problem!" + errorThrown + status );
                }
            });        
        e.stopPropagation();
     });
});

$('#refreshbutton').click(function() {
    location.reload(true);
});