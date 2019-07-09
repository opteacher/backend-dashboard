<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <el-button-group class="p-0">
            <el-button class="p-7" icon="el-icon-plus" size="mini"/>
        </el-button-group>
        <el-button-group class="p-0">
            <el-input class="input-with-select" placeholder="未选定接口" v-model="selItf.name" size="mini" disabled>
                <el-button class="p-7" icon="el-icon-menu" size="mini" slot="prepend" @click="showSelItfDlg = true"/>
                <el-button class="p-7" icon="el-icon-warning" size="mini" slot="append" :disabled="selItf.name.length === 0" @click="showItfInfo"/>
            </el-input>
        </el-button-group>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="选择接口" :visible.sync="showSelItfDlg" :modal-append-to-body="false" width="40vw">
        <sel-interface ref="sel-itf-form"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showSelItfDlg = false">取 消</el-button>
            <el-button type="primary" @click="selInterface">确 定</el-button>
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import _ from "lodash"

import apisBkd from "../async/api"
import selInterface from "../forms/selInterface"
import interfaceInfo from "../forms/interfaceInfo"

export default {
    components: {
        "sel-interface": selInterface
    },
    data() {
        return {
            showSelItfDlg: false,
            selItf: {
                name: ""
            },
            index: 1
        }
    },
    async created() {
        let res = await apisBkd.qry()
        if (typeof res === "string") {
            this.$message(`查询接口失败：${res}`)
        } else {
            let interfaces = res.data.data.infos || []
            if (interfaces.length === 0) {
                return
            }
            this.selItf = interfaces[0]
            this.$emit("sel-interface", this.selItf)
        }
    },
    methods: {
        selInterface() {
            let selItf = this.$refs["sel-itf-form"].selItf
            if (selItf) {
                this.selItf = _.clone(selItf)
                this.$emit("sel-interface", selItf)
            }
            this.showSelItfDlg = false
        },
        showItfInfo() {
            this.index++
            this.$msgbox({
                title: "接口信息",
                message: this.$createElement(interfaceInfo, {
                    props: {
                        interface: this.selItf,
                    },
                    key: this.index
                }),
                showConfirmButton: false
            })
        }
    }
}
</script>

