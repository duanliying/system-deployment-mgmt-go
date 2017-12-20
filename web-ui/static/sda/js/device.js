// JS for SDA Manager device

function device_onLoad() {
    get_apps();
}

function get_apps() {
    sda_manager_control_hide();
    $("#app_tbody").empty();
    $("#service_tbody").empty();

    $.ajax({
        url: base_url +"/sdamanager/apps",
        type: "GET",
        contentType: "application/json",
        dataType: "json",
        error: function(error) {
            swal("Server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                var listLen = list.apps.length;
                $("#device_info").text(list.device);
                for (var i = 0; i < listLen; i++) {
                    var No = i + 1;
                    $("#app_table tbody").append('<tr title="' + list.apps[i].id + '">'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.apps[i].name + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.apps[i].services + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.apps[i].state + '</td>'
                    + '</tr>');
                }
            } else {
                swal("Server return error", "", "error");
            }
        }
    });
}

$(function() {
    $("#app_tbody").on("click", "tr", function(e) {
        sda_manager_control_show();
        $("#service_tbody").empty();
        $("tr").removeClass("active");
        $(this).addClass("active");
        $("#app_name").text($(this).find("td:eq(1)").text());

        var obj = new Object();
        obj.id = $(this).attr('title');
        $.ajax({
            url: base_url +"/sdamanager/app",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("Server return error",  "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var list = $.parseJSON(data);
                    var listLen = list.services.length;
                    for (var i = 0; i < listLen; i++) {
                        var No = i + 1;
                        $("#service_table tbody").append('<tr>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.services[i].name + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.services[i].state+ '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.services[i].exitcode + '</td>'
                        + '</tr>');
                    }
                } else {
                    swal("Server return error", "", "error");
                }
            }
        });
    });

    $('#btn_delete_app').click(function(){
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
                    url: base_url +"/sdamanager/app",
                    type: "DELETE",
                    error: function(error) {
                        swal("Server return error", "", "error");
                    },
                    success: function(data, code) {
                        if (code == "success") {
                            get_apps();
                            swal("Deleted!", "", "success");
                        } else {
                            swal("Server return error", "", "error");
                        }
                    }
                });
              }, 2000);
            }
        );
    });

    $("#btn_update_app").click(function() {
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
                        url: base_url +"/sdamanager/app/update",
                        type: "GET",
                        error: function(error) {
                            swal("Server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_apps();
                                swal("Updated!", "", "success");
                            } else {
                                swal("Server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_start_app").click(function() {
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
                        url: base_url +"/sdamanager/app/start",
                        type: "GET",
                        error: function(error) {
                            swal("Server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_apps();
                                swal("Started!", "", "success");
                            } else {
                                swal("Server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_stop_app").click(function() {
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
                        url: base_url +"/sdamanager/app/stop",
                        type: "GET",
                        error: function(error) {
                            swal("Server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                get_apps();
                                swal("Stoped!", "", "success");
                            } else {
                                swal("Server return error", "", "error");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_update_YAML").click(function() {
        $.ajax({
            url: base_url +"/sdamanager/app/yaml",
            type: "GET",
            error: function(error) {
                swal("Server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var obj = $.parseJSON(data);
                    obj = $.parseJSON(obj);
                    $("#textarea_update_yaml").val(obj);
                } else {
                    swal("Server return error", "", "error");
                }
            }
        });
    });

    $("#btn_confirm_update_YAML").click(function() {
        var obj = new Object();
        obj.data = $("#textarea_update_yaml").val();

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
                        url: base_url +"/sdamanager/app/yaml",
                        type: "POST",
                        contentType: "application/json",
                        data: JSON.stringify(obj),
                        error: function(error) {
                            swal("Server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                swal("Updated!", "", "success");
                                $('#dig_update_YAML').modal('toggle');
                            } else {
                                alert("Server return error.");
                            }
                        }
                    });
                }, 2000);
            }
        );
    });

    $("#btn_delete_device").click(function() {
        var isSuccess = false;
        swal({
            title:"Do you want to unregister this device?",
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
                        url: base_url +"/sdamanager/unregister",
                        type: "POST",
                        error: function(error) {
                            swal("Server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {

                                swal({
                                    title:"Unregister!",
                                    type:"success",
                                    allowOutsideClick:false,
                                    confirmButtonText:"Ok",
                                    closeOnConfirm: false,
                                    },
                                    function(){
                                        setTimeout(function(){
                                            $(location).attr("href", "/sdamanager");
                                        }, 0);
                                    }
                                );
                            }
                        }
                    });

                }, 2000);
            }
        );
    });
});
