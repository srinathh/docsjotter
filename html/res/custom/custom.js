$(function(){
    $("#foldertree").jstree({
        "core" : {
            data :{
                url:"/servetree",
                dataType : "json"
            },
            multiple:false
        },
        "plugins":["wholerow","types"], //"types"
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
        //"types": {"Archive":{"icon":"fa fa-file-archive-o"},"Audio File":{"icon":"fa fa-file-audio-o"},"Excel Spreadsheet":{"icon":"fa fa-file-excel-o"},"File":{"icon":"fa fa-file-o"},"Folder":{"icon":"fa fa-folder-open-o"},"HTML File":{"icon":"fa fa-html5"},"Image File":{"icon":"fa fa-file-image-o"},"Markdown Document":{"icon":"fa fa-file-text-o"},"PDF File":{"icon":"fa fa-file-pdf-o"},"Powerpoint Presentation":{"icon":"fa fa-file-powerpoint-o"},"Programming Code":{"icon":"fa fa-file-code-o"},"Text File":{"icon":"fa fa-file-text-o"},"Video File":{"icon":"fa fa-file-video-o"},"Word Document":{"icon":"fa fa-file-word-o"}}
    });
});

$(function(){
    $('#foldertree').on("ready.jstree", function(){
        $('#foldertree').jstree("select_node","0")
    });
});

$(function(){
    $('#foldertree').on("select_node.jstree", function(e, data){
        $.ajax({
            url:"/servenode",
            dataType:"json",
            data:{
                "id":data.node.id
            },
            success:function(nodedata){
                $("#NodeName").html('<button data.id="'+nodedata.Id+'" class="btn btn-default" role="button"><span class="glyphicon glyphicon-open" /></button>&nbsp;'+nodedata.Name);
                $("#NodeType").text(nodedata.Type);
                $("#NodeSize").text(nodedata.Size);
                $("#NodeModTime").text(nodedata.ModTime);
                $("#Comments").children().remove();
                //tk add handler for the button

                for (var i = 0; i < nodedata.Comments.length; i++){
                    $("#Comments").append('<article class="panel"><header><div class="panel-heading">Updated '+nodedata.Comments[i].ModTime+'<a data.id="'+nodedata.Comments[i].Id+'" href="#"><span style="color:white;" class="glyphicon glyphicon-pencil pull-right"></span></a></div></header><div class="panel-body">'+nodedata.Comments[i].Text+'</div></article>');
                    //tk add handler for the edit icon
                } 


            },
            error: function( xhr, status, errorThrown ) {
                alert( "Sorry, there was a problem!" );
                console.log( "Error: " + errorThrown );
                console.log( "Status: " + status );
                console.dir( xhr );
            },
        })
    })
})

/*
$(function() {
    $("#foldertree").jstree({
        "core" : {
            "data" :{
                "url":"/dirtree"
            }
        },
        "plugins":["sort", "types", "wholerow"],
        "types": {"Archive":{"icon":"fa fa-file-archive-o"},"Audio File":{"icon":"fa fa-file-audio-o"},"Excel Spreadsheet":{"icon":"fa fa-file-excel-o"},"File":{"icon":"fa fa-file-o"},"Folder":{"icon":"fa fa-folder-open-o"},"HTML File":{"icon":"fa fa-html5"},"Image File":{"icon":"fa fa-file-image-o"},"Markdown Document":{"icon":"fa fa-file-text-o"},"PDF File":{"icon":"fa fa-file-pdf-o"},"Powerpoint Presentation":{"icon":"fa fa-file-powerpoint-o"},"Programming Code":{"icon":"fa fa-file-code-o"},"Text File":{"icon":"fa fa-file-text-o"},"Video File":{"icon":"fa fa-file-video-o"},"Word Document":{"icon":"fa fa-file-word-o"}}
    });

    $('#foldertree').on("select_node.jstree", function (evt, data) {
        $.ajax({
            url: "/node",
            data: { 
                "id": data.node.id, 
            },
            cache: false,
            type: "POST",
            success: function(response) {
                $("#contents").html(response)
            },
            error: function(xhr) {
                alert(xhr)
            }
        });
    });

    //$('#refresh').click( function() { window.location = 'http://www.google.com'; });


});


var OpenNode = function(id){
     $.ajax({
        url: "/startnode",
        data: { 
           "id": id, 
        },
        cache: false,
        type: "POST",
        error: function(xhr) {
            alert(xhr)
        }
    });

};

var EditNode = function(id){

    if(id=="MDCREATENEW"){
        alert("not implemented")
    }else{
         $.ajax({
            url: "/editmd",
            data: { 
               "id": id, 
            },
            cache: false,
            type: "POST",
            error: function(xhr) {
                alert(xhr)
            },
            success: function(ret){
                $('#EditModal').html(ret)
                $('#EditModal').modal() 
            }
        });
    }

}
*/