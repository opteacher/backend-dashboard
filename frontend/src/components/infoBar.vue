<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <el-dropdown split-button trigger="click" size="mini" type="primary" @click="showAddApiDlg = true" @command="hdlSelApi">
            添加接口
            <el-dropdown-menu slot="dropdown" >
                <el-dropdown-item v-for="api in recentApis" :key="api.api.name" :command="api.api.name">
                    {{api.api.name}}
                </el-dropdown-item>
                <el-dropdown-item command="*more">...</el-dropdown-item>
            </el-dropdown-menu>
        </el-dropdown>
        <el-button size="mini" @click="hdlDelApi">删除接口</el-button>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="选择接口" :visible.sync="showSelApiDlg" :modal-append-to-body="false" width="40vw">
        <select-api ref="sel-api-form"/>
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
            showSelApiDlg: false,
            showAddApiDlg: false,
            selApi: {
                name: ""
            },
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
            let res = await backend.qryAllApis()
            if (typeof res === "string") {
                this.$message.error(`查询接口失败：${res}`)
            } else {
                let apis = res.infos || []
                if (apis.length === 0) {
                    return
                }
                this.selApiLocal(apis[0])
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
            if (!form.api.enableHttp) {
                form.api.route = "/"
                form.api.method = "GET"
            }
            form.$refs["form"].validate(async valid => {
                if (valid) {
                    // 根据需要消除http相关信息
                    let api = _.cloneDeep(form.api)
                    if (!api.enableHttp) {
                        delete(api.route)
                        delete(api.method)
                    }
                    delete(api.enableHttp)
                    // 将params转化成一个对象
                    let params = {}
                    api.params.map(param => {
                        params[param.name] = param.type
                    })
                    api.params = params
                    // 将API信息发送给后台
                    let res = await backend.addApi(api)
                    if (typeof res === "string") {
                        this.$message.error(`添加接口失败：${res}`)
                    } else {
                        this.selApiLocal(res)
                    }
                    this.showAddApiDlg = false
                } else {
                    return false
                }
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
                            message: `模块（${this.selApi.name}）删除成功！`
                        })
                        await this.refresh()
                    }
                }
            })
        }
    }
}
</script>

