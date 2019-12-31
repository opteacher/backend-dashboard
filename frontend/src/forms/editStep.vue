<template>
<div>
    <el-form :inline="true" ref="form" :model="step" label-width="80px">
        <el-form-item label="操作KEY">
            <el-select class="w-100" v-model="step.key" placeholder="选择既存操作" @change="hdlSelStep">
                <el-option v-for="step in groupByLang[selLang]" :key="step.key" :label="step.key" :value="step.key"/>
            </el-select>
        </el-form-item>
        <el-form-item label="语言">
            <el-select class="float-right w-90" v-model="selLang" placeholder="选择语言" @change="step = {key: ''}">
                <el-option v-for="lang in Object.keys(groupByLang)" :key="lang" :label="lang" :value="lang"/>
            </el-select>
        </el-form-item>
        <el-form-item class="float-right">
            <el-button @click="hdlAddStep">添加模板步骤</el-button>
        </el-form-item>
    </el-form>
    <step-detail v-if="Object.keys(stepMapper).includes(step.key)" :selStep="step" preMode="editing-step" :locVars="locVars"/>
</div>
</template>

<script>
import _ from "lodash"

import backend from "../backend"
import stepDetail from "../forms/stepDetail"

export default {
    components: {
        "step-detail": stepDetail
    },
    props: {
        "prevIdx": Number,
        "locVars": Array
    },
    data() {
        return {
            steps: [],
            groupByLang: {},
            stepMapper: {},
            step: {key: ""},
            selLang: ""
        }
    },
    async created() {
        let res = await backend.qryAllTempStep()
        if (typeof res === "string") {
            this.$message.error(`查询模板步骤失败：${res}`)
        } else {
            this.steps = res.steps || []
            for (let step of this.steps) {
                this.stepMapper[step.key] = step
            }
            this.groupByLang = _.groupBy(this.steps, "lang")
            const langs = Object.keys(this.groupByLang)
            if (langs.length !== 0) {
                this.selLang = langs[0]
            }
        }
    },
    methods: {
        hdlAddStep() {

        },
        hdlSelStep() {
            this.step = _.cloneDeep(this.stepMapper[this.step.key])
        }
    }
}
</script>
