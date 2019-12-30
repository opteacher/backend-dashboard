<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <div v-show="qryApiFun === 'qryAllApis'">
            <el-dropdown split-button trigger="click" size="mini" type="primary" @click="showAddApiDlg = true" @command="hdlSelApi">
                添加接口
                <el-dropdown-menu slot="dropdown" >
                    <el-dropdown-item v-for="api in recentApis" :key="api.api.name" :command="api.api.name">
                        {{api.api.name}}
                    </el-dropdown-item>
                    <el-dropdown-item command="*more" icon="el-icon-more"/>
                </el-dropdown-menu>
            </el-dropdown>
            <el-button size="mini" @click="hdlDelApi" :disabled="selApi.name.length === 0">删除接口</el-button>
            <el-button size="mini" @click="showTempApis(true)">显示模板接口</el-button>
        </div>
        <el-dropdown v-show="qryApiFun === 'qryAllTempApis'" split-button trigger="click" size="mini" @click="showTempApis(false)" @command="hdlSelApi">
            隐藏模板接口
            <el-dropdown-menu slot="dropdown" >
                <el-dropdown-item v-for="api in recentApis" :key="api.api.name" :command="api.api.name">
                    {{api.api.name}}
                </el-dropdown-item>
                <el-dropdown-item command="*more" icon="el-icon-more"/>
            </el-dropdown-menu>
        </el-dropdown>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="选择接口" :visible.sync="showSelApiDlg" :modal-append-to-body="false" width="70vw">
        <select-api ref="sel-api-form" :showFlag="showSelApiDlg" :qryApiFun="qryApiFun"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showSelApiDlg = false">取 消</el-button>
            <el-button type="primary" @click="selectApi">确 定</el-button>
        </div>
    </el-dialog>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="新建接口" :visible.sync="showAddApiDlg" :modal-append-to-body="false" width="40vw">
        <edit-api ref="add-api-form"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddApiDlg = false">取 消</el-button>
            <el-button type="primary" @click="addApi">确 定</el-button>
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import _ from "lodash"
import backend from "../backend"
import selectApi from "../forms/selectApi"
import editApi from "../forms/editApi"

export default {
    components: {
        "select-api": selectApi,
        "edit-api": editApi
    },
    data() {
        return {
            qryApiFun: "qryAllApis",
            showSelApiDlg: false,
            showAddApiDlg: false,
            selApi: {name: ""},
            recentApis: [],
            index: 1
        }
    },
    async created() {
        await this.refresh()
    },
    methods: {
        async refresh() {
            this.recentApis = []
            let res = await backend[this.qryApiFun]()
            if (typeof res === "string") {
                this.$message.error(`查询接口失败：${res}`)
            } else {
                let apis = res.infos || []
                if (apis.length === 0) {
                    this.selApiLocal({name: ""})
                } else {
                    this.selApiLocal(apis[0])
                }
            }
        },
        selApiLocal(api) {
            this.selApi = api
            this.addToRecentApis(this.selApi)
            this.$emit("select-api", this.selApi)
        },
        addToRecentApis(api) {
            let holder = this.recentApis.find(ele => ele.api.name === api.name)
            if (holder) {
                holder.num++
                return
            } else if (this.recentApis.length > 5) {
                // 找出num最小的api，从最近api表中剔除
                let minIdx = 0
                for (let i = 1; i < this.recentApis.length; i++) {
                    if (this.recentApis[i].num < this.recentApis[minIdx].num) {
                        minIdx = i
                    }
                }
                this.recentApis = this.recentApis.slice(0, minIdx)
                    .concat(this.recentApis.slice(minIdx + 1))
            }
            this.recentApis.push({
                api: api,
                num: 1
            })
        },
        selectApi() {
            this.showSelApiDlg = false
            let selApi = this.$refs["sel-api-form"].selApi
            if (selApi) {
                this.selApiLocal(selApi)
            }
        },
        hdlSelApi(apiName) {
            if (apiName === "*more") {
                this.showSelApiDlg = true
            } else {
                let holder = this.recentApis.find(ele => ele.api.name === apiName)
                this.selApiLocal(holder.api)
            }
        },
        addApi() {
            let form = this.$refs["add-api-form"]
            form.$refs["form"].validate(async valid => {
                if (!valid) {
                    return false
                }
                // 根据需要消除http相关信息
                let api = _.cloneDeep(form.api)
                // 删除不需要的激活方式
                switch (form.activeType) {
                    case "http":
                        delete api.timing
                        delete api.subscribe
                        break
                    case "timing":
                        delete api.http
                        delete api.subscribe
                        if (api.timing.type == "interval" || api.timing.type == "timeout") {
                            api.timing.mseconds = form.tempTime[0] * Number(form.tempTime[1])
                            delete api.timing.hms
                        } else {
                            delete api.timing.mseconds
                        }
                        break
                    case "subscribe":
                        delete api.http
                        delete api.timing
                        break
                    case "interface":
                    default:
                        delete api.http
                        delete api.timing
                        delete api.subscribe
                }
                // 将API信息发送给后台
                let res = await backend.addApi(api)
                if (typeof res === "string") {
                    this.$message.error(`添加接口失败：${res}`)
                } else {
                    this.selApiLocal(res)
                }
                this.showAddApiDlg = false
            })
        },
        hdlDelApi() {
            this.$alert("确定删除接口？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delApiByName(this.selApi.name)
                    if (typeof res === "string") {
                        this.$message.error(`删除接口时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `接口（${this.selApi.name}）删除成功！`
                        })
                        await this.refresh()
                        this.$emit("del-api")
                    }
                }
            })
        },
        async showTempApis(show) {
            this.qryApiFun = show ? "qryAllTempApis" : "qryAllApis"
            await this.refresh()
        }
    }
}
</script>

