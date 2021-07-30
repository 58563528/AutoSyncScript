package models

var Admin = `<html lang="zh-cn">

    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>账号管理</title>
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/cdle/static/layui/css/layui.css">
        <script src="https://cdn.jsdelivr.net/gh/cdle/static/layui/layui.all.js"></script>
    </head>
    
    <body>
        <div class="layui-tab">
            <ul class="layui-tab-title">
                <li class="layui-this">账号管理</li>
                <li>系统设置</li>
            </ul>
            <div class="layui-tab-content">
                <div class="layui-tab-item layui-show">
                    <table id="accounts" lay-filter="accounts"></table>
                </div>
                <div class="layui-tab-item">
                   啥都没有
                </div>
            </div>
    </body>
    <script>
        var table = layui.table;
        table.render({
            elem: '#accounts',
            height: "auto",
            url: '/api/account',
            toolbar: 'default',
            response: {
                statusName: 'code',
                statusCode: 200,
                msgName: 'code',
                countName: 'message',
                dataName: 'data'
            },
            title: '账号列表',
            page: true,
            limit: 15,
            cols: [
                [ //表头
                    {
                        field: 'ID',
                        title: 'ID',
                        width: 100,
                        align: 'center',
                    }, {
                        field: 'Nickname',
                        title: '用户昵称',
                        width: 120,
                        align: 'center',
                    }, {
                        field: 'ScanedAt',
                        title: '扫码时间',
                        width: 110,
                        align: 'center',
                    }, {
                        field: 'BeanNum',
                        title: '京豆数目',
                        width: 90,
                        align: 'center',
                    }, {
                        field: 'Priority',
                        title: '优先级',
                        width: 80,
                        edit: 'text',
                        align: 'center',
                    }, {
                        field: 'Available',
                        title: '可用',
                        edit: 'text',
                        width: 80,
                        align: 'center',
                    }, {
                        field: 'Note',
                        title: '备注',
                        width: 120,
                        edit: 'text',
                        align: 'center',
                    }, {
                        field: 'PtPin',
                        title: 'PtPin',
                        width: 150,
                        align: 'center',
                    }, {
                        field: 'PtKey',
                        title: 'PtKey',
                        edit: 'text',
                        width: 500,
                        align: 'center',
                    }
                ]
            ]
        });
    
        table.on('edit(accounts)', function(obj) {
            obj.data.Priority = +obj.data.Priority
            obj.data.JinLi = +obj.data.JinLi
            layui.$.ajax({
                url: '/api/account',
                type: 'POST',
                contentType: "application/json",
                data: JSON.stringify(obj.data),
                dataType: 'json',
                timeout: 1000,
                cache: false,
                error: function() {
                    table.reload('accounts');
                }, //错误执行方法
                success: function(data) {
                    layer.msg(data["msg"])
                    table.reload('accounts');
                },
            });
        });
    </script>
    
    
    </html>`
