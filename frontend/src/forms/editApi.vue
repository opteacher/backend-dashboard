<template>
<el-form ref="form" :model="api" label-width="80px">
    <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="base">
            <el-form-item label="接口名称" :rules="[
                { required: true, message: '请输入接口名称', trigger: 'blur' },
                { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
            ]" prop="name">
                <el-input v-model="api.name"/>
            </el-form-item>
            <el-form ref="add-param-form" :model="newParam" label-width="80px">
                <el-form-item label="新增参数" :rules="[
                    { required: true, message: '请输入参数名称', trigger: 'blur' },
                    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
                ]" prop="paramName">
                    <el-input v-model="newParam.paramName"/>
                </el-form-item>
                <el-form-item label="参数类型" :rules="[
                    { required: true, message: '请输入参数类型', trigger: 'blur' }
                ]" prop="paramType">
                    <el-col :span="18">
                        <el-select class="w-100" v-model="newParam.paramType" placeholder="请选择">
                            <el-option v-for="typ in types" :key="typ.name" :label="typ.name" :value="typ.name"/>
                        </el-select>
                    </el-col>
                    <el-col :span="6">
                        <el-button class="float-right" @click="hdlAddParam">添加参数</el-button>
                    </el-col>
                </el-form-item>
            </el-form>
            <el-form-item label="参数表" v-if="api.params.length !== 0">
                <el-table :data="api.params" :show-header="false">
                    <el-table-column prop="type" width="100">
                        <template slot-scope="scope">
                            <el-tag size="medium">{{scope.row.type}}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="name"/>
                    <el-table-column width="80">
                        <template slot-scope="scope">
                            <el-button circle type="danger" size="mini" @click.native.prevent="removeParam(scope.row.name)" icon="el-icon-minus"/>
                        </template>
                    </el-table-column>
                </el-table>
            </el-form-item>
            <el-form :model="newReturn" ref="add-return-form" label-width="80px">
                <el-form-item label="返回类型" :rules="[
                    { required: true, message: '请选择返回值类型', trigger: 'blur' }
                ]" prop="type">
                    <el-col :span="18">
                        <el-select class="w-100" v-model="newReturn.type" placeholder="请选择">
                            <el-option v-for="typ in types" :key="typ.name" :label="typ.name" :value="typ.name"/>
                        </el-select>
                    </el-col>
                    <el-col :span="6">
                        <el-button class="float-right" @click="hdlAddReturn">添加返回值</el-button>
                    </el-col>
                </el-form-item>
            </el-form>
            <el-form-item label="返回值" v-if="api.returns.length !== 0">
                <el-tag class="mr-2" v-for="ret in api.returns" :key="ret" closable :disable-transitions="false" @close="hdlRemoveRet(ret)">
                    {{ret}}
                </el-tag>
            </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="激活方式" name="active">
            <el-form-item label="激活方式">
                <el-select class="w-100" v-model="activeType" placeholder="请选择">
                    <el-option label="普通/RPC接口" value="interface"/>
                    <el-option label="HTTP接口" value="http"/>
                    <el-option label="定时任务" value="timing"/>
                    <el-option label="订阅频道" value="subscribe"/>
                </el-select>
            </el-form-item>
            <div v-show="activeType == 'http'">
                <el-form-item label="路由">
                    <el-input v-model="api.http.route"/>
                </el-form-item>
                <el-form-item label="方法">
                    <el-select class="w-100" v-model="api.http.method">
                        <el-option label="GET" value="GET"/>
                        <el-option label="POST" value="POST"/>
                        <el-option label="PUT" value="PUT"/>
                        <el-option label="DELETE" value="DELETE"/>
                        <el-option label="PATCH" value="PATCH"/>
                    </el-select>
                </el-form-item>
            </div>
            <div v-show="activeType == 'timing'">
                <el-form-item label="定时方式">
                    <el-select class="w-100" v-model="api.timing.type" placeholder="请选择" @change="hdlChgTimeType">
                        <el-option label="间隔" value="interval"/>
                        <el-option label="延时" value="timeout"/>
                        <el-option label="定期" value="crontab"/>
                    </el-select>
                </el-form-item>
                <div v-show="api.timing.type == 'interval'">
                    <el-form-item label="间隔执行">
                        <el-input v-model="tempTime[0]" class="input-with-select" text="number">
                            <el-select slot="append" v-model="tempTime[1]">
                                <el-option label="天" value="28800000"/>
                                <el-option label="时" value="1200000"/>
                                <el-option label="分" value="60000"/>
                                <el-option label="秒" value="1000" />
                                <el-option label="毫秒" value="1"/>
                            </el-select>
                        </el-input>
                    </el-form-item>
                </div>
                <div v-show="api.timing.type == 'timeout'">
                    <el-form-item label="延时执行">
                        <el-input v-model="tempTime[0]" class="input-with-select" text="number">
                            <el-select slot="append" v-model="tempTime[1]">
                                <el-option label="天" value="28800000"/>
                                <el-option label="时" value="1200000"/>
                                <el-option label="分" value="60000"/>
                                <el-option label="秒" value="1000"/>
                                <el-option label="毫秒" value="1"/>
                            </el-select>
                        </el-input>
                    </el-form-item>
                </div>
                <div v-show="api.timing.type == 'crontab'">
                    <el-form-item label="定期执行">
                        <el-time-select class="w-100" v-model="api.timing.hms" :picker-options="{step: '00:15'}" placeholder="选择时间"/>
                    </el-form-item>
                </div>
            </div>
            <div v-show="activeType == 'subscribe'">
                <el-form-item label="频道">
                    <el-col :span="13">
                        <el-input v-model="api.subscribe.channel"/>
                    </el-col>
                    <el-col :span="11">
                        <el-select class="float-right" v-model="api.subscribe.daoGroup" placeholder="请选择通讯DAO">
                            <el-option v-for="daoGrp in subscbDaos" :key="daoGrp.name" :label="daoGrp.name" :value="daoGrp"/>
                        </el-select>
                    </el-col>
                </el-form-item>
            </div>
        </el-tab-pane>
        <el-tab-pane label="其他" name="others">
            <el-form-item label="绑定模块">
                <el-select class="w-100" v-model="api.model" placeholder="可选">
                    <el-option v-for="model in models" :key="model.name" :label="model.name" :value="model.name"/>
                </el-select>
            </el-form-item>
        </el-tab-pane>
    </el-tabs>
