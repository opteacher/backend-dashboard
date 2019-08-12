<template>
<dashboard>
    <div class="dao-container">
        <el-button type="primary" @click="showAddDaoGroup = true">添加DAO组</el-button>
        <el-table class="mt-10" :data="daoGroups" style="width: 100%">
            <el-table-column type="expand">
                <template slot-scope="scope">
                    <el-form label-position="left" inline class="demo-table-expand">
                        <el-form-item label="接口名">
                            <span>{{scope.row.name}}</span>
                        </el-form-item>
                    </el-form>
                </template>
            </el-table-column>
            <el-table-column label="组名" prop="name"/>
            <el-table-column label="类别" prop="category">
                <template slot-scope="scope">
                    <el-tag>{{scope.row.category}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="实现" prop="implement">
                <template slot-scope="scope">
                    <el-select v-model="scope.row.implement" size="mini" placeholder="请选择">
                        <el-option v-for="item in categories[scope.row.category]" :key="item.name" :label="item.name" :value="item.name"/>
                    </el-select>
                </template>
            </el-table-column>
            <el-table-column label="配置" prop="setting">
                <template>
                    <el-button size="mini">详细</el-button>
                </template>
            </el-table-column>
        </el-table>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="添加DAO组" :visible.sync="showAddDaoGroup" :modal-append-to-body="false" width="40vw">
            <edit-dao-group ref="add-dao-group-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showAddDaoGroup = false">取 消</el-button>
                <el-button type="primary">确 定</el-button>
            </div>
        </el-dialog>
    </div>
</dashboard>
</template>

<script>
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
            daoGroups: [{
                name: "abcd",
                category: "databases",
                implement: "",
                interfaces: [{
                    name: "SaveTx"
                }],
                setting: {}
            }],
            categories: {
                databases: [{
                    name: "cctv"
                }]
            }
        }
    },
    created() {

    }
}
</script>

<style lang="scss">
.dao-container {
    padding: 5vh 5vw;
}
</style>

