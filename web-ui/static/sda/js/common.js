var base_url = window.location.protocol + "//" + window.location.host;
var yaml_index;
var yaml_name;

function sda_manager_control_show() {
        $("#sda_main").removeClass("sda_container")
        $("#sda_main").addClass("sda_container_with_control")
        $("#sda_main").removeClass("col-lg-12")
        $("#sda_main").addClass("col-lg-9")

        $("#sda_side").removeClass("div_hidden")
        $("#sda_side").addClass("div_show")
        $("#sda_side").removeClass("col-lg-0")
        $("#sda_side").addClass("col-lg-3")
}

function sda_manager_control_hide() {
        $("#sda_main").removeClass("sda_container_with_control");
        $("#sda_main").addClass("sda_container");
        $("#sda_main").removeClass("col-lg-9")
        $("#sda_main").addClass("col-lg-12")

        $("#sda_side").removeClass("div_show");
        $("#sda_side").addClass("div_hidden");
        $("#sda_side").removeClass("col-lg-3");
        $("#sda_side").addClass("col-lg-0");
}

function get_yamls() {
    $("#yaml_tbody").empty();
    $("#textarea_yaml").val("");

    $.ajax({
        url: base_url +"/sdamanager/yaml",
        type: "GET",
        dataType: "json",
        error: function(error) {
            swal("server return error", "", "error");
        },
        success: function(data, code) {
            if (code == "success") {
                var list = $.parseJSON(data);
                var listLen = list.yamls.length;
                for (var i = 0; i < listLen; i++) {
                    var No = i + 1;
                    $("#yaml_table tbody").append('<tr>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + No + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;"><img height="50" width="50" src="' + base_url + '/static/user/' + list.yamls[i].img + '"></td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.yamls[i].name + '</td>'
                    + '<td align="center" class="table_sda" style="vertical-align: middle;">' + list.yamls[i].description + '</td>'
                    + '</tr>');
                }
            } else {
                swal("server return error", "", "error");
            }
        }
    });
}

$(function() {
    $("#yaml_tbody").on("click", "tr", function(e) {
        $("tr").removeClass("active");
        $(this).addClass("active");

        $("#btn_deploy_app").removeAttr("disabled");
        $("#btn_deploy_group_app").removeAttr("disabled");

        yaml_index = parseInt($(this).find("td:eq(0)").text()) - 1;
        yaml_name = $(this).find("td:eq(2)").text();
        $.ajax({
            url: base_url +"/sdamanager/yaml",
            type: "GET",
            contentType: "application/json",
            dataType: "json",
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    var list = $.parseJSON(data);
                    $("#textarea_yaml").val(list.yamls[yaml_index].yaml);
                } else {
                    swal("server return error.");
                }
            }
        });
    });

    $("#btn_deploy_app").click(function() {
        $.ajaxSetup({
        beforeSend: function() {
            $("div[name=loading_bar]").show();
        },
        complete: function() {
            $("div[name=loading_bar]").hide();
        }
        });
        var obj = new Object();
        obj.data = $("#textarea_yaml").val();
        obj.name = yaml_name;
        $.ajax({
            url: base_url +"/sdamanager/app/install",
            type: "POST",
            contentType: "application/json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    swal("Installed!", "", "success");
                    $('#dig_install').modal('toggle');
                    get_apps();
                } else {
                    swal("server return error.");
                }
            }
        });
        $.ajaxSetup({
            beforeSend: function() {},
            complete: function() {}
        });
    });

    $("#btn_deploy_group_app").click(function() {
        $.ajaxSetup({
        beforeSend: function() {
            $("div[name=loading_bar]").show();
        },
        complete: function() {
            $("div[name=loading_bar]").hide();
        }
        });
        var obj = new Object();
        obj.data = $("#textarea_yaml").val();
        obj.name = yaml_name;
        $.ajax({
            url: base_url +"/sdamanager/group/deploy",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    swal("Installed!", "", "success");
                    $('#dig_group_install').modal('toggle');
                    get_groups();
                } else {
                    swal("server return error.");
                }
            }
        });
        $.ajaxSetup({
            beforeSend: function() {},
            complete: function() {}
        });
    });

    $("#btn_install_new_app").click(function() {
        get_yamls();
    });

    $("#btn_deploy_app_group").click(function() {
        get_yamls();
    });

    $("#yaml_tbody").on("dblclick", "tr", function(e) {
        $("tr").removeClass("active");
        $(this).addClass("active");
        swal("Edit/Remove will be supported!!");
    });

    $("#btn_confirm_add_new_yaml").click(function() {
        var obj = new Object();
        //obj.img = $("#create_new_yaml_icon").val();
        obj.img = "sample.png"
        obj.name = $("#create_new_yaml_name").val();
        obj.description = $("#create_new_yaml_description").val();
        obj.yaml = $("#create_new_yaml_yaml").val();

        $.ajax({
            url: base_url +"/sdamanager/yaml",
            type: "POST",
            contentType: "application/json",
            dataType: "text",
            data: JSON.stringify(obj),
            error: function(error) {
                swal("server return error", "", "error");
            },
            success: function(data, code) {
                if (code == "success") {
                    get_yamls();
                    $("#btn_deploy_app").attr("disabled", "disabled");
                    $("#btn_deploy_group_app").attr("disabled", "disabled");
                } else {
                    swal("server return error.");
                }
            }
        });
        $.ajaxSetup({
            beforeSend: function() {},
            complete: function() {}
        });
    });
});
