<template>
<dashboard>
    <div class="dao-container">
        <el-button type="primary" @click="showAddDaoGroup = true">添加DAO组</el-button>
        <el-table class="mt-10" :data="daoGroups" style="width: 100%">
            <el-table-column type="expand">
                <template slot-scope="scope">
                    <el-table class="demo-table-expand" border :data="scope.row.interfaces">
                        <el-table-column label="接口名" prop="name"/>
                        <el-table-column label="参数" prop="name"/>
                        <el-table-column label="返回值" prop="name"/>
                        <el-table-column label="需包含的模块" prop="name"/>
                        <el-table-column label="描述" prop="name"/>
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
                <template>
                    <el-button size="mini">添加接口</el-button>
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
    </div>
</dashboard>
</template>

<script>
import backend from "../backend"
import dashboard from "../layouts/dashboard"
import editDaoGroup from "../forms/editDaoGroup"

export default {
    components: {
        "dashboard": dashboard,
        "edit-dao-group": editDaoGroup
    },
    data() {
        return {
            showAddDaoGroup: false,
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
        }
    }
}
</script>

<style lang="scss">
.dao-container {
    padding: 5vh 5vw;
}
</style>

