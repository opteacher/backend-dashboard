<template>
<dashboard>
    <div class="dao-container">
        <el-button type="primary">添加DAO组</el-button>
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
                    <el-select v-model="value" placeholder="请选择">
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
                category: "databases",
                implement: "无",
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

