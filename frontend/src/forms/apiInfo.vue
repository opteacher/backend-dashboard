<template>
<dl class="row">
    <dt class="col-sm-3 text-right">接口名</dt>
    <dd class="col-sm-9">{{api.name}}</dd>
    <dt class="col-sm-3 text-right" v-if="api.model">模块名</dt>
    <dd class="col-sm-9">{{api.model}}</dd>
    <dt class="col-sm-3 text-right" v-if="api.table">表名</dt>
    <dd class="col-sm-9">{{api.table}}</dd>
    <el-divider v-if="api.http">HTTP请求方式</el-divider>
    <dt v-if="api.http" class="col-sm-3 text-right">方法</dt>
    <dd v-if="api.http" class="col-sm-9">{{api.http.method}}</dd>
    <dt v-if="api.http" class="col-sm-3 text-right">路由</dt>
    <dd v-if="api.http" class="col-sm-9">{{api.http.route}}</dd>
    <el-divider v-if="api.timing">定时任务方式</el-divider>
    <dt v-if="api.timing" class="col-sm-3 text-right">定时方式</dt>
    <dd v-if="api.timing" class="col-sm-9">{{api.timing.type}}</dd>
    <dt v-if="api.timing" class="col-sm-3 text-right">定时规则</dt>
    <dd v-if="api.timing && (api.timing.type === 'interval' || api.timing.type === 'timeout')" class="col-sm-9">
        {{(api.timing.type === 'interval' ? "每" : "延时") + api.timing.mseconds}}毫秒执行
    </dd>
    <dd v-if="api.timing && api.timing.type === 'crontab'" class="col-sm-9">
        每天{{api.timing.hms}}准点执行
    </dd>
    <el-divider>请求/响应参数</el-divider>
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
