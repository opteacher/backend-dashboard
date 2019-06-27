<template>
    <el-row class="toolbar">
        <el-col class="p-10" :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
        </el-col>
        <el-col class="p-10" :span="22">
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-plus" size="mini" @click="showAddMdlDlg = true"/>
                <el-button class="p-7" type="primary" icon="el-icon-share" size="mini" @click="showAddRelDlg = true" :disabled="disableAddRelBtn"/>
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
        <el-dialog title="新建关联" :visible.sync="showAddRelDlg" :modal-append-to-body="false" width="50vw">
            <edit-relation ref="add-relation-form" :models="models"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="resetRelation">重 置</el-button>
                <el-button @click="showAddRelDlg = false">取 消</el-button>
                <el-button type="primary" @click="addRelation">确 定</el-button>
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
import editRelation from "../forms/editRelation"
import expProject from "../forms/expProject"

export default {
    props: {
        "models": Array
    },
    components: {
        "edit-model": editModel,
        "edit-relation": editRelation,
        "exp-project": expProject
    },
    data() { return {
        disableAddRelBtn: true,
        showAddMdlDlg: false,
        showAddRelDlg: false,
        showExportDlg: false
    }},
    watch: {
        models(nv, ov) {
            this.chkAddRelBtn()
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
        resetRelation() {
            this.$refs["add-relation-form"].resetRelation()
        },
        async addRelation() {
            let form = this.$refs["add-relation-form"]
            form.$refs["edit-relation-form"].validate(valid => {
                if (valid) {
                    let newRelation = _.clone(form.relation)
                    this.$emit("add-relation", newRelation)
                    form.resetRelation()
                    this.showAddRelDlg = false
                } else {
                    return false
                }
            })
        },
        chkAddRelBtn() {
            this.disableAddRelBtn = !this.models || this.models.length < 2
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
