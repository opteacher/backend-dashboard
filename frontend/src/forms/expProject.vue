<template>
<el-form :model="exportOption" label-width="80px" label-position="right" ref="exp-project-form">
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
    <el-divider>
        <a href="#" @click.prevent="clkSeniorLnk"><i :class="openSenior ? 'el-icon-arrow-up' : 'el-icon-arrow-down'"></i></a>
    </el-divider>
    <div v-show="openSenior">
        <el-form-item label="微服务">
            <el-checkbox v-model="exportOption.isMicoServ">生成为微服务</el-checkbox>
        </el-form-item>
    </div>
</el-form>
</template>

<script>
export default {
    data() {
        return {
            openSenior: false,
            exportOption: {
                name: "",
                routePrefix: "/api/v1",
                isMicoServ: false,
                database: {
                    name: "",
                    host: "127.0.0.1",
                    port: "3306",
                    type: "mysql",
                    username: "root",
                    password: "12345"
                },
            },
            exportTypes: [{
                title: "bl-kratos",
                value: "kratos"
            }],
            databaseTypes: [{
                title: "无数据库",
                value: ""
            }, {
                title: "MySQL",
                value: "mysql"
            }]
        }
    },
    methods: {
        clkSeniorLnk() {
            this.openSenior = !this.openSenior
        },
        chgProjName() {
            this.exportOption.database.name = this.exportOption.name
        }
    }
}
</script>
