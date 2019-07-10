<template>
<el-form ref="form" label-width="80px">
    <el-form-item label="操作标识">
        {{selStep.operKey}}
    </el-form-item>
    <el-form-item label="依赖" v-show="selStep.requires && selStep.requires.length !== 0">
        <el-tag v-for="require in selStep.requires" :key="require" :closable="mode === 'editing'">
            {{require}}
        </el-tag>
    </el-form-item>
    <el-form-item label="备注">
        <el-input v-if="mode === 'editing'" v-model="selStep.desc"/>
        <p v-else>{{selStep.desc}}</p>
    </el-form-item>
    <el-form-item label="代码">
        <el-input type="textarea" v-model="selStep.code" :autosize="{minRows: 5, maxRows: 100}" :disabled="mode !== 'editing'"/>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    props: {
        "selStep": Object
    },
    data() {
        return {
            mode: "display"
        }
    }
}
</script>

