<template>
<el-form ref="form" label-width="80px">
    <el-form-item label-width="0">
        <el-input v-model="searchTxt">
            <i class="el-icon-search el-input__icon" slot="prefix"></i>
        </el-input>
    </el-form-item>
    <el-form-item label-width="0">
        <el-table :data="itfList" highlight-current-row @current-change="handleCurrentChange" style="width: 100%">
            <el-table-column label="接口名" property="name"/>
            <el-table-column label="所属" property="model"/>
            <el-table-column label="HTTP方法" property="method"/>
            <el-table-column label="HTTP路径" property="route"/>
        </el-table>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    data() {
        return {
            searchTxt: "",
            itfList: [{
                name: "abcd",
                model: "User",
                method: "POST",
                route: "/api/v1/users"
            }]
        }
    },
    methods: {
        handleCurrentChange() {
            console.log(arguments)
        }
    }
}
</script>

