<template>
<dl class="row">
    <dt class="col-sm-3 text-right">接口名</dt>
    <dd class="col-sm-9">{{interface.name}}</dd>
    <dt class="col-sm-3 text-right">模块名</dt>
    <dd class="col-sm-9">{{interface.model}}</dd>
    <dt class="col-sm-3 text-right">表名</dt>
    <dd class="col-sm-9">{{interface.table}}</dd>
    <el-divider>HTTP请求方式</el-divider>
    <dt class="col-sm-3 text-right">方法</dt>
    <dd class="col-sm-9">{{interface.method}}</dd>
    <dt class="col-sm-3 text-right">路由</dt>
    <dd class="col-sm-9">{{interface.route}}</dd> 
    <el-divider>HTTP请求/响应</el-divider>
    <dt class="col-sm-3 text-right">参数</dt>
    <dd class="col-sm-9">{{interface.params}}</dd>
    <dt class="col-sm-3 text-right">返回值</dt>
    <dd class="col-sm-9">{{interface.return}}</dd>
</dl>
</template>

<script>
export default {
    props: {
        "interface": Object
    }
}
</script>
