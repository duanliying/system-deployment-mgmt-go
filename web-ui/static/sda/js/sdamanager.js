// JS for SDA Manager

var not_include_devices;
var include_devices;
var not_include_devices_agents;
var include_devices_agents;

var SDAManagerAddress;

function sdamanager_onLoad() {
    get_groups();
}

function isSDAManagerAddress(){
    SDAManagerAddress = "";
    $.ajax({
        url: base_url + "/sdamanager/check/address",
        type: "GET",
        error: function(error) {
            swal("Please input SDA Manager Address!", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                SDAManagerAddress = data;
            }
            else {
                swal("server return error", "", "error");
            }
        }
    });
}

function get_groups() {
    sda_manager_control_hide();
    $("#group_tbody").empty();

    $.ajax({
        url: base_url + "/sdamanager/groups",
        type: "GET",
        contentType: "application/json",
        dataType: "json",
        error: function(error) {
            swal("server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                var listLen = list.groups.length;
                $("#sda_manager_ip").val(list.address);
                for (var i = 0; i < listLen; i++) {
                    var No = i + 1;
                    $("#group_table tbody").append('<tr title="' + list.groups[i].id + '">'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.groups[i].groupname + '</td>'
                        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.groups[i].members + '</td>'
                        + '</tr>');
                }
                get_devices();
            }
            else {
                    swal("server return error", "", "error");
            }
        }
    });
}

