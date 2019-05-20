<template>
    <el-row class="toolbar">
        <el-col :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
        </el-col>
        <el-col :span="22">
            <el-button-group class="p-0">
                <el-button type="primary" icon="el-icon-plus" size="mini" @click="addModel"/>
            </el-button-group>
        </el-col>
        <el-col :span="1">
            <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
        </el-col>
    </el-row>
</template>

<script>
export default {
    methods: {
        async addModel() {
            const h = this.$createElement;
            try {
                let action = await this.$msgbox({
                    title: "新建",
                    message: h("p", "新建一个模块"),
                    showCancelButton: true,
                    confirmButtonText: "确认",
                    cancelButtonText: "取消"
                })
                this.$message({
                    type: "info",
                    message: `action: ${action}`
                })
            } catch (e) {
                this.$message({
                    type: "error",
                    message: `error: ${e.message || e}`
                })
            }
        }
    }
}
</script>


<style lang="scss">
.toolbar {
    margin-top: -20px;
    margin-left: -20px;
    margin-right: -20px;
    margin-bottom: 20px;
    background-color: white;

    * {
        padding: 10px;
    }

    button i {
        padding: 0;
    }
}
</style>
