<template>
    <el-row class="toolbar">
        <el-col class="p-10" :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
        </el-col>
        <el-col class="p-10" :span="22">
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-plus" size="mini" @click="showAddMdlDlg = true"/>
            </el-button-group>
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-share" size="mini" @click="showAddLnkDlg = true" :disabled="disableAddLnkBtn"/>
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
            <edit-link ref="add-link-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showAddLnkDlg = false">取 消</el-button>
                <el-button type="primary" @click="addModel">确 定</el-button>
            </div>
        </el-dialog>
    </el-row>
</template>

<script>
import _ from "lodash"
import glbVar from "../global"

import editModel from "../forms/editModel"
import editLink from "../forms/editLink"

export default {
    components: {
        "edit-model": editModel,
        "edit-link": editLink,
    },
    data() { return {
        models: glbVar.models,
        disableAddLnkBtn: true,
        showAddMdlDlg: false,
        showAddLnkDlg: false
    }},
    watch: {
        models(newMdl, oldMdl) {
            this.disableAddLnkBtn = newMdl.length < 2
        }
    },
    methods: {
        async addModel() {
            this.showAddMdlDlg = false
            let form = this.$refs["add-model-form"]
            let newModel = _.clone(form.model)
            delete newModel.propName
            this.$emit("add-model", newModel)
            form.resetModel()
        },
        async addLink() {
            this.showAddLnkDlg = false
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
