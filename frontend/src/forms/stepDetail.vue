<template>
<el-form ref="form" label-width="80px">
    <el-form-item label="操作标识">
        {{selStep.operKey}}
    </el-form-item>
    <el-form-item label="依赖" v-show="selStep.requires && selStep.requires.length !== 0">
        <el-tag v-for="require in selStep.requires" :key="require" :closable="mode !== 'display'">
            {{require}}
        </el-tag>
    </el-form-item>
    <el-form-item label="备注">
        <el-input v-if="mode !== 'display'" v-model="selStep.desc"/>
        <p v-else>{{selStep.desc}}</p>
    </el-form-item>
    <el-form-item label="输入" v-show="mode === 'editing-flow'">
        <div class="card input-card" v-for="(input, symbol) in selStep.inputs" :key="symbol">
            <div class="card-body input-card-body row">
                {{symbol}}
                <el-dropdown class="float-right" trigger="click">
                    <span class="el-dropdown-link">
                        选择局部变量<i class="el-icon-arrow-down el-icon--right"></i>
                    </span>
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item>黄金糕</el-dropdown-item>
                        <el-dropdown-item>狮子头</el-dropdown-item>
                        <el-dropdown-item>螺蛳粉</el-dropdown-item>
                        <el-dropdown-item disabled>双皮奶</el-dropdown-item>
                        <el-dropdown-item divided>蚵仔煎</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
            </div>
        </div>
    </el-form-item>
    <el-form-item label="输出" v-show="mode === 'editing-flow'">

    </el-form-item>
    <el-form-item label="代码">
        <el-input type="textarea" v-model="selStep.code" :autosize="{minRows: 5, maxRows: 100}" :disabled="mode === 'display'"/>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    props: {
        "selStep": Object,
        "preMode": String,
    },
    data() {
        return {
            mode: "display",
            options: [{
                label: "test",
                value: 1234
            }]
        }
    },
    created() {
        if (this.preMode && this.preMode.length !== 0) {
            this.mode = this.preMode
        }
    }
}
</script>

<style lang="scss">
.input-card {
    margin-bottom: 1vh;
}
.input-card:last-child {
    margin-bottom: 0 !important;
}
.input-card-body {
    padding: 0.1vh 0.5vw;
    .el-dropdown {
        cursor: pointer;
    }
}
</style>
