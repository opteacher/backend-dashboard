<template>
<el-form label-width="80px" :rules="formRules" ref="edit-relation-form" :model="relation">
    <el-form-item label="模块关系">
        <el-col :span="8">
            <el-form-item prop="model1">
                <el-select placeholder="起始模块" v-model="relation.model1">
                    <el-option v-for="model in model1s" :key="model.id" :label="model.name" :value="model.id"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col :span="3">
            <el-form-item prop="model1n">
                <el-select placeholder="数量" v-model="relation.model1n">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="line text-center" :span="2">-</el-col>
        <el-col :span="3">
            <el-form-item prop="model2n">
                <el-select placeholder="数量" v-model="relation.model2n">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="text-right" :span="8">
            <el-form-item prop="model2">
                <el-select placeholder="目标模块" v-model="relation.model2">
                    <el-option v-for="model in model2s" :key="model.id" :label="model.name" :value="model.id"/>
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
            model1: [{ required: true, message: "需要选择起始模块", trigger: "blur" }],
            model1n: [{ required: true, message: "需要指定起始模块数量", trigger: "blur" }],
            model2: [{ required: true, message: "需要选择目标模块", trigger: "blur" }],
            model2n: [{ required: true, message: "需要指定目标模块数量", trigger: "blur" }],
        },
        model1s: [],
        model2s: [],
        relation: {
            model1: null,
            model1n: null,
            model2: null,
            model2n: null,
        }
    }},
    created() {
        this.model1s = this.models
        this.model2s = this.models
    },
    watch: {
        model1(selMdlID) {
            this.model2s = _.filter(this.models, ele => ele.id !== selMdlID)
        },
        model2(selMdlID) {
            this.model1s = _.filter(this.models, ele => ele.id !== selMdlID)
        }
    },
    methods: {
        resetRelation() {
            this.model1s = this.models
            this.model2s = this.models
            this.relation = {
                model1: null,
                model1n: null,
                model2: null,
                model2n: null,
            }
        }
    }
}
</script>

