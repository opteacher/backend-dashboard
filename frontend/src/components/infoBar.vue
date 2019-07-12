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
            <el-button type="primary" @click="selectApi">确 定</el-button>
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import _ from "lodash"

import apisBkd from "../async/api"
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
        let res = await apisBkd.qry()
        if (typeof res === "string") {
            this.$message(`查询接口失败：${res}`)
        } else {
            let apis = res.data.data.infos || []
            if (apis.length === 0) {
                return
            }
            this.selApi = apis[0]
            this.addToRecentApis(this.selApi)
            this.$emit("select-api", this.selApi)
        }
    },
    methods: {
        addToRecentApis(api) {
            let apiInRec = this.recentApis.find(ele => ele.api.name === api.name)
            if (apiInRec) {
                apiInRec.num++
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
            let selApi = this.$refs["sel-api-form"].selApi
            if (selApi) {
                this.selApi = _.clone(selApi)
                this.addToRecentApis(this.selApi)
                this.$emit("select-api", selApi)
            }
            this.showSelApiDlg = false
        },
        hdlSelApi(apiName) {
            if (apiName === "*more") {
                this.showSelApiDlg = true
            } else {
                this.selApi = this.recentApis.find(ele => ele.api.name === apiName)
                this.selApi = this.selApi.api
                this.addToRecentApis(this.selApi)
                this.$emit("select-api", this.selApi)
            }
        },
    }
}
</script>

