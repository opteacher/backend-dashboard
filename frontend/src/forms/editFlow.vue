<template>
<el-form ref="form" :model="step" label-width="100px">
    <el-form-item label="操作KEY">
        <el-col :span="18">
            <el-select class="w-100" v-model="step.operKey" placeholder="选择既存操作" @change="hdlSelOper">
                <el-option v-for="oper in opers" :key="oper.operKey" :label="oper.operKey" :value="oper.operKey"/>
            </el-select>
        </el-col>
        <el-col :span="6">
            <el-button class="float-right" @click="hdlAddOper">添加操作</el-button>
        </el-col>
    </el-form-item>
    <step-detail v-if="operMap && _.keys(operMap).includes(step.operKey)"/>
</el-form>
</template>

<script>
import _ from "lodash"

import stepBkd from "../async/step"
import stepDetail from "../forms/stepDetail"

export default {
    components: {
        "step-detail": stepDetail
    },
    data() {
        return {
            opers: [],
            operMap: {},
            step: {
                operKey: "",
                requires: []
            }
        }
    },
    async created() {
        let res = await stepBkd.qryTmp()
        if (typeof res === "string") {
            this.$message(`查询模板步骤失败：${res}`)
        } else {
            this.opers = res.data.data.steps || []
            for (let oper of this.opers) {
                this.operMap[oper.operKey] = oper
            }
        }
    },
    methods: {
        hdlAddOper() {

        },
        hdlSelOper() {
            console.log(this.step.operKey)
        }
    }
}
</script>
