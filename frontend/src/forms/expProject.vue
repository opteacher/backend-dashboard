<template>
<el-form :model="exportOption" label-width="80px" label-position="right" ref="exp-project-form">
    <el-form-item label="项目名" :rules="[
        { required: true, message: '请输入项目名称', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
    ]" prop="name">
        <el-input v-model="exportOption.name"/>
    </el-form-item>
    <el-form-item label="项目类别" :rules="[
        { required: true, message: '请选择项目类别', trigger: 'change' }
    ]" prop="type">
        <el-select v-model="exportOption.type" style="width: 100%">
            <el-option v-for="typ in exportTypes" :key="typ.title" :label="typ.title" :value="typ.value"/>
        </el-select>
    </el-form-item>
    <el-form-item label="路由前缀" :rules="[
        { required: true, message: '请输入路由前缀', trigger: 'blur' }
    ]" prop="prefix">
        <el-input v-model="exportOption.prefix"/>
    </el-form-item>
    <el-form-item label="微服务">
        <el-checkbox v-model="exportOption.micoSvc">生成为微服务</el-checkbox>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    data() {return {
        exportOption: {
            name: "",
            prefix: "/api/v1",
            micoSvc: false
        },
        exportTypes: [{
            title: "bl-kratos",
            value: "kratos-demo"
        }],
    }}
}
</script>
