// JS for SDA Manager app

var selected_app_name;

function app_onLoad() {
    get_group_apps();
}

function get_group_apps() {
    sda_manager_control_hide();
    $("#group_app_tbody").empty();
    $("#group_device_tbody").empty();

    $.ajax({
        url: base_url +"/sdamanager/group/apps",
        type: "GET",
        contentType: "application/json",
        dataType: "json",
        error: function(error) {
            swal("server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                var listLen = list.apps.length;
                $("#group_info").text(list.group);
                for (var i = 0; i < listLen; i++) {
                    var No = i + 1;
                    $("#group_app_table tbody").append('<tr title="' + list.apps[i].id + '">'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.apps[i].name + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.apps[i].services + '</td>'
                    + '</tr>');
                }
            } else {
                swal("server return error", "", "error");
            }
        }
    });
}

$(function() {
    $("#group_app_tbody").on("click", "tr", function() {
        sda_manager_control_show();
        $("#group_device_tbody").empty();
        $("#group_app_name").text($(this).find("td:eq(1)").text());
        // focus clicked tr
        $("tr").removeClass("active");
        $(this).addClass("active");

        var obj = new Object();
        obj.id = $(this).attr('title');
        selected_app_name = $(this).find("td:eq(1)").text();

        $.ajax({
            url: base_url + "/sdamanager/group/app",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var list = $.parseJSON(data);
                    var listLen = list.devices.length;
                    for (var i = 0; i < listLen; i++) {
                        var No = i + 1;
                        $("#group_device_table tbody").append('<tr title="' + list.devices[i].id + '">'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.devices[i].host + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.devices[i].port + '</td>'
                        + '</tr>');
                    }
                }
                else {
                    swal("server return error", "", "error");
                }
            }
        });
    });

    $("#group_device_tbody").on("click", "tr", function() {
        var obj = new Object();
        obj.id = $(this).attr('title');
        obj.host = $(this).find("td:eq(1)").text();
        obj.port = $(this).find("td:eq(2)").text();

        $.ajax({
            url: base_url + "/sdamanager/device",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    $(location).attr("href", data);
                }
                else {
                    swal("server return error", "", "error");
                }
            }
        });
    });

     $('#btn_delete_group_app').click(function(){
		swal({
            title:"Are you sure?",
            text:"The app will not be able to recover this.",
            type:"warning",
            allowOutsideClick:false,
            showCancelButton: true,
            confirmButtonText:"Yes",
            confirmButtonColor: "#DD6B55",
            cancelButtonText:"No",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
            },
            function(){
                setTimeout(function(){
                    $.ajax({
                    url: base_url +"/sdamanager/group/app",
                    type: "DELETE",
                    error: function(error) {
                        swal("server return error", "", "error");
                    },
                    success: function(data, code) {
                        if (code == "success") {
                            get_group_apps();
                            swal("Deleted!", "", "success");
                        } else {
                            swal("server return error", "", "error");
                        }
                    }
                });
              }, 2000);
            }
        );
	});

    $("#btn_update_group_app").click(function() {
        swal({
            title:"Do you want to update?",
            text:"It only supports latest tag",
            type:"info",
            allowOutsideClick:false,
            showCancelButton: true,
            confirmButtonText:"Yes",
            cancelButtonText:"No",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
            },
            function(){
                setTimeout(function(){
                    $.ajax({
                        url: base_url +"/sdamanager/group/app/update",
                        type: "GET",
                        error: function(error) {
                            swal("server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_group_apps();
                                swal("Updated!", "", "success");
                            } else {
                                swal("server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_start_group_app").click(function() {
        swal({
            title:"Do you want to start?",
            type:"info",
            allowOutsideClick:false,
            showCancelButton: true,
            confirmButtonText:"Yes",
            cancelButtonText:"No",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
            },
            function(){
                setTimeout(function(){
                    $.ajax({
                        url: base_url +"/sdamanager/group/app/start",
                        type: "GET",
                        error: function(error) {
                            swal("server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_group_apps();
                                swal("Started!", "", "success");
                            } else {
                                swal("server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_stop_group_app").click(function() {
        swal({
            title:"Do you want to stop?",
            type:"info",
            allowOutsideClick:false,
            showCancelButton: true,
            confirmButtonText:"Yes",
            cancelButtonText:"No",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
            },
            function(){
                setTimeout(function(){
                    $.ajax({
                        url: base_url +"/sdamanager/group/app/stop",
                        type: "GET",
                        error: function(error) {
                            swal("server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_group_apps();
                                swal("Stoped!", "", "success");
                            } else {
                                swal("server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_confirm_update_group_YAML").click(function() {
        var obj = new Object();
        obj.data = $("#textarea_update_group_yaml").val();

        swal({
            title:"Do you want to update YAML?",
            type:"info",
            allowOutsideClick:false,
            showCancelButton: true,
            confirmButtonText:"Yes",
            cancelButtonText:"No",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
            },
            function(){
                setTimeout(function(){
                    $.ajax({
                        url: base_url +"/sdamanager/group/app/yaml",
                        type: "POST",
                        contentType: "application/json",
                        data: JSON.stringify(obj),
                        error: function(error) {
                            swal("server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                swal("Updated!", "", "success");
                                $('#dig_update_group_YAML').modal('toggle');
                            } else {
                                alert("server return error.");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });
});