function get_devices() {
    $("#device_tbody").empty();

    $.ajax({
        url: base_url + "/sdamanager/devices",
        type: "GET",
        contentType: "application/json",
        dataType: "json",
        error: function(error) {
            swal("server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                var listLen = list.devices.length;

                for (var i = 0; i < listLen; i++) {
                    var No = i + 1;
                    $("#device_table tbody").append('<tr title="' + list.devices[i].id + '">'
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
}

function get_group_devices(not_include_devices_agents) {
    $.ajax({
        url: base_url + "/sdamanager/group/devices",
        type: "GET",
        contentType: "application/json",
        dataType: "json",
        error: function(error) {
            swal("server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                include_devices_agents = list.devices;
                for (var i = 0; i < include_devices_agents.length; i++) {
                    for (var j = 0; j < not_include_devices_agents.length; j++) {
                        if (not_include_devices_agents[j].id == include_devices_agents[i].id){
                            not_include_devices_agents.splice(j, 1);
                            break;
                        }
                    }
                }
                include_devices.agents = include_devices_agents;
                not_include_devices.agents = not_include_devices_agents;
                show_devices();
            }
            else {
                swal("server return error", "", "error");
            }
        }
    });
}

function show_devices(){
    $("tbody[name=not_include_device_tbody]").empty();
    $("tbody[name=include_device_tbody]").empty();
    var listLen = not_include_devices.agents.length;

    for (var i = 0; i < listLen; i++) {
        var No = i + 1;
        $("table[name=not_include_device_table] tbody").append('<tr title="' + not_include_devices.agents[i].id + '">'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + not_include_devices.agents[i].host + '</td>'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + not_include_devices.agents[i].port + '</td>'
        + '</tr>');
    }

    listLen = include_devices.agents.length;
    for (var i = 0; i < listLen; i++) {
        var No = i + 1;
        $("table[name=include_device_table] tbody").append('<tr title="' + include_devices.agents[i].id + '">'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + include_devices.agents[i].host + '</td>'
        + '<td align="center" class="table_sda" style="vertical-align: middle;">' + include_devices.agents[i].port + '</td>'
        + '</tr>');
    }
}

function get_agent_list(data){
    var obj = new Object();
    var arr = new Array();

    for (var i = 0; i < data.agents.length; i++) {
        arr.push(data.agents[i].id);
    }

    obj.agents = arr;
    return obj;
}

$(function() {
    $("#group_tbody").on("click", "tr", function() {
        sda_manager_control_show();
        $("#device_tbody").empty();
        $("#group_name").text($(this).find("td:eq(1)").text());
        // focus clicked tr
        $("tr").removeClass("active");
        $(this).addClass("active");

        var obj = new Object();
        obj.id = $(this).attr('title');
        obj.groupname = $(this).find("td:eq(1)").text();
        selected_group_name = $(this).find("td:eq(1)").text();

        $.ajax({
            url: base_url + "/sdamanager/group",
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
                        $("#device_table tbody").append('<tr title="' + list.devices[i].id + '">'
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

    $("#device_tbody").on("click", "tr", function() {
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

    $("#btn_set_sdam_address").click(function() {
        $("#sda_manager_section").removeClass("section_with_control")
        var obj = new Object();
        obj.ip = $("#sda_manager_ip").val();

        $.ajax({
            url: base_url + "/sdamanager/address",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    $("#sda_manager_ip").val(data);
                    get_groups();
                }
                else {
                    swal("server return error", "", "error");
                }
            }
        });
    });

    $("#btn_show_all_device").click(function() {
        get_groups();
    });

    $("#btn_create_new_group").click(function() {
        if(SDAManagerAddress == "") return;

        not_include_devices = new Object();
        include_devices = new Object();
        not_include_devices_agents = new Array();
        include_devices_agents = new Array();

        include_devices.agents = include_devices_agents;
        not_include_devices.agents = not_include_devices_agents;

        $.ajax({
            url: base_url + "/sdamanager/devices",
            type: "GET",
            contentType: "application/json",
            dataType: "json",
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var list = $.parseJSON(data);
                    var listLen = list.devices.length;
                    for (var i = 0; i < listLen; i++) {
                        var No = i + 1;
                        not_include_devices_agents.push(list.devices[i]);
                    }
                    not_include_devices.agents = not_include_devices_agents;
                    show_devices();
                }
                else {
                    swal("server return error", "", "error");
                }
            }
        });
    });

    $("tbody[name=not_include_device_tbody]").on("dblclick", "tr", function() {
        var obj = new Object();
        obj.id = $(this).attr('title');

        for (var i = 0; i < not_include_devices.agents.length; i++) {
            if(not_include_devices.agents[i].id == obj.id){
                include_devices_agents.push(not_include_devices.agents[i]);
                not_include_devices.agents.splice(i, 1);
                break;
            }
        }
        include_devices.agents = include_devices_agents;
        show_devices();
    });

    $("tbody[name=include_device_tbody]").on("dblclick", "tr", function() {
        var obj = new Object();
        obj.id = $(this).attr('title');

        for (var i = 0; i < include_devices.agents.length; i++) {
            if(include_devices.agents[i].id == obj.id){
                not_include_devices.agents.push(include_devices.agents[i]);
                include_devices.agents.splice(i, 1);
                break;
            }
        }
        not_include_devices.agents = not_include_devices_agents;
        show_devices();
    });

    $("#btn_confirm_create_new_group").click(function() {
        if(SDAManagerAddress == "") return ;

        var condition = 1;
        var obj = new Object();
        obj.members = get_agent_list(include_devices);
        obj.groupname = $("#create_new_group_name").val();

        if(obj.groupname == ""){
            swal("Please input a group name", "", "error");
            condition = 0;
        }

        if(obj.members.agents.length == 0){
            swal("Please add at least one device", "", "error");
            condition = 0;
        }

        if(condition){
            $.ajax({
                url: base_url + "/sdamanager/group/create",
                type: "POST",
                contentType: "application/json",
                data: JSON.stringify(obj),
                error: function(error) {
                    swal("server return error", "", "error");
                },
                success: function(data, code) {
                    if (code == "success") {
                        swal("Created!", "", "success");
                        $('#dig_create_group').modal('toggle');
                        get_groups();
                    }
                    else {
                        swal("server return error", "", "error");
                    }
                }
            });
        }
    });

    $("#btn_delete_group").click(function() {
        swal({
            title:"Are you sure?",
            text:"The group will not be able to recover this.",
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
                    url: base_url +"/sdamanager/group/delete",
                    type: "DELETE",
                    error: function(error) {
                        swal("server return error", "", "error");
                    },
                    success: function(data, code) {
                        if (code == "success") {
                            get_groups();
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

    $("#btn_members_group").click(function() {
        not_include_devices = new Object();
        include_devices = new Object();
        not_include_devices_agents = new Array();
        include_devices_agents = new Array();

        include_devices.agents = include_devices_agents;
        not_include_devices.agents = not_include_devices_agents;

        $("#members_group_name").text(" - " + selected_group_name);

        $.ajax({
            url: base_url + "/sdamanager/devices",
            type: "GET",
            contentType: "application/json",
            dataType: "json",
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var list = $.parseJSON(data);
                    not_include_devices_agents = list.devices;
                    get_group_devices(not_include_devices_agents);
                }
                else {
                    swal("server return error", "", "error");
                }
            }
        });
    });

    $("#btn_save_group_members").click(function() {
        var condition = 1;
        var obj = new Object();
        obj = get_agent_list(include_devices);

        if(obj.agents.length == 0){
            swal("Please add at least one device", "", "error");
            condition = 0;
        }

        if(condition){
            $.ajax({
                url: base_url + "/sdamanager/group/members",
                type: "POST",
                contentType: "application/json",
                data: JSON.stringify(obj),
                error: function(error) {
                    swal("server return error", "", "error");
                },
                success: function(data, code) {
                    if (code == "success") {
                        swal("Saved!", "", "success");
                        $('#dig_members_group').modal('toggle');
                        get_groups();
                    }
                    else {
                        swal("server return error", "", "error");
                    }
                }
            });
        }
    });

    $("#btn_group_app").click(function() {
        $(location).attr("href", "/app");
    });

    $("#btn_add_device").click(function() {
        if(SDAManagerAddress == "") return ;
    });

    $("#btn_confirm_register_new_device").click(function() {
        if(SDAManagerAddress == "") return ;

        var obj = new Object();
        obj.ip = $("#register_new_device_ip").val();
        obj.interval = $("#register_new_device_interval").val();

        swal({
            title:"Do you want to register this device?",
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
                        url: base_url +"/sdamanager/register",
                        type: "POST",
                        contentType: "application/json",
                        data: JSON.stringify(obj),
                        error: function(error) {
                            swal("server return error", "", "error");
                        },
                        success: function(data, code) {
                            if (code == "success") {
                                swal("Register!", "", "success");
                                $('#dig_register_device').modal('toggle');
                                get_devices();
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