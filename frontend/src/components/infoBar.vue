<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <el-dropdown size="mini" split-button type="primary">
            超小尺寸
            <el-dropdown-menu slot="dropdown">
                <el-dropdown-item>黄金糕</el-dropdown-item>
                <el-dropdown-item>狮子头</el-dropdown-item>
                <el-dropdown-item>螺蛳粉</el-dropdown-item>
                <el-dropdown-item>双皮奶</el-dropdown-item>
                <el-dropdown-item>蚵仔煎</el-dropdown-item>
            </el-dropdown-menu>
        </el-dropdown>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
</el-row>
</template>

<script>
import modelBkd from "../async/model"

export default {
    data() {
        return {

        }
    },
    created() {
        let res = await modelBkd.qry()
        if (typeof res === "string") {
            this.$message(`查询模块失败：${res}`)
        } else {
            this.models = (res.data.data && res.data.data.models) || []
        }
    }
}
</script>

