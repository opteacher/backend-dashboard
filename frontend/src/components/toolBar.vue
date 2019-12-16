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
                <el-dropdown-item command="showAddMdlDlg">模块</el-dropdown-item>
                <el-dropdown-item command="showAddSttDlg">结构</el-dropdown-item>
                <el-dropdown-item command="showAddLnkDlg" v-show="!disableAddLnkBtn">关联</el-dropdown-item>
            </el-dropdown-menu>
        </el-dropdown>
        <el-badge :value="newStructsNum" class="item" :hidden="newStructsNum === 0" :max="5">
            <el-button size="mini" @click="hdlOpenSttLstDlg">结构列表</el-button>
        </el-badge>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="新建模块" :visible.sync="showAddMdlDlg" :modal-append-to-body="false" width="40vw">
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
    <el-dialog title="新建结构" :visible.sync="showAddSttDlg" :modal-append-to-body="false" width="40vw">
        <edit-model ref="add-struct-form" :structFlag="true"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddSttDlg = false">取 消</el-button>
            <el-button type="primary" @click="addStruct">确 定</el-button>
        </div>
    </el-dialog>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="结构列表" :visible.sync="showSttLstDlg" :modal-append-to-body="false" width="30vw">
        <list-struct ref="list-struct-form" :showFlag="showSttLstDlg"/>
        <div slot="footer" class="dialog-footer">
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import _ from "lodash"

import editModel from "../forms/editModel"
import editLink from "../forms/editLink"
import listStruct from "../forms/listStruct"
import backend from "../backend"

export default {
    props: {
        "models": Array
    },
    components: {
        "edit-model": editModel,
        "edit-link": editLink,
        "list-struct": listStruct
    },
    data() {
        return {
            checkboxGroup4: [],
            disableAddLnkBtn: true,
            showAddMdlDlg: false,
            showAddLnkDlg: false,
            showAddSttDlg: false,
            showSttLstDlg: false,
            newStructsNum: 0
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
                    let newModel = _.cloneDeep(form.model)
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
                    let newLink = _.cloneDeep(form.link)
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
        },
        addStruct() {
            let form = this.$refs["add-struct-form"]
            form.model.propName = "test"
            form.$refs["edit-model-form"].validate(async valid => {
                if (valid) {
                    let newStruct = _.cloneDeep(form.model)
                    newStruct = _.pick(newStruct, ["name", "props"])
                    let res = await backend.addModel(Object.assign(newStruct, {
                        type: "struct"
                    }))
                    if (typeof res === "string") {
                        this.$message.error(`添加结构时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "success",
                            message: "新增结构成功！"
                        })
                        form.resetModel()
                        this.newStructsNum++
                    }
                    this.showAddSttDlg = false
                } else {
                    return false
                }
            })
        },
        hdlOpenSttLstDlg() {
            this.newStructsNum = 0
            this.showSttLstDlg = true
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