</el-form>
</template>

<script>
import backend from "../backend"

export default {
    data() {
        return {
            activeTab: "base",
            api: {
                name: "",
                params: [],
                returns: [],
                http: {
                    route: "/",
                    method: "GET"
                },
                timing: {
                    type: "interval",
                    mseconds: 0,
                    hms: ""
                },
                subscribe: {
                    channel: "",
                    daoGroup: ""
                },
                model: ""
            },
            newParam: {
                paramName: "",
                paramType: ""
            },
            newReturn: {
                type: ""
            },
            models: [],
            types: [],
            activeType: "interface",
            tempTime: [0, "1000"],
            subscbDaos: []
        }
    },
    async created() {
        let res = await backend.qryAllModels()
        if (typeof res === "string") {
            this.$message.error(`查询模块失败：${res}`)
        } else {
            this.models = res.models || []
        }
        res = await backend.qryAllAvaTypes()
        if (typeof res === "string") {
            this.$message.error(`查询可用类型失败：${res}`)
        } else {
            this.types = res.models || []
        }
    },
    methods: {
        hdlAddParam() {
            this.$refs["add-param-form"].validate(valid => {
                if (valid) {
                    this.api.params.push({
                        name: this.newParam.paramName,
                        type: this.newParam.paramType
                    })
                    this.newParam.paramName = ""
                    this.newParam.paramType = ""
                } else {
                    return false
                }
            })
        },
        removeParam(pname) {
            for (let i = 0; i < this.api.params.length; i++) {
                if (this.api.params[i].name !== pname) {
                    continue
                }
                this.api.params = this.api.params.slice(0, i).concat(
                    this.api.params.slice(i + 1)
                )
                break
            }
        },
        hdlAddReturn() {
            this.$refs["add-return-form"].validate(valid => {
                if (valid) {
                    this.api.returns.push(this.newReturn.type)
                    this.newReturn.type = ""
                } else {
                    return false
                }
            })
        },
        hdlRemoveRet(ret) {
            this.api.returns.splice(this.api.returns.indexOf(ret), 1)
        },
        hdlChgTimeType() {
            if (this.timing.type === "interval" || this.timing.type === "timeout") {
                this.timing.hms = ""
            } else {
                this.timing.msecond = 0
            }
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
