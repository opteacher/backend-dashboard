<template>
<el-form :model="model" label-width="80px" label-position="right" ref="edit-model-form">
    <el-form-item label="模块名" :rules="[
        { required: true, message: '请输入模块名称', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
    ]" prop="name">
        <el-input v-model="model.name"/>
    </el-form-item>
    <el-form-item label="属性">
        <el-row class="m-0" :gutter="5">
            <el-col class="p-0" :span="18">
                <el-form-item :rules="[
                    { required: true, message: '请输入属性名称', trigger: 'blur' },
                    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
                ]" prop="propName">
                    <el-input v-model="model.propName"/>
                </el-form-item>
            </el-col>
            <el-col class="p-0" :span="6">
                <el-button class="float-right" @click="addProp">添加</el-button>
            </el-col>
        </el-row>
    </el-form-item>
    <el-form-item v-for="(prop, index) in model.props" :key="prop.name" :label="prop.name">
        <el-row class="m-0" :gutter="5">
            <el-col class="p-0" :span="18">
                <el-form-item :rules="[
                    { required: true, message: '请选择属性类型', trigger: 'change' }
                ]" :prop="`props.${index}.type`">
                    <el-select v-model="prop.type" style="width: 100%">
                        <el-option v-for="option in propOptions" :key="option.title" :label="option.title" :value="option.value"/>
                    </el-select>
                </el-form-item>
            </el-col>
            <el-col class="p-0" :span="6">
                <el-button class="float-right" @click.prevent="deleteProp(prop.name)">删除</el-button>
            </el-col>
        </el-row>
    </el-form-item>
    <el-form-item label="RPC接口">
        <el-checkbox-group class="w-100" v-model="model.methods">
            <el-checkbox-button label="POST" name="methods"/>
            <el-checkbox-button label="DELETE" name="methods"/>
            <el-checkbox-button label="PUT" name="methods"/>
            <el-checkbox-button label="GET" name="methods"/>
            <el-checkbox-button label="ALL" name="methods"/>
        </el-checkbox-group>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    props: {
        "input-model": Object
    },
    data() { return {
        model: {
            id: "",
            name: "",
            props: [],
            propName: "",
            methods: [],
            x: 0,
            y: 0,
            width: 400,
            height: 300
        },
        propOptions: [{
            title: "文本",
            value: "string"
        }, {
            title: "数字",
            value: "int32"
        }, {
            title: "日期",
            value: "uint64"
        }, {
            title: "布尔",
            value: "bool"
        }]
    }},
    created() {
        if (this["input-model"] && this["input-model"].name) {
            this.model.name = this["input-model"].name
        }
        if (this["input-model"] && this["input-model"].props) {
            this.model.props = this["input-model"].props
        }
    },
    methods: {
        addProp() {
            this.$refs["edit-model-form"].validate(valid => {
                if (valid) {
                    this.model.props.push({
                        name: this.model.propName
                    })
                    this.model.propName = ""
                } else {
                    return false
                }
            })
        },
        deleteProp(propName) {
            this.model.props.pop(ele => ele.name === propName)
        },
        resetModel() {
            this.model = {
                id: "",
                name: "",
                props: [],
                propName: "",
                methods: [],
                x: 0,
                y: 0,
                width: 400,
                height: 300
            }
        }
    }
}
</script>
