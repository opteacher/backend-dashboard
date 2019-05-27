<template>
<el-form label-width="80px">
    <el-form-item label="模块关系">
        <el-col :span="8">
            <el-select placeholder="起始模块" v-model="selBegMdl">
                <el-option v-for="model in begModels" :key="model.id" :label="model.name" :value="model.id"/>
            </el-select>
        </el-col>
        <el-col :span="3">
            <el-input v-model="begMdlNum"/>
        </el-col>
        <el-col class="line text-center" :span="2">-</el-col>
        <el-col :span="3">
            <el-input v-model="endMdlNum"/>
        </el-col>
        <el-col class="text-right" :span="8">
            <el-select placeholder="结束模块" v-model="selEndMdl">
                <el-option v-for="model in endModels" :key="model.id" :label="model.name" :value="model.id"/>
            </el-select>
        </el-col>
    </el-form-item>
</el-form>
</template>

<script>
import glbVar from "../global"

export default {
    data() { return {
        begModels: [],
        selBegMdl: null,
        begMdlNum: 0,
        endModels: [],
        selEndMdl: null,
        endMdlNum: 0
    }},
    created() {
        this.begModels = glbVar.models
        this.endModels = glbVar.models
    },
    watch: {
        selBegMdl(selMdlID) {
            this.endModels = _.filter(glbVar.models, ele => ele.id !== selMdlID)
        },
        selEndMdl(selMdlID) {
            this.begModels = _.filter(glbVar.models, ele => ele.id !== selMdlID)
        }
    }
}
</script>

