<template>
<el-form :model="daoGroup" :rules="formRules" label-width="80px" ref="edit-dao-group-form">
    <el-form-item label="组名" prop="name">
        <el-input placeholder="输入组名" v-model="daoGroup.name"/>
    </el-form-item>
    <el-form-item label="类别" prop="category">
        <el-select class="w-100" v-model="daoGroup.category">
            <el-option v-for="category in categories" :key="category.value" :label="category.title" :value="category.value"/>
        </el-select>
    </el-form-item>
    <el-form-item label="语言" prop="language">
        <el-select class="w-100" v-model="daoGroup.language">
            <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang"/>
        </el-select>
    </el-form-item>
    <el-form-item label="实现" prop="implement">
        <el-col :span="20">
            <el-select class="w-100" v-model="daoGroup.implement">
                <el-option v-for="impl in impls" :key="impl.value" :label="impl.title" :value="impl.value"/>
            </el-select>
        </el-col>
        <el-col :span="4">
            <el-button class="float-right" icon="el-icon-refresh" circle/>
        </el-col>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    data() {
        return {
            formRules: {
                name: [{ required: true, message: "需要指定组名", trigger: "blur" }],
                category: [{ required: true, message: "需要选择类别", trigger: "blur" }],
                language: [{ required: true, message: "需要选择语言", trigger: "blur" }]
            },
            daoGroup: {
                name: "",
                category: "",
                language: "",
                implement: ""
            },
            categories: [{
                title: "数据库",
                value: "database"
            }, {
                title: "缓存",
                value: "cache"
            }],
            impls: [{
                title: "MySQL",
                value: "mysql"
            }],
            languages: ["golang", "javascript", "java"]
        }
    }
}
</script>

