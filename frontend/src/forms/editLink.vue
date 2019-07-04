<template>
<el-form label-width="80px" :rules="formRules" ref="edit-link-form" :model="link">
    <el-form-item label="模块关系">
        <el-col :span="8">
            <el-form-item prop="modelName1">
                <el-select placeholder="起始模块" v-model="link.modelName1">
                    <el-option v-for="model in model1s" :key="model.id" :label="model.name" :value="model.name"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col :span="3">
            <el-form-item prop="modelNumber1">
                <el-select placeholder="数量" v-model="link.modelNumber1">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="line text-center" :span="2">-</el-col>
        <el-col :span="3">
            <el-form-item prop="modelNumber2">
                <el-select placeholder="数量" v-model="link.modelNumber2">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="text-right" :span="8">
            <el-form-item prop="modelName2">
                <el-select placeholder="目标模块" v-model="link.modelName2">
                    <el-option v-for="model in model2s" :key="model.id" :label="model.name" :value="model.name"/>
                </el-select>
            </el-form-item>
        </el-col>
    </el-form-item>
</el-form>
</template>

<script>
export default {
    props: {
        "models": Array
    },
    data() { return {
        formRules: {
            modelName1: [{ required: true, message: "需要选择起始模块", trigger: "blur" }],
            modelNumber1: [{ required: true, message: "需要指定起始模块数量", trigger: "blur" }],
            modelName2: [{ required: true, message: "需要选择目标模块", trigger: "blur" }],
            modelNumber2: [{ required: true, message: "需要指定目标模块数量", trigger: "blur" }],
        },
        model1s: [],
        model2s: [],
        link: {
            symbol: "",
            modelName1: null,
            modelNumber1: null,
            modelName2: null,
            modelNumber2: null,
        }
    }},
    created() {
        this.model1s = this.models
        this.model2s = this.models
    },
    computed: {
        modelName1() {
            return this.link.modelName1
        },
        modelName2() {
            return this.link.modelName2
        }
    },
    watch: {
        modelName1(name) {
            this.model2s = _.filter(this.models, ele => ele.name !== name)
        },
        modelName2(name) {
            this.model1s = _.filter(this.models, ele => ele.name !== name)
        }
    },
    methods: {
        resetLink() {
            this.model1s = this.models
            this.model2s = this.models
            this.link = {
                symbol: "",
                modelName1: null,
                modelNumber1: null,
                modelName2: null,
                modelNumber2: null,
            }
        }
    }
}
</script>

