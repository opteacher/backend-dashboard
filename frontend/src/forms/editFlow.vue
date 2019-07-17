<template>
<el-form ref="form" :model="flow" label-width="80px">
    <el-form-item label="操作KEY">
        <el-col :span="18">
            <el-select class="w-100" v-model="flow.operKey" placeholder="选择既存操作" @change="hdlSelOper">
                <el-option v-for="oper in opers" :key="oper.operKey" :label="oper.operKey" :value="oper.operKey"/>
            </el-select>
        </el-col>
        <el-col :span="6">
            <el-button class="float-right" @click="hdlAddOper">添加操作模板</el-button>
        </el-col>
    </el-form-item>
    <step-detail v-if="Object.keys(operMap).includes(flow.operKey)" :selStep="flow" preMode="editing-flow" :locVars="flowInfo.locVars"/>
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
    props: {
        "flowInfo": Object
    },
    data() {
        return {
            opers: [],
            operMap: {},
            flow: {
                operKey: ""
            }
        }
    },
    async created() {
        let res = await stepBkd.qryTmp()
        if (typeof res === "string") {
            this.$message.error(`查询模板步骤失败：${res}`)
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
            this.flow = _.cloneDeep(this.operMap[this.flow.operKey])
        }
    }
}
</script>
