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
            <el-table-column label="所属" property="model" width="100"/>
            <el-table-column label="HTTP方法" property="method" width="100"/>
            <el-table-column label="HTTP路径" property="route"/>
        </el-table>
    </el-form-item>
</el-form>
</template>

<script>
import apisBkd from "../async/api"

export default {
    data() {
        return {
            searchTxt: "",
            itfList: [],
            selItf: null,
        }
    },
    async created() {
        await this.queryApis()
    },
    methods: {
        handleCurrentChange(selItf) {
            this.selItf = selItf
        },
        async queryApis() {
            let res = await apisBkd.qry()
            if (typeof res === "string") {
                this.$message(`查询接口失败：${res}`)
            } else {
                this.itfList = res.data.data.infos || []
            }
        }
    }
}
</script>

