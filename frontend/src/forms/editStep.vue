<template>
<el-form ref="form" :model="step" label-width="80px">
    <el-form-item label="操作KEY">
        <el-col :span="18">
            <el-select class="w-100" v-model="step.key" placeholder="选择既存操作" @change="hdlSelStep">
                <el-option v-for="step in steps" :key="step.key" :label="step.key" :value="step.key"/>
            </el-select>
        </el-col>
        <el-col :span="6">
            <el-button class="float-right" @click="hdlAddStep">添加操作模板</el-button>
        </el-col>
    </el-form-item>
    <step-detail v-if="Object.keys(stepMap).includes(step.key)" :selStep="step" preMode="editing-step" :locVars="stepInfo.locVars"/>
</el-form>
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
        "stepInfo": Object
    },
    data() {
        return {
            steps: [],
            stepMap: {},
            step: {
                key: ""
            }
        }
    },
    async created() {
        let res = await backend.qryAllTempStep()
        if (typeof res === "string") {
            this.$message.error(`查询模板步骤失败：${res}`)
        } else {
            this.steps = res.steps || []
            for (let step of this.steps) {
                this.stepMap[step.key] = step
            }
        }
    },
    methods: {
        hdlAddStep() {

        },
        hdlSelStep() {
            this.step = _.cloneDeep(this.stepMap[this.step.key])
        }
    }
}
</script>
