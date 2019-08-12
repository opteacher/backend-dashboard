<template>
<el-form :model="daoGroup" label-width="80px" ref="edit-dao-group-form">
    <el-form-item label="组名">
        <el-input placeholder="输入组名" v-model="daoGroup.name"/>
    </el-form-item>
    <el-form-item label="类别">
        <el-select class="w-100" v-model="daoGroup.category">
            <el-option v-for="category in categories" :key="category.value" :label="category.title" :value="category.value"/>
        </el-select>
    </el-form-item>
    <el-form-item label="实现"></el-form-item>
</el-form>
</template>

<script>
export default {
    data() {
        return {
            daoGroup: {
                name: "",
                category: ""
            },
            categories: [{
                title: "数据库",
                value: "database"
            }, {
                title: "缓存",
                value: "cache"
            }]
        }
    }
}
</script>

