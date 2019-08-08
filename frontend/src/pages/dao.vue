<template>
<dashboard>
    <div class="dao-container">
        <el-button type="primary">添加接口</el-button>
        <el-table :data="daoGroups" style="width: 100%">
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
            <el-table-column label="类别" prop="category"/>
            <el-table-column label="实现" prop="implement"/>
        </el-table>
    </div>
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"

export default {
    components: {
        "dashboard": dashboard
    },
    data() {
        return {
            daoGroups: [{
                name: "abcd",
                category: "数据源",
                implement: "无",
                interfaces: [{
                    name: "SaveTx"
                }]
            }]
        }
    }
}
</script>

<style lang="scss">
.dao-container {
    padding: 5vh 5vw;
}
</style>

