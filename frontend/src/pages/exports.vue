<template>
<dashboard>
    <div class="h-100 w-100" style="display: flex">
        <el-button style="align-self: center; margin: 0 auto" type="primary" @click="showExportDlg = true">导出</el-button>
    </div>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="导出项目" :visible.sync="showExportDlg" :modal-append-to-body="false" width="50vw">
        <exp-project ref="exp-project-form"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showExportDlg = false">取 消</el-button>
            <el-button type="primary" @click="exportProject">导 出</el-button>
        </div>
    </el-dialog>
</dashboard>
</template>

<script>
import backend from "../backend"
import dashboard from "../layouts/dashboard"
import expProject from "../forms/expProject"

export default {
    components: {
        "dashboard": dashboard,
        "exp-project": expProject
    },
    data() {
        return {
            showExportDlg: false
        }
    },
    methods: {
        async exportProject() {
            let form = this.$refs["exp-project-form"]
            form.$refs["exp-project-form"].validate(async valid => {
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
