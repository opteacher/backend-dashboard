<template>
<dl class="row">
    <dt class="col-sm-3 text-right">接口名</dt>
    <dd class="col-sm-9">{{api.name}}</dd>
    <dt class="col-sm-3 text-right" v-if="api.model">模块名</dt>
    <dd class="col-sm-9">{{api.model}}</dd>
    <dt class="col-sm-3 text-right" v-if="api.table">表名</dt>
    <dd class="col-sm-9">{{api.table}}</dd>
    <el-divider>HTTP请求方式</el-divider>
    <dt class="col-sm-3 text-right">方法</dt>
    <dd class="col-sm-9">{{api.method}}</dd>
    <dt class="col-sm-3 text-right">路由</dt>
    <dd class="col-sm-9">{{api.route}}</dd> 
    <el-divider>HTTP请求/响应</el-divider>
    <dt class="col-sm-3 text-right">参数</dt>
    <dd class="col-sm-9">{{api.params}}</dd>
    <dt class="col-sm-3 text-right">返回值</dt>
    <dd class="col-sm-9">{{api.returns}}</dd>
</dl>
</template>

<script>
export default {
    props: {
        "api": Object
    }
}
</script>
