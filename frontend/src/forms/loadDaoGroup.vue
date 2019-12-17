<template>
<el-form :inline="true">
    <el-form-item label="语言">
        <el-select v-model="value" placeholder="请选择">
            <el-option v-for="lang in supportLangs" :key="lang.value" :label="lang.title" :value="lang.value"/>
        </el-select>
    </el-form-item>
    <el-form-item label="类别">
        <el-select v-model="value" placeholder="请选择">
            <el-option v-for="lang in supportLangs" :key="lang.value" :label="lang.title" :value="lang.value"/>
        </el-select>
    </el-form-item>
</el-form>
</template>

<script>
import backend from "../backend"

export default {
    data() {
        return {
            tempDaoGroups: {},
            supportLangs: [],
            supportCategories: []
        }
    },
    async created() {
        let res = await backend.qryAllDaoGroups()
        if (typeof res === "string") {
            this.$message.error(`查询所有模板DAO组时发生错误：${res}`)
        } else {
            for (let group of res.groups) {
                // TODO: 
            }
        }
    }
}
</script>