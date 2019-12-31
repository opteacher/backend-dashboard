<template>
<dashboard>
    <el-container>
        <el-main class="mx-auto">
            <el-form :model="options" label-width="80px" label-position="right" ref="exp-project-form" style="width: 50vw">
                <el-form-item label="项目名" :rules="[
                    { required: true, message: '请输入项目名称', trigger: 'blur' },
                    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
                ]" prop="name">
                    <el-input v-model="options.name"/>
                </el-form-item>
                <el-form-item label="项目类别" :rules="[
                    { required: true, message: '请选择项目类别', trigger: 'change' }
                ]" prop="type">
                    <el-select v-model="options.type" style="width: 100%" @change="hdlSelFramework">
                        <el-option v-for="fwk in frameworks" :key="fwk.id" :label="fwk.name" :value="fwk.id"/>
                    </el-select>
                </el-form-item>
                <el-form-item v-if="fwkSupports.includes('micoService')" label="微服务">
                    <el-checkbox v-model="options.components.micoService.enable">生成为微服务</el-checkbox>
                </el-form-item>
                <el-form-item v-if="fwkSupports.includes('frontendEnv')" label="前端支持">
                    <el-checkbox v-model="options.components.frontendEnv.enable">配置前端支持</el-checkbox>
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
            options: {
                name: "",
                type: "",
                components: {
                    micoService: {enable: false},
                    frontendEnv: {enable: false}
                }
            },
            frameworks: [],
            fwkSupports: []
        }
    },
    async created() {
        const res = await backend.qryAllExpTypes()
        if (typeof res === "string") {
            this.$message.error(`查询导出项目类别时发生错误：${res}`)
        } else {
            this.frameworks = res.frameworks
        }
    },
    methods: {
        hdlSelFramework() {
            const selFramework = this.frameworks.find(fwk => fwk.id === this.options.type)
            this.fwkSupports = selFramework.supports
        },
        async exportProject() {
            let form = this.$refs["exp-project-form"]
            form.validate(async valid => {
                if (valid) {
                    let res = await backend.export(this.options)
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
