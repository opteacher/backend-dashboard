<template>
<el-form label-width="80px" :rules="formRules" ref="edit-relation-form" :model="relation">
    <el-form-item label="模块关系">
        <el-col :span="8">
            <el-form-item prop="selBegMdl">
                <el-select placeholder="起始模块" v-model="relation.selBegMdl">
                    <el-option v-for="model in begModels" :key="model.id" :label="model.name" :value="model.id"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col :span="3">
            <el-form-item prop="begMdlNum">
                <el-select placeholder="数量" v-model="relation.begMdlNum">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="line text-center" :span="2">-</el-col>
        <el-col :span="3">
            <el-form-item prop="endMdlNum">
                <el-select placeholder="数量" v-model="relation.endMdlNum">
                    <el-option key="1" label="1" :value="1"/>
                    <el-option key="2" label="*" :value="0"/>
                </el-select>
            </el-form-item>
        </el-col>
        <el-col class="text-right" :span="8">
            <el-form-item prop="selEndMdl">
                <el-select placeholder="目标模块" v-model="relation.selEndMdl">
                    <el-option v-for="model in endModels" :key="model.id" :label="model.name" :value="model.id"/>
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
            selBegMdl: [{ required: true, message: "需要选择起始模块", trigger: "blur" }],
            begMdlNum: [{ required: true, message: "需要指定起始模块数量", trigger: "blur" }],
            selEndMdl: [{ required: true, message: "需要选择目标模块", trigger: "blur" }],
            endMdlNum: [{ required: true, message: "需要指定目标模块数量", trigger: "blur" }],
        },
        begModels: [],
        endModels: [],
        relation: {
            selBegMdl: null,
            begMdlNum: null,
            selEndMdl: null,
            endMdlNum: null,
        }
    }},
    created() {
        this.begModels = this.models
        this.endModels = this.models
    },
    watch: {
        selBegMdl(selMdlID) {
            this.endModels = _.filter(this.models, ele => ele.id !== selMdlID)
        },
        selEndMdl(selMdlID) {
            this.begModels = _.filter(this.models, ele => ele.id !== selMdlID)
        }
    },
    methods: {
        resetRelation() {
            this.begModels = this.models
            this.endModels = this.models
            this.relation = {
                selBegMdl: null,
                begMdlNum: null,
                selEndMdl: null,
                endMdlNum: null,
            }
        }
    }
}
</script>

