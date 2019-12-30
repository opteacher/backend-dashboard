<template>
<dashboard>
    <el-container>
        <el-main class="mx-auto">
            <el-form :model="exportOption" label-width="80px" label-position="right" ref="exp-project-form" style="width: 50vw">
                <el-form-item label="项目名" :rules="[
                    { required: true, message: '请输入项目名称', trigger: 'blur' },
                    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
                ]" prop="name">
                    <el-input v-model="exportOption.name" @change="chgProjName"/>
                </el-form-item>
                <el-form-item label="项目类别" :rules="[
                    { required: true, message: '请选择项目类别', trigger: 'change' }
                ]" prop="type">
                    <el-select v-model="exportOption.type" style="width: 100%">
                        <el-option v-for="typ in exportTypes" :key="typ.title" :label="typ.title" :value="typ.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="微服务">
                    <el-checkbox v-model="exportOption.isMicoServ">生成为微服务</el-checkbox>
                </el-form-item>
            </el-form>
        </el-main>
        <el-footer>
            <el-button class="float-right" type="primary" @click="exportProject">导出</el-button>
        </el-footer>
    </el-container>
</dashboard>
</template>

<script>
import backend from "../backend"
import dashboard from "../layouts/dashboard"

export default {
    components: {
        "dashboard": dashboard
    },
    data() {
        return {
            exportOption: {
                name: "",
                isMicoServ: false
            },
            exportTypes: [{
                title: "bl-kratos",
                value: "kratos"
            }]
        }
    },
    methods: {
        async exportProject() {
            let form = this.$refs["exp-project-form"]
            form.validate(async valid => {
                if (valid) {
                    // if (form.exportOption.name.slice(-4).toLowerCase() !== ".zip") {
                    //     form.exportOption.name += ".zip"
                    // }
                    let res = await backend.export(form.exportOption)
                    if (typeof res === "string") {
                        this.$message.error(`导出时发生错误：${res}`)
                    } else {
                        window.open(res.url, "_blank")
                    }
                    this.showExportDlg = false
                } else {
                    return false
                }
            })
        }
    }
}
</script>
