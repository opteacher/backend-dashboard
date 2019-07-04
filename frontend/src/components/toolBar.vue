<template>
    <el-row class="toolbar">
        <el-col class="p-10" :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
        </el-col>
        <el-col class="p-10" :span="22">
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-plus" size="mini" @click="showAddMdlDlg = true"/>
                <el-button class="p-7" type="primary" icon="el-icon-share" size="mini" @click="showAddLnkDlg = true" :disabled="disableAddLnkBtn"/>
            </el-button-group>
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-download" size="mini" @click="showExportDlg = true"/>
            </el-button-group>
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
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="导出项目" :visible.sync="showExportDlg" :modal-append-to-body="false" width="50vw">
            <exp-project ref="exp-project-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showExportDlg = false">取 消</el-button>
                <el-button type="primary" @click="exportProject">导 出</el-button>
            </div>
        </el-dialog>
    </el-row>
</template>

<script>
import _ from "lodash"

import backend from "../async/backend"
import editModel from "../forms/editModel"
import editLink from "../forms/editLink"
import expProject from "../forms/expProject"

export default {
    props: {
        "models": Array
    },
    components: {
        "edit-model": editModel,
        "edit-link": editLink,
        "exp-project": expProject
    },
    data() { return {
        disableAddLnkBtn: true,
        showAddMdlDlg: false,
        showAddLnkDlg: false,
        showExportDlg: false
    }},
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
        async exportProject() {
            let form = this.$refs["exp-project-form"]
            form.$refs["exp-project-form"].validate(async valid => {
                if (valid) {
                    if (form.exportOption.name.slice(-4).toLowerCase() !== ".zip") {
                        form.exportOption.name += ".zip"
                    }
                    let res = await backend.export(form.exportOption)
                    if (res.data.data && res.data.data.url) {
                        window.open(res.data.data.url, "_blank")
                    }
                    this.showExportDlg = false
                } else {
                    return false
                }
            })
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
