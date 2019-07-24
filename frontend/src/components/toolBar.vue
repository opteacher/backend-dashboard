<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <el-dropdown trigger="click" @command="hdlSelAdd">
            <el-button size="mini" type="primary">
                添加组件<i class="el-icon-arrow-down el-icon--right"></i>
            </el-button>
            <el-dropdown-menu slot="dropdown">
                <el-dropdown-item icon="el-icon-menu" command="showAddMdlDlg">模块</el-dropdown-item>
                <el-dropdown-item icon="el-icon-share" command="showAddLnkDlg" v-show="!disableAddLnkBtn">关联</el-dropdown-item>
            </el-dropdown-menu>
        </el-dropdown>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="新建模块" :visible.sync="showAddMdlDlg" :modal-append-to-body="false" width="50vw">
        <edit-model ref="add-model-form"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddMdlDlg = false">取 消</el-button>
            <el-button type="primary" @click="addModel">确 定</el-button>
        </div>
    </el-dialog>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="新建关联" :visible.sync="showAddLnkDlg" :modal-append-to-body="false" width="50vw">
        <edit-link ref="add-link-form" :models="models"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="resetLink">重 置</el-button>
            <el-button @click="showAddLnkDlg = false">取 消</el-button>
            <el-button type="primary" @click="addLink">确 定</el-button>
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import _ from "lodash"

import editModel from "../forms/editModel"
import editLink from "../forms/editLink"

export default {
    props: {
        "models": Array
    },
    components: {
        "edit-model": editModel,
        "edit-link": editLink
    },
    data() {
        return {
            disableAddLnkBtn: true,
            showAddMdlDlg: false,
            showAddLnkDlg: false
        }
    },
    watch: {
        models(nv, ov) {
            this.chkAddLnkBtn()
        }
    },
    methods: {
        async addModel() {
            let form = this.$refs["add-model-form"]
            form.model.propName = "test"
            form.$refs["edit-model-form"].validate(valid => {
                if (valid) {
                    let newModel = _.clone(form.model)
                    delete newModel.propName
                    this.$emit("add-model", newModel)
                    form.resetModel()
                    this.showAddMdlDlg = false
                } else {
                    return false
                }
            })
        },
        resetLink() {
            this.$refs["add-link-form"].resetLink()
        },
        async addLink() {
            let form = this.$refs["add-link-form"]
            form.$refs["edit-link-form"].validate(valid => {
                if (valid) {
                    let newLink = _.clone(form.link)
                    this.$emit("add-link", newLink)
                    form.resetLink()
                    this.showAddLnkDlg = false
                } else {
                    return false
                }
            })
        },
        chkAddLnkBtn() {
            this.disableAddLnkBtn = !this.models || this.models.length < 2
        },
        hdlSelAdd(cmd) {
            this[cmd] = true
        }
    }
}
</script>


<style lang="scss">
.toolbar {
    background-color: white;
    button i {
        padding: 0;
    }
}
</style>
