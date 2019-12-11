<template>
<dashboard>
    <div class="table-container">
        <el-button type="primary" @click="showAddDaoGroup = true">添加DAO组</el-button>
        <el-button type="primary" @click="showLoadDaoGroup = true">导入DAO组</el-button>
        <el-table class="mt-10" :data="daoGroups" style="width: 100%">
            <el-table-column type="expand">
                <template slot-scope="scope">
                    <el-table class="demo-table-expand" :data="scope.row.interfaces">
                        <el-table-column label="接口名" prop="name"/>
                        <el-table-column label="参数" prop="params">
                            <template slot-scope="subScope">
                                <div class="interval-container">
                                    <el-popover
                                        class="interval-item"
                                        v-for="param in subScope.row.params" :key="param.name"
                                        placement="top-start"
                                        trigger="hover"
                                        :content="param.type"
                                    >
                                        <el-tag size="small" slot="reference">{{param.name}}</el-tag>
                                    </el-popover>
                                </div>
                            </template>
                        </el-table-column>
                        <el-table-column label="返回值" prop="returns"/>
                        <el-table-column label="依赖模块" prop="requires"/>
                        <el-table-column label="描述" prop="desc"/>
                        <el-table-column label="配置" prop="setting">
                            <template slot-scope="subScope">
                                <el-button size="mini" type="danger" @click="delDaoInterface(scope.row.name, subScope.row.name)">删除</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </template>
            </el-table-column>
            <el-table-column label="组名" prop="name"/>
            <el-table-column label="类别" prop="category">
                <template slot-scope="scope">
                    <el-tag type="info">{{scope.row.category}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="语言" prop="language">
                <template slot-scope="scope">
                    <el-tag>{{scope.row.language}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="实现" prop="implement">
                <template slot-scope="scope">
                    <el-button type="text">{{scope.row.implement}}</el-button>
                </template>
            </el-table-column>
            <el-table-column label="配置" prop="setting">
                <template slot-scope="scope">
                    <el-button size="mini" @click="editingGroup = scope.row">添加接口</el-button>
                    <el-button size="mini" type="danger" @click="delDaoGroup(scope.row.name)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="添加DAO组" :visible.sync="showAddDaoGroup" :modal-append-to-body="false" width="40vw">
            <edit-dao-group ref="add-dao-group-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showAddDaoGroup = false">取 消</el-button>
                <el-button type="primary" @click="addDaoGroup">确 定</el-button>
            </div>
        </el-dialog>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="添加DAO接口" :visible="editingGroup !== null" :modal-append-to-body="false" width="40vw" @close="editingGroup = null">
            <edit-dao-interface ref="add-dao-interface-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="editingGroup = null">取 消</el-button>
                <el-button type="primary" @click="addDaoInterface">确 定</el-button>
            </div>
        </el-dialog>
    </div>
</dashboard>
</template>

<script>
import backend from "../backend"
import dashboard from "../layouts/dashboard"
import editDaoGroup from "../forms/editDaoGroup"
import editDaoInterface from "../forms/editDaoInterface"

export default {
    components: {
        "dashboard": dashboard,
        "edit-dao-group": editDaoGroup,
        "edit-dao-interface": editDaoInterface
    },
    data() {
        return {
            showAddDaoGroup: false,
            showLoadDaoGroup: false,
            editingGroup: null,
            daoGroups: [],
            categories: {
                databases: [{
                    name: "cctv"
                }]
            }
        }
    },
    async created() {
        await this.refresh()
    },
    methods: {
        async refresh() {
            let res = await backend.qryAllDaoGroups()
            if (typeof res === "string") {
                this.$message.error(`查询所有DAO组时发生错误：${res}`)
            } else {
                this.daoGroups = res.groups
            }
        },
        addDaoGroup() {
            let addForm = this.$refs["add-dao-group-form"]
            addForm.$refs["edit-dao-group-form"].validate(async valid => {
                if (!valid) {
                    return false
                }
                let form = addForm.$refs["edit-dao-group-form"]
                let res = await backend.addDaoGroup(form.model)
                if (typeof res === "string") {
                    this.$message.error(`添加DAO组发生错误：${res}`)
                } else {
                    this.showAddDaoGroup = false
                    await this.refresh()
                }
            })
        },
        addDaoInterface() {
            let addForm = this.$refs["add-dao-interface-form"]
            addForm.$refs["edit-dao-interface-form"].validate(async valid => {
                if (!valid) {
                    return false
                }
                let form = addForm.$refs["edit-dao-interface-form"]
                let res = await backend.addDaoInterface({
                    gpname: this.editingGroup.name,
                    interface: form.model
                })
                if (typeof res === "string") {
                    this.$message.error(`添加DAO接口发生错误：${res}`)
                } else {
                    this.editingGroup = null
                    await this.refresh()
                }
            })
        },
        delDaoInterface(gpname, ifname) {
            this.$alert("确定删除接口？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delDaoInterface({gpname, ifname})
                    if (typeof res === "string") {
                        this.$message.error(`删除接口时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `接口（${ifname}）删除成功！`
                        })
                        await this.refresh()
                    }
                }
            })
        },
        delDaoGroup(gpname) {
            this.$alert("确定删除组？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delDaoGroup(gpname)
                    if (typeof res === "string") {
                        this.$message.error(`删除DAO组时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `DAO组（${gpname}）删除成功！`
                        })
                        await this.refresh()
                    }
                }
            })
        }
    }
}
</script>

