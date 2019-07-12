<template>
<el-form ref="form" :model="api" label-width="80px">
    <el-form-item label="接口名称">
        <el-input v-model="api.name"></el-input>
    </el-form-item>
    <el-form-item label="新增参数">
        <el-input v-model="paramName"/>
    </el-form-item>
    <el-form-item label="参数类型">
        <el-col :span="18">
            <el-select class="w-100" v-model="paramType" placeholder="请选择">
                <el-option v-for="typ in typeMap" :key="typ.value" :label="typ.title" :value="typ.value"/>
            </el-select>
        </el-col>
        <el-col :span="6">
            <el-button class="float-right" @click="hdlAddParam">添加参数</el-button>
        </el-col>
    </el-form-item>
    <el-form-item label="参数表" v-if="api.params.length !== 0">
        <el-table :data="api.params" :show-header="false">
            <el-table-column prop="type" width="40%">
                <template slot-scope="scope">
                    <el-tag size="medium">{{scope.row.type}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="name"/>
        </el-table>
    </el-form-item>
    <el-form-item label="HTTP">
        <el-switch v-model="enableHttp"/>
    </el-form-item>
    <el-form-item v-show="enableHttp" label="路由">
        <el-input v-model="api.route"/>
    </el-form-item>
    <el-form-item v-show="enableHttp" label="方法">
        <el-select class="w-100" v-model="api.method" placeholder="请选择">
            <el-option v-for="method in methodMap" :key="method" :label="method" :value="method"/>
        </el-select>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    data() {
        return {
            typeMap: [{
                title: "文本",
                value: "string"
            }, {
                title: "数字",
                value: "int32"
            }, {
                title: "日期",
                value: "uint64"
            }, {
                title: "布尔",
                value: "bool"
            }],
            methodMap: [
                "GET", "POST", "PUT", "DELETE", "PATCH"
            ],
            api: {
                name: "",
                params: [],
                route: "",
                method: ""
            },
            paramName: "",
            paramType: "",
            enableHttp: true,
        }
    },
    methods: {
        hdlAddParam() {
            this.api.params.push({
                name: this.paramName,
                type: this.paramType
            })
        }
    }
}
</script>

<style lang="scss">
.btn-append-input:hover, .btn-append-input:focus {
    color: #409eff;
    border-color: #c6e2ff;
    background-color: #ecf5ff;
}
.btn-append-input:active {
    color: #3a8ee6;
    border-color: #3a8ee6;
    outline: none;
}
</style>
