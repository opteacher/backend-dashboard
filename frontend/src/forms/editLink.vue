<template>
<el-form label-width="80px" :rules="formRules" ref="edit-link-form" :model="link">
    <el-form-item label="模块关系">
        <el-col :span="8">
            <el-form-item prop="mname1">
                <el-select placeholder="起始模块" v-model="link.mname1">
                    <el-option v-for="model in model1s" :key="model.id" :label="model.name" :value="model.name"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col :span="3">
            <el-form-item prop="mnumber1">
                <el-select placeholder="数量" v-model="link.mnumber1">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="line text-center" :span="2">-</el-col>
        <el-col :span="3">
            <el-form-item prop="mnumber2">
                <el-select placeholder="数量" v-model="link.mnumber2">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="text-right" :span="8">
            <el-form-item prop="mname2">
                <el-select placeholder="目标模块" v-model="link.mname2">
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
            mname1: [{ required: true, message: "需要选择起始模块", trigger: "blur" }],
            mnumber1: [{ required: true, message: "需要指定起始模块数量", trigger: "blur" }],
            mname2: [{ required: true, message: "需要选择目标模块", trigger: "blur" }],
            mnumber2: [{ required: true, message: "需要指定目标模块数量", trigger: "blur" }],
        },
        model1s: [],
        model2s: [],
        link: {
            symbol: "",
            mname1: null,
            mnumber1: null,
            mname2: null,
            mnumber2: null,
        }
    }},
    created() {
        this.model1s = this.models
        this.model2s = this.models
    },
    computed: {
        mname1() {
            return this.link.mname1
        },
        mname2() {
            return this.link.mname2
        }
    },
    watch: {
        mname1(name) {
            this.model2s = _.filter(this.models, ele => ele.name !== name)
        },
        mname2(name) {
            this.model1s = _.filter(this.models, ele => ele.name !== name)
        }
    },
    methods: {
        resetLink() {
            this.model1s = this.models
            this.model2s = this.models
            this.link = {
                symbol: "",
                mname1: null,
                mnumber1: null,
                mname2: null,
                mnumber2: null,
            }
        }
    }
}
</script>

