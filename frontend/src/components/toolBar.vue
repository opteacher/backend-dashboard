<template>
    <el-row class="toolbar">
        <el-col :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
        </el-col>
        <el-col :span="22">
            <el-button-group class="p-0">
                <el-button class="p-7" type="primary" icon="el-icon-plus" size="mini" @click="showAddModelDlg = true"/>
            </el-button-group>
        </el-col>
        <el-col :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
        </el-col>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="新建模块" :visible.sync="showAddModelDlg" :modal-append-to-body="false" width="50vw">
            <edit-model ref="add-model-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showAddModelDlg = false">取 消</el-button>
                <el-button type="primary" @click="addModel">确 定</el-button>
            </div>
        </el-dialog>
    </el-row>
</template>

<script>
import _ from "lodash"

import editModel from "../forms/editModel"

export default {
    components: {
        "edit-model": editModel
    },
    data() { return {
        showAddModelDlg: false
    }},
    methods: {
        async addModel() {
            this.showAddModelDlg = false
            let form = this.$refs["add-model-form"]
            let newModel = _.clone(form.model)
            try {
                console.log(await this.axios.post("http://127.0.0.1:8000/backend-dashboard/backend/models", newModel))
            } catch(e) {
                this.$message(`创建模块失败：${e}`)
            }
            this.$emit("add-model", newModel)
            form.resetModel()
        }
    }
}
</script>


<style lang="scss">
.toolbar {
    background-color: white;

    .el-col {
        padding: 10px;
    }

    button i {
        padding: 0;
    }
}
</style>
