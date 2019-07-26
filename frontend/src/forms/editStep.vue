<template>
<el-form ref="form" :model="step" label-width="80px">
    <el-form-item label="操作KEY">
        <el-col :span="18">
            <el-select class="w-100" v-model="step.operKey" placeholder="选择既存操作" @change="hdlSelOper">
                <el-option v-for="oper in opers" :key="oper.operKey" :label="oper.operKey" :value="oper.operKey"/>
            </el-select>
        </el-col>
        <el-col :span="6">
            <el-button class="float-right" @click="hdlAddOper">添加操作模板</el-button>
        </el-col>
    </el-form-item>
    <step-detail v-if="Object.keys(operMap).includes(step.operKey)" :selStep="step" preMode="editing-step" :locVars="stepInfo.locVars"/>
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
            opers: [],
            operMap: {},
            step: {
                operKey: ""
            }
        }
    },
    async created() {
        let res = await backend.qryStepTmp()
        if (typeof res === "string") {
            this.$message.error(`查询模板步骤失败：${res}`)
        } else {
            this.opers = res.steps || []
            for (let oper of this.opers) {
                this.operMap[oper.operKey] = oper
            }
        }
    },
    methods: {
        hdlAddOper() {

        },
        hdlSelOper() {
            this.step = _.cloneDeep(this.operMap[this.step.operKey])
        }
    }
}
</script>